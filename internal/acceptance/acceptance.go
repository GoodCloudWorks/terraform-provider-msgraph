package acceptance

import (
	"fmt"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func Config(config string) string {
	return fmt.Sprintf(`
	provider "msgraph" {
	}

	%s
	`, config)
}

func NewProtoV6ProviderFactories(dataSources []func() datasource.DataSource, resources []func() resource.Resource) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"msgraph": providerserver.NewProtocol6WithError(provider.New(dataSources, resources)),
	}
}
