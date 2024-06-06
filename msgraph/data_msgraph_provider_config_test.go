package msgraph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMsGraphProviderConfigDataSource(t *testing.T) {
	const resourceName = "data.msgraph_provider_config.this"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: defaultProviderConfigWith(`
					data "msgraph_provider_config" "this" {}
					`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
		},
	})
}
