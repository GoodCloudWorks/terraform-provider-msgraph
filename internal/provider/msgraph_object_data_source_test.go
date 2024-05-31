package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMsGraphObjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig + `
					data "msgraph_provider_config" "this" {}
					data "msgraph_object" "organization" {
						id = "organization/${data.msgraph_provider_config.this.tenant_id}" 
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.msgraph_object.organization", "output.@odata.context", "https://graph.microsoft.com/v1.0/$metadata#organization/$entity"),
				),
			},
			{
				Config: testProviderConfig + `
					data "msgraph_provider_config" "this" {}
					data "msgraph_object" "organization" {
						id          = "organization/${data.msgraph_provider_config.this.tenant_id}" 
						api_version = "beta"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.msgraph_object.organization", "output.@odata.context", "https://graph.microsoft.com/beta/$metadata#organization/$entity"),
				),
			},
		},
	})
}
