package msgraph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccMsGraphObjectResource(t *testing.T) {
	const resourceName = "msgraph_object.group"
	groupName := acctest.RandString(10)
	updatedGroupName := fmt.Sprintf("%s-updated", groupName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: msGraphGroupResourceConfig(groupName, groupName), ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "output.@odata.context", "https://graph.microsoft.com/v1.0/$metadata#groups/$entity"),
					resource.TestCheckResourceAttr(resourceName, "output.displayName", groupName),
					resource.TestCheckResourceAttr(resourceName, "output.mailEnabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "output.mailNickname", groupName),
					resource.TestCheckResourceAttr(resourceName, "output.mailEnabled", "false"),
				),
			},
			{
				Config: msGraphGroupResourceConfig(groupName, groupName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
			{
				Config: msGraphGroupResourceConfig(updatedGroupName, groupName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "output.displayName", updatedGroupName),
					resource.TestCheckResourceAttr(resourceName, "output.mailNickname", groupName),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
		},
	})
}

func msGraphGroupResourceConfig(displayName string, mailNickname string) string {
	return defaultProviderConfigWith(`
	resource "msgraph_object" "group" {
		collection = "groups"
		properties = {
			displayName = "%s"
			mailEnabled = false
			mailNickname = "%s"
			securityEnabled = true
		}
	}
	`, displayName, mailNickname)
}
