package msgraph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMsGraphObjectDataSource(t *testing.T) {
	const resourceName = "data.msgraph_object.organization"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: defaultProviderConfigWith(`
					data "msgraph_provider_config" "this" {}
					data "msgraph_object" "organization" {
						id = "organization/${data.msgraph_provider_config.this.tenant_id}" 
					}
					`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "output.@odata.context", "https://graph.microsoft.com/v1.0/$metadata#organization/$entity"),
					resource.TestCheckResourceAttr(resourceName, "collection", "organization"),
				),
			},
			{
				Config: defaultProviderConfigWith(`
					data "msgraph_provider_config" "this" {}
					data "msgraph_object" "organization" {
						id          = "organization/${data.msgraph_provider_config.this.tenant_id}" 
						api_version = "beta"
					}
					`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "output.@odata.context", "https://graph.microsoft.com/beta/$metadata#organization/$entity"),
				),
			},
		},
	})
}
