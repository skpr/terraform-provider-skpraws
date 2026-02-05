package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccManagedLoginBrandingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccManagedLoginBrandingResourceConfig("{}"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"skpraws_managed_login_branding.test",
						tfjsonpath.New("settings"),
						knownvalue.StringExact("{}"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "skpraws_managed_login_branding.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"settings"},
			},
			// Update and Read testing
			{
				Config: testAccManagedLoginBrandingResourceConfig("{test = true}"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"skpraws_managed_login_branding.test",
						tfjsonpath.New("settings"),
						knownvalue.StringExact("{test = true}"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccManagedLoginBrandingResourceConfig(settings string) string {
	return fmt.Sprintf(`
resource "skpraws_managed_login_branding" "test" {
  client_id = "abc"
  user_pool_id = "abc"
  settings = "%s"
  assets = [
    {
      category = "PAGE_BACKGROUND"
      color_mode = "LIGHT"
      bytes = "ABCDEFGH"
    },
  ]
}
`, settings)
}
