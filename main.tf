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

output "provider_config" {
  value = data.msgraph_provider_config.this
}

output "me" {
    value = data.msgraph_object.me.output
}

output "organization" {
    value = data.msgraph_object.organization.output
}