package data

import (
	"testing"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMsGraphProviderConfigDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptance.Config(`
					data "msgraph_provider_config" "this" {}
					`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.msgraph_provider_config.this", "tenant_id"),
					resource.TestCheckResourceAttrSet("data.msgraph_provider_config.this", "client_id"),
					resource.TestCheckResourceAttrSet("data.msgraph_provider_config.this", "object_id"),
				),
			},
		},
	})
}
