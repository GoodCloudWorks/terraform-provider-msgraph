terraform {
    required_providers {
      msgraph = {
        source = "msgraph"
      }
    }
}

provider "msgraph" {}

data "msgraph_object" "me" {
    id = "me"
}

output "me" {
    value = data.msgraph_object.me.output
}