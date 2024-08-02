package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSCIMSourcePropertyMapping(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSCIMSourcePropertyMapping(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_source_scim.name", "name", rName),
				),
			},
			{
				Config: testAccResourceSCIMSourcePropertyMapping(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_source_scim.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceSCIMSourcePropertyMapping(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_source_scim" "name" {
  name         = "%[1]s"
  expression   = "return True"
}
`, name)
}
