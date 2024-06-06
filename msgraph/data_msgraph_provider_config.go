package msgraph

import (
	"context"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/client"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &msGraphProviderConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &msGraphProviderConfigDataSource{}
)

type msGraphProviderConfigDataSource struct {
	client client.MsGraphClient
}

type msGraphProviderConfigDataSourceModel struct {
	TenantID types.String `tfsdk:"tenant_id"`
	ClientID types.String `tfsdk:"client_id"`
	ObjectID types.String `tfsdk:"object_id"`
}

func NewMsGraphProviderConfigDataSource() datasource.DataSource {
	return &msGraphProviderConfigDataSource{}
}

func (r *msGraphProviderConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if v, ok := req.ProviderData.(client.MsGraphClient); ok {
		r.client = v
	}
}

func (r *msGraphProviderConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_config"
}

func (r *msGraphProviderConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides access to Microsoft Graph objects.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "The client ID (application ID) linked to the authenticated principal, or the application used for delegated authentication",
			},

			"client_id": schema.StringAttribute{
				Computed:    true,
				Description: "The client ID (application ID) linked to the authenticated principal, or the application used for delegated authentication",
			},

			"object_id": schema.StringAttribute{
				Computed:    true,
				Description: "The object ID of the authenticated principal",
			},
		},
	}
}

func (r *msGraphProviderConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model msGraphProviderConfigDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessToken, err := r.client.GetToken(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get access token.", err.Error())
		return
	}

	var claims jwt.MapClaims
	_, err = jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		claims = token.Claims.(jwt.MapClaims)
		return nil, nil
	}, jwt.WithoutClaimsValidation())
	if claims == nil && err != nil {
		resp.Diagnostics.AddError("Failed to parse access token.", err.Error())
		return
	}

	oid := claims["oid"].(string)
	tid := claims["tid"].(string)
	clientID := claims["appid"].(string)

	model.TenantID = types.StringValue(tid)
	model.ObjectID = types.StringValue(oid)
	model.ClientID = types.StringValue(clientID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
