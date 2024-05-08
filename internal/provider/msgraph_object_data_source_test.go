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
				Config: testProviderConfig + `data "msgraph_object" "me" { id = "me" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.msgraph_object.me", "output.@odata.context", "https://graph.microsoft.com/v1.0/$metadata#users/$entity"),
				),
			},
		},
	})
}
