package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-msgraph/internal/client"
	"terraform-provider-msgraph/internal/dynamic"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MsGraphObjectDataSourceModel struct {
	ID         types.String  `tfsdk:"id"`
	ApiVersion types.String  `tfsdk:"api_version"`
	Output     types.Dynamic `tfsdk:"output"`
}

type MsGraphObjectDataSource struct {
	ProviderData *client.MsGraphClient
}

var (
	_ datasource.DataSource              = &MsGraphObjectDataSource{}
	_ datasource.DataSourceWithConfigure = &MsGraphObjectDataSource{}
)

func NewMsGraphObjectDataSource() datasource.DataSource {
	return &MsGraphObjectDataSource{}
}

func (r *MsGraphObjectDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*client.MsGraphClient); ok {
		r.ProviderData = v
	}
}

func (r *MsGraphObjectDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_object"
}

func (r *MsGraphObjectDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This data source provides access to Microsoft Graph objects.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the object to retrieve.",
			},

			"api_version": schema.StringAttribute{
				Optional:    true,
				Description: "Override the provider Microsoft Graph API version.",
			},

			"output": schema.DynamicAttribute{
				Computed:    true,
				Description: "The object retrieved from Microsoft Graph.",
			},
		},
	}
}

func (r *MsGraphObjectDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model MsGraphObjectDataSourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)

	if response.Diagnostics.HasError() {
		return
	}

	client := r.ProviderData
	id := model.ID.ValueString()

	httpRequest := client.Request(id)

	if !model.ApiVersion.IsNull() {
		httpRequest = httpRequest.ApiVersion(model.ApiVersion.ValueString())
	}

	result, err := httpRequest.Get(ctx)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Failed to get resource with ID %q.", id), err.Error())
		return
	}

	jsonStr, err := json.Marshal(result)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Failed to marshal resource with ID %q.", id), err.Error())
		return
	}

	output, err := dynamic.FromJSONImplied(jsonStr)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Failed to unmarshal resource with ID %q.", id), err.Error())
		return
	}

	model.Output = types.DynamicValue(output)

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
