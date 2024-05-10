package provider

import (
	"context"
	"os"
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
	ClientID   types.String `tfsdk:"client_id"`
	TenantID   types.String `tfsdk:"tenant_id"`

	UseOIDC           types.Bool   `tfsdk:"use_oidc"`
	OIDCRequestToken  types.String `tfsdk:"oidc_request_token"`
	OIDCRequestURL    types.String `tfsdk:"oidc_request_url"`
	OIDCToken         types.String `tfsdk:"oidc_token"`
	OIDCTokenFilePath types.String `tfsdk:"oidc_token_file_path"`
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

			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Client ID which should be used.",
			},

			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Tenant ID which should be used.",
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

			"use_oidc": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow OpenID Connect to be used for authentication",
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

	readProviderData(&data)
	writeDefaultAzureCredentialEnvironmentVariables(&data)

	options := &client.MsGraphClientOptions{
		ApiVersion: data.ApiVersion.ValueString(),
		TenantID:   data.TenantID.ValueString(),
		ClientID:   data.ClientID.ValueString(),

		UseOIDC:           data.UseOIDC.ValueBool(),
		OIDCRequestToken:  data.OIDCRequestToken.ValueString(),
		OIDCRequestURL:    data.OIDCRequestURL.ValueString(),
		OIDCToken:         data.OIDCToken.ValueString(),
		OIDCTokenFilePath: data.OIDCTokenFilePath.ValueString(),
	}
	client, err := client.NewMsGraphClient(options)
	if err != nil {
		resp.Diagnostics.AddError("Failed to obtain a msgraph client.", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func readProviderData(data *MsGraphProviderData) {
	if data.ApiVersion.IsNull() {
		if v := os.Getenv("MSGRAPH_API_VERSION"); v != "" {
			data.ApiVersion = types.StringValue(v)
		} else {
			data.ApiVersion = types.StringValue("v1.0")
		}
	}

	if data.ClientID.IsNull() {
		if v := os.Getenv("ARM_CLIENT_ID"); v != "" {
			data.ClientID = types.StringValue(v)
		}
	}

	if data.TenantID.IsNull() {
		if v := os.Getenv("ARM_TENANT_ID"); v != "" {
			data.TenantID = types.StringValue(v)
		}
	}

	readOidcOptions(data)
}

func readOidcOptions(data *MsGraphProviderData) {
	if data.OIDCRequestToken.IsNull() {
		if v := os.Getenv("ARM_OIDC_REQUEST_TOKEN"); v != "" {
			data.OIDCRequestToken = types.StringValue(v)
		} else if v := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"); v != "" {
			data.OIDCRequestToken = types.StringValue(v)
		}
	}

	if data.OIDCRequestURL.IsNull() {
		if v := os.Getenv("ARM_OIDC_REQUEST_URL"); v != "" {
			data.OIDCRequestURL = types.StringValue(v)
		} else if v := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"); v != "" {
			data.OIDCRequestURL = types.StringValue(v)
		}
	}

	if data.OIDCToken.IsNull() {
		if v := os.Getenv("ARM_OIDC_TOKEN"); v != "" {
			data.OIDCToken = types.StringValue(v)
		}
	}

	if data.OIDCTokenFilePath.IsNull() {
		if v := os.Getenv("ARM_OIDC_TOKEN_FILE_PATH"); v != "" {
			data.OIDCTokenFilePath = types.StringValue(v)
		}
	}

	if data.UseOIDC.IsNull() {
		if v := os.Getenv("ARM_USE_OIDC"); v != "" {
			data.UseOIDC = types.BoolValue(v == "true")
		} else {
			data.UseOIDC = types.BoolValue(false)
		}
	}
}

func writeDefaultAzureCredentialEnvironmentVariables(data *MsGraphProviderData) {
	if v := data.TenantID.ValueString(); v != "" {
		_ = os.Setenv("AZURE_TENANT_ID", v)
	}
	if v := data.ClientID.ValueString(); v != "" {
		_ = os.Setenv("AZURE_CLIENT_ID", v)
	}
}

func (p *MsGraphProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *MsGraphProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMsGraphProviderConfigDataSource,
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
