package msgraph

import (
	msgraphprovider "github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/provider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var dataSources = []func() datasource.DataSource{
	NewMsGraphProviderConfigDataSource,
	NewMsGraphObjectDataSource,
}

var resources = []func() resource.Resource{
	NewMsGraphObjectResource,
}

func NewProvider() provider.Provider {
	return msgraphprovider.New(dataSources, resources)
}
