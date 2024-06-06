package msgraph

import (
	"fmt"

	msgraphprovider "github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var protoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"msgraph": providerserver.NewProtocol6WithError(msgraphprovider.New(dataSources, resources)),
}

func defaultProviderConfig() string {
	return `
		provider "msgraph" {}
		`
}

func defaultProviderConfigWith(config string, params ...any) string {
	return fmt.Sprintf(`
		%s
		%s
		`, defaultProviderConfig(), fmt.Sprintf(config, params...))
}
