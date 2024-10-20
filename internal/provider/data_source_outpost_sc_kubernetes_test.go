package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOutpostServiceConnectionsKubernetes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutpostServiceConnectionKubernetesSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.local-cluster", "name", "Local Kubernetes Cluster"),
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.local-cluster", "local", "true"),
				),
			},
			{
				Config: testAccDataSourceOutpostServiceConnectionKubernetesSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.remote-cluster", "name", "Remote Kubernetes Cluster"),
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.remote-cluster", "local", "false"),
				),
			},
		},
	})
}

const testAccDataSourceOutpostServiceConnectionKubernetesSimple = `
data "authentik_service_connection_kubernetes" "local-cluster" {
  name = "Local Kubernetes Cluster"
  local = true
}

data "authentik_service_connection_kubernetes" "remote-cluster" {
  name = "Remote Kubernetes Cluster"
}
`
