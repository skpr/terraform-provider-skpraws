// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"fmt"
	"github.com/skpr/terraform-provider-skpraws/internal/provider"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccManagedLoginBrandingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.testAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccManagedLoginBrandingResourceConfig(true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"cognito_managed_login_branding.test",
						tfjsonpath.New("use_cognito_provided_values"),
						knownvalue.Bool(true),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "cognito_managed_login_branding.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"use_cognito_provided_values"},
			},
			// Update and Read testing
			{
				Config: testAccManagedLoginBrandingResourceConfig(false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"cognito_managed_login_branding.test",
						tfjsonpath.New("use_cognito_provided_values"),
						knownvalue.Bool(false),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccManagedLoginBrandingResourceConfig(useCognitoProvidedValues bool) string {
	return fmt.Sprintf(`
resource "cognito_managed_login_branding" "test" {
  client_id = "abc"
  user_pool_id = "abc"
  use_cognito_provided_values = %t
  settings = "{}"
}
`, useCognitoProvidedValues)
}
