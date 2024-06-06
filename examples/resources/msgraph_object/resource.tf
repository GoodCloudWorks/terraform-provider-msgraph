resource "msgraph_object" "group" {
  collection = "groups"
  properties = {
    "displayName"     = "My Group"
    "mailEnabled"     = false
    "mailNickname"    = "mygroup"
    "securityEnabled" = true
  }
}
