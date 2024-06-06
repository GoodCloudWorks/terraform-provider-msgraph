package provider

import (
	"context"

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
	dataSources []func() datasource.DataSource
	resources   []func() resource.Resource
}

func New(dataSources []func() datasource.DataSource, resources []func() resource.Resource) provider.Provider {
	return &MsGraphProvider{
		dataSources: dataSources,
		resources:   resources,
	}
}

func (*MsGraphProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "msgraph"
}

func (*MsGraphProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
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

			"scopes": schema.SetAttribute{
				Description: "The scopes to request when authenticating.",
				Optional:    true,
				ElementType: types.StringType,
			},

			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Client ID used for authentication.",
			},

			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Tenant ID to authenticate against.",
			},

			"use_oidc": schema.BoolAttribute{
				Optional:    true,
				Description: "Attempt to use OpenID Connect Federated authentication.",
			},

			"use_msi": schema.BoolAttribute{
				Optional:    true,
				Description: "Attempt to use Managed Service Identity authentication.",
			},

			"use_cli": schema.BoolAttribute{
				Optional:    true,
				Description: "Attempt to use Azure CLI for authentication.",
			},

			"oidc_request_token": schema.StringAttribute{
				Optional:    true,
				Description: "The bearer token for the request to the OIDC provider. For use When authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_request_url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL for the OIDC provider from which to request an ID token. For use When authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token": schema.StringAttribute{
				Optional:    true,
				Description: "The OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token_file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to a file containing an OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},
		},
	}
}

func (*MsGraphProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MsGraphProviderData

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Configure()...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := data.NewClient()
	if err != nil {
		resp.Diagnostics.AddError("Failed to obtain a client.", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (provider *MsGraphProvider) Resources(ctx context.Context) []func() resource.Resource {
	return provider.resources
}

func (provider *MsGraphProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return provider.dataSources
}
