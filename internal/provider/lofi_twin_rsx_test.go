// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("{}", "Terraform", "0000000000000000:0000000000000000:0000000000000000"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tensor9_lofi_twin.test_twin",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
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
				ImportStateVerifyIgnore: []string{"template", "template_fmt", "projection_id"},
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig("{}", "Terraform", "0000000000000000:0000000000000000:0000000000000000"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"tensor9_lofi_twin.test_twin",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
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

func testAccExampleResourceConfig(template string, templateFmt string, projectionId string) string {
	return fmt.Sprintf(`
provider "tensor9" {
  endpoint = "http://localhost:9000"
  api_key  = "deadbeef"
}
resource "tensor9_lofi_twin" "test_twin" {
  template = %[1]q
  template_fmt = %[2]q
  projection_id = %[3]q
}
`, template, templateFmt, projectionId)
}
