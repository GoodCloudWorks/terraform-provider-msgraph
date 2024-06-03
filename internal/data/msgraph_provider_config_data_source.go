package data

import (
	"context"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/client"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MsGraphProviderConfigDataSourceModel struct {
	TenantID types.String `tfsdk:"tenant_id"`
	ClientID types.String `tfsdk:"client_id"`
	ObjectID types.String `tfsdk:"object_id"`
}

type MsGraphProviderConfigDataSource struct {
	Client client.MsGraphClient
}

var (
	_ datasource.DataSource              = &MsGraphProviderConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &MsGraphProviderConfigDataSource{}
)

func NewMsGraphProviderConfigDataSource() datasource.DataSource {
	return &MsGraphProviderConfigDataSource{}
}

func (r *MsGraphProviderConfigDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(client.MsGraphClient); ok {
		r.Client = v
	}
}

func (r *MsGraphProviderConfigDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_provider_config"
}

func (r *MsGraphProviderConfigDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
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

func (r *MsGraphProviderConfigDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model MsGraphProviderConfigDataSourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)

	if response.Diagnostics.HasError() {
		return
	}

	client := r.Client

	accessToken, err := client.GetToken(ctx)
	if err != nil {
		response.Diagnostics.AddError("Failed to get access token.", err.Error())
		return
	}

	var claims jwt.MapClaims
	_, err = jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		claims = token.Claims.(jwt.MapClaims)
		return nil, nil
	}, jwt.WithoutClaimsValidation())
	if claims == nil && err != nil {
		response.Diagnostics.AddError("Failed to parse access token.", err.Error())
		return
	}

	oid := claims["oid"].(string)
	tid := claims["tid"].(string)
	clientID := claims["appid"].(string)

	model.TenantID = types.StringValue(tid)
	model.ObjectID = types.StringValue(oid)
	model.ClientID = types.StringValue(clientID)

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
