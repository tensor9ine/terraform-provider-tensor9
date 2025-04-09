// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccExampleResource(t *testing.T) {

	reactor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/stack/tf/react" {
			w.Header().Set("Content-Type", "application/json")
			// Respond with a dummy infra_id
			fmt.Fprint(w, `{"infra_id": "0000000000000000000000000000beef"}`)
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
				Config: testAccExampleResourceConfig(reactor.URL, "{}", "Terraform", "0000000000000000:0000000000000000:0000000000000000", "rsx_a", map[string]string{"key1": "value1"}),
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
				ImportStateVerifyIgnore: []string{"template", "template_fmt", "projection_id", "rsx_id", "infra_id", "properties"}, // TODO: remove this once ::Read is implemented
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig(reactor.URL, "{}", "Terraform", "0000000000000000:0000000000000000:0000000000000000", "rsx_a", map[string]string{"key1": "value1"}),
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

func testAccExampleResourceConfig(endpoint string, template string, templateFmt string, projectionId string, rsxId string, properties map[string]string) string {
	var props string
	if len(properties) == 0 {
		props = "{}"
	} else {
		props = "{\n"
		for k, v := range properties {
			props += fmt.Sprintf(`    %s = %q`+"\n", k, v)
		}
		props += "}\n"
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
  properties = %s
}
`, endpoint, template, templateFmt, projectionId, rsxId, props)
}
