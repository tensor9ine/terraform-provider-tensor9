// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestLoFiTwinRsx(t *testing.T) {

	reactor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/stack/tf/react" {
			w.Header().Set("Content-Type", "application/json")

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read body", http.StatusBadRequest)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					println(fmt.Sprintf("Error closing request body: %v\n", err))
				}
			}(r.Body)

			var evt TfRsxEvt

			err = json.Unmarshal(bodyBytes, &evt)
			if err != nil {
				http.Error(w, "failed to unmarshal evt", http.StatusBadRequest)
				return
			}

			println(fmt.Sprintf("Reactor received evt: %s %s", evt.EvtType, evt.RsxType))

			rsx := evt.LoFiTwinRsx
			infraId := "000000000000000000000000deadbeef"
			propertiesOut := make(map[string]string)
			for k, v := range *rsx.Properties {
				propertiesOut[k] = v
			}
			propertiesOut["new1"] = "value1"
			propertiesOut["new2"] = "value2"

			for k, v := range propertiesOut {
				println(fmt.Sprintf("%s = %s", k, v))
			}

			evtResult := TfRsxEvtResult{
				ResultType: "Created",
				EvtType:    "Create",
				BeforeRsx:  *rsx,
				AfterRsx: TfLoFiTwinRsx{
					RsxId:        rsx.RsxId,
					Template:     rsx.Template,
					ProjectionId: rsx.ProjectionId,
					Properties:   &propertiesOut,
					InfraId:      &infraId,
				},
			}

			evtResultJson, err := json.Marshal(evtResult)
			if err != nil {
				http.Error(w, "failed to marshal evt result", http.StatusBadRequest)
				return
			}

			println("Reactor sending result", string(evtResultJson))

			_, err = fmt.Fprint(w, string(evtResultJson))
			if err != nil {
				http.Error(w, "failed to write response", http.StatusInternalServerError)
				return
			}

			return
		}

		http.NotFound(w, r)
	}))
	defer reactor.Close()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig(reactor.URL, "{}", "Terraform", "0000000000000000:0000000000000000:0000000000000000", "rsx_a", map[string]string{"original1": "value1"}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tensor9_lofi_twin.test_twin",
						tfjsonpath.New("template"),
						knownvalue.StringExact("{}"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "tensor9_lofi_twin.test_twin",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"template", "template_fmt", "projection_id", "rsx_id", "infra_id", "properties_in", "properties_out"}, // TODO: remove this once ::Read is implemented
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig(reactor.URL, "{}", "Terraform", "0000000000000000:0000000000000000:0000000000000000", "rsx_a", map[string]string{"original1": "value1"}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tensor9_lofi_twin.test_twin",
						tfjsonpath.New("template"),
						knownvalue.StringExact("{}"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig(endpoint string, template string, templateFmt string, projectionId string, rsxId string, propertiesIn map[string]string) string {
	var propsIn string
	if len(propertiesIn) == 0 {
		propsIn = "{}"
	} else {
		propsIn = "{\n"
		for k, v := range propertiesIn {
			propsIn += fmt.Sprintf(`    %s = %q`+"\n", k, v)
		}
		propsIn += "}\n"
	}

	return fmt.Sprintf(`
provider "tensor9" {
  endpoint = %[1]q
  api_key  = "deadbeef"
}
resource "tensor9_lofi_twin" "test_twin" {
  template = %[2]q
  template_fmt = %[3]q
  projection_id = %[4]q
  rsx_id = %[5]q
  properties_in = %s
}
`, endpoint, template, templateFmt, projectionId, rsxId, propsIn)
}
