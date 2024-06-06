package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MsGraphProviderData struct {
	ApiVersion types.String `tfsdk:"api_version"`
	Scopes     types.Set    `tfsdk:"scopes"`

	TenantID types.String `tfsdk:"tenant_id"`
	ClientID types.String `tfsdk:"client_id"`

	UseOIDC types.Bool `tfsdk:"use_oidc"`
	UseMSI  types.Bool `tfsdk:"use_msi"`
	UseCLI  types.Bool `tfsdk:"use_cli"`

	OIDCRequestToken  types.String `tfsdk:"oidc_request_token"`
	OIDCRequestURL    types.String `tfsdk:"oidc_request_url"`
	OIDCToken         types.String `tfsdk:"oidc_token"`
	OIDCTokenFilePath types.String `tfsdk:"oidc_token_file_path"`
}

func (data *MsGraphProviderData) Configure() diag.Diagnostics {
	diag := data.read()
	if diag.HasError() {
		return diag
	}

	data.writeEnvironmentVariables()

	return diag
}

func (data *MsGraphProviderData) read() diag.Diagnostics {
	var diag diag.Diagnostics

	data.ApiVersion = readStringFromEnvironment(data.ApiVersion, "MSGRAPH_API_VERSION")
	if data.ApiVersion.IsNull() {
		data.ApiVersion = types.StringValue("v1.0")
	}

	if data.Scopes.IsNull() || len(data.Scopes.Elements()) == 0 {
		data.Scopes, diag = types.SetValue(types.StringType, []attr.Value{types.StringValue("https://graph.microsoft.com/.default")})
		if diag.HasError() {
			return diag
		}
	}

	data.ClientID = readStringFromEnvironment(data.ClientID, "ARM_CLIENT_ID")
	data.TenantID = readStringFromEnvironment(data.TenantID, "ARM_TENANT_ID")

	// OIDC
	data.UseOIDC = readBoolFromEnvironment(data.UseOIDC, "ARM_USE_OIDC")
	data.OIDCRequestToken = readStringFromEnvironment(data.OIDCRequestToken, "ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	data.OIDCRequestURL = readStringFromEnvironment(data.OIDCRequestURL, "ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL")
	data.OIDCToken = readStringFromEnvironment(data.OIDCToken, "ARM_OIDC_TOKEN")
	data.OIDCTokenFilePath = readStringFromEnvironment(data.OIDCTokenFilePath, "ARM_OIDC_TOKEN_FILE_PATH")

	// MSI
	data.UseMSI = readBoolFromEnvironment(data.UseMSI, "ARM_USE_MSI")

	// CLI
	data.UseCLI = defaultIsTrue(readBoolFromEnvironment(data.UseCLI, "ARM_USE_CLI"))

	return diag
}

func (data *MsGraphProviderData) writeEnvironmentVariables() {
	tryWriteEnvironmentVariable("AZURE_TENANT_ID", data.TenantID)
	tryWriteEnvironmentVariable("AZURE_CLIENT_ID", data.ClientID)
}

func tryWriteEnvironmentVariable(name string, value types.String) {
	if v := value.ValueString(); v != "" {
		_ = os.Setenv(name, v)
	}
}

func defaultIsTrue(data types.Bool) types.Bool {
	if data.IsNull() {
		return types.BoolValue(true)
	}
	return data
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
