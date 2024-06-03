package data

import "github.com/hashicorp/terraform-plugin-framework/datasource"

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMsGraphProviderConfigDataSource,
		NewMsGraphObjectDataSource,
	}
}
