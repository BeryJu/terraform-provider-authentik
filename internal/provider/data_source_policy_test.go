package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicy(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicySimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy.default-authentication-flow-password-stage", "name", "default-authentication-flow-password-stage"),
				),
			},
			{
				Config:      testAccDataSourcePolicyNotExisting,
				ExpectError: regexp.MustCompile(`No matching policy found`),
			},
		},
	})
}

const testAccDataSourcePolicySimple = `
data "authentik_policy" "default-authentication-flow-password-stage" {
  name = "default-authentication-flow-password-stage"
}
`

const testAccDataSourcePolicyNotExisting = `
data "authentik_policy" "not-exiting-policy" {
  name = "not-exiting-policy"
}
`
