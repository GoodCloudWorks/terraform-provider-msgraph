package provider

import (
	"context"
	"terraform-provider-msgraph/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &MsGraphProvider{}

type MsGraphProvider struct {
	version string
}

type MsGraphProviderData struct {
	ApiVersion types.String `tfsdk:"api_version"`
}

func (p *MsGraphProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "msgraph"
	resp.Version = p.version
}

func (p *MsGraphProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The provider allows you to interact with Microsoft Graph.",
		Attributes: map[string]schema.Attribute{
			"api_version": schema.StringAttribute{
				Description: "The Microsoft Graph API version to use, default is v1.0.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("v1.0", "beta"),
				},
			},
		},
	}
}

func (p *MsGraphProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MsGraphProviderData

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ApiVersion.IsNull() {
		data.ApiVersion = types.StringValue("v1.0")
	}

	credential, err := client.NewTokenCredential()
	if err != nil {
		resp.Diagnostics.AddError("Failed to obtain a credential.", err.Error())
		return
	}

	client := client.NewMsGraphClient(data.ApiVersion.ValueString(), credential)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MsGraphProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *MsGraphProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMsGraphObjectDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MsGraphProvider{
			version: version,
		}
	}
}
