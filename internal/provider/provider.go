package provider

import (
	"context"
	"os"
	"terraform-provider-msgraph/internal/client"
	"terraform-provider-msgraph/internal/client/credentials"

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

	credentialOptions := &credentials.CredentialOptions{
		TenantID: data.TenantID.ValueString(),
		ClientID: data.ClientID.ValueString(),

		UseOIDC:           data.UseOIDC.ValueBool(),
		OIDCRequestToken:  data.OIDCRequestToken.ValueString(),
		OIDCRequestURL:    data.OIDCRequestURL.ValueString(),
		OIDCToken:         data.OIDCToken.ValueString(),
		OIDCTokenFilePath: data.OIDCTokenFilePath.ValueString(),
	}

	options := &client.MsGraphClientOptions{
		ApiVersion:  data.ApiVersion.ValueString(),
		Credentials: credentialOptions,
	}

	client, err := client.NewMsGraphClient(options)
	if err != nil {
		resp.Diagnostics.AddError("Failed to obtain a msgraph client.", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func readStringFromEnvironment(data types.String, names ...string) types.String {
	if data.IsNull() {
		for _, name := range names {
			if value := os.Getenv(name); value != "" {
				return types.StringValue(value)
			}
		}
	}
	return data
}

func readBoolFromEnvironment(data types.Bool, names ...string) types.Bool {
	if data.IsNull() {
		for _, name := range names {
			if value := os.Getenv(name); value != "" {
				return types.BoolValue(value == "true")
			}
		}
		return types.BoolValue(false)
	}

	return data
}

func readProviderData(data *MsGraphProviderData) {
	data.ApiVersion = readStringFromEnvironment(data.ApiVersion, "MSGRAPH_API_VERSION")
	if data.ApiVersion.IsNull() {
		data.ApiVersion = types.StringValue("v1.0")
	}

	data.ClientID = readStringFromEnvironment(data.ClientID, "ARM_CLIENT_ID")
	data.TenantID = readStringFromEnvironment(data.TenantID, "ARM_TENANT_ID")

	// OIDC
	data.OIDCRequestToken = readStringFromEnvironment(data.OIDCRequestToken, "ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	data.OIDCRequestURL = readStringFromEnvironment(data.OIDCRequestURL, "ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL")
	data.OIDCToken = readStringFromEnvironment(data.OIDCToken, "ARM_OIDC_TOKEN")
	data.OIDCTokenFilePath = readStringFromEnvironment(data.OIDCTokenFilePath, "ARM_OIDC_TOKEN_FILE_PATH")
	data.UseOIDC = readBoolFromEnvironment(data.UseOIDC, "ARM_USE_OIDC")
}

func tryWriteEnvironmentVariable(name string, value types.String) {
	if v := value.ValueString(); v != "" {
		_ = os.Setenv(name, v)
	}
}

func writeDefaultAzureCredentialEnvironmentVariables(data *MsGraphProviderData) {
	tryWriteEnvironmentVariable("AZURE_TENANT_ID", data.TenantID)
	tryWriteEnvironmentVariable("AZURE_CLIENT_ID", data.ClientID)
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
