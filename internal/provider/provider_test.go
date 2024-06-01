package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const testProviderConfig = `
	provider "msgraph" {
	}
	`

var testProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"msgraph": providerserver.NewProtocol6WithError(&MsGraphProvider{}),
}
