// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

			switch evt.EvtType {
			case "Create":
				println("Reactor handling Create event")

				rsx := evt.LoFiTwinRsx
				infraId := "000000000000000000000000deadbeef"
				propertiesOut := make(map[string]string)
				for k, v := range *rsx.Vars {
					propertiesOut[k] = v
				}
				propertiesOut["new1"] = "value1"
				propertiesOut["new2"] = "value2"

				for k, v := range propertiesOut {
					println(fmt.Sprintf("%s = %s", k, v))
				}

				evtResult := TfRsxEvtResult{
					EvtType:    evt.EvtType,
					RsxType:    evt.RsxType,
					ResultType: "Created",
					LoFiTwinRsx: &Delta[TfLoFiTwinRsx]{
						Before: rsx,
						After: &TfLoFiTwinRsx{
							RsxId:        rsx.RsxId,
							Template:     rsx.Template,
							ProjectionId: rsx.ProjectionId,
							Vars:         &propertiesOut,
							Schema:       rsx.Schema,
							InfraId:      &infraId,
						},
					},
					Reason: nil,
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
			case "Read":
				println("Reactor handling Read event")
				println("TODO: implement Read event handling")
				return
			case "Update":
				println("Reactor handling Update event")
				println("TODO: implement Update event handling")
				return
			case "Delete":
				println("Reactor handling Delete event")
				println("TODO: implement Delete event handling")
				return
			default:
				http.Error(w, "unknown evt type", http.StatusBadRequest)
				return
			}
		}

		http.NotFound(w, r)
	}))
	defer reactor.Close()

	schema := map[string]TfRsxPropType{"original1": "Str", "new1": "Str", "new2": "Str"}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig(
					reactor.URL,
					"{}",
					"Terraform",
					"0000000000000000:0000000000000000:0000000000000000",
					"rsx_a",
					map[string]string{"original1": "value1"},
					schema,
				),
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
				ImportStateVerifyIgnore: []string{"template", "template_fmt", "projection_id", "rsx_id", "infra_id", "vars", "schema", "outputs"}, // TODO: remove this once ::Read is implemented
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig(
					reactor.URL,
					"{}",
					"Terraform",
					"0000000000000000:0000000000000000:0000000000000000",
					"rsx_a",
					map[string]string{"original1": "value1"},
					schema,
				),
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

func testAccExampleResourceConfig(
	endpoint string,
	template string,
	templateFmt string,
	projectionId string,
	rsxId string,
	vars map[string]string,
	schema map[string]TfRsxPropType,
) string {
	var varsStr string
	if len(vars) == 0 {
		varsStr = "{}"
	} else {
		varsStr = "{\n"
		for k, v := range vars {
			varsStr += fmt.Sprintf(`    %s = %q`+"\n", k, v)
		}
		varsStr += "}\n"
	}

	var schemaStr string
	if len(schema) == 0 {
		schemaStr = "{}"
	} else {
		schemaStr = "{\n"
		for k, v := range schema {
			schemaStr += fmt.Sprintf(`    %q = "%s"`+"\n", k, v)
		}
		schemaStr += "}\n"
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
  vars = %s
  schema = %s
}
`, endpoint, template, templateFmt, projectionId, rsxId, varsStr, schemaStr)
}
