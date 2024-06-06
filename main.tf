terraform {
    required_providers {
      msgraph = {
        source = "msgraph"
      }
    }
}

provider "msgraph" {}

data "msgraph_provider_config" "this" {}

data "msgraph_object" "me" {
    id = "me"
}

data "msgraph_object" "organization" {
  id = "organization/${data.msgraph_provider_config.this.tenant_id}" 
}

resource "msgraph_object" "group" {
  collection = "groups"
  properties = {
    displayName = "terraform-provider-msgraph"
    mailEnabled = false
    mailNickname = "terraform-provider-msgraph"
    securityEnabled = true
  }
}

output "provider_config" {
  value = data.msgraph_provider_config.this
}

output "me" {
    value = data.msgraph_object.me.output
}

output "organization" {
    value = data.msgraph_object.organization.output
}

output "group" {
  value = msgraph_object.group.output
}