package data

import (
	"context"
	"fmt"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/client"
	"github.com/GoodCloudWorks/terraform-provider-msgraph/internal/dynamic"

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
	Client client.MsGraphClient
}

var (
	_ datasource.DataSource              = &MsGraphObjectDataSource{}
	_ datasource.DataSourceWithConfigure = &MsGraphObjectDataSource{}
)

func NewMsGraphObjectDataSource() datasource.DataSource {
	return &MsGraphObjectDataSource{}
}

func (r *MsGraphObjectDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(client.MsGraphClient); ok {
		r.Client = v
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

	client := r.Client
	http := client.R(ctx, model.ApiVersion)
	path := client.URL(model.ID)

	result, err := http.Get(path)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Failed to get resource with ID %q.", path), err.Error())
		return
	}

	if result.IsError() {
		response.Diagnostics.AddError(fmt.Sprintf("Failed (%d) to get resource with ID %q.", result.StatusCode(), path), string(result.Body()))
		return
	}

	output, err := dynamic.FromJSONImplied(result.Body())
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Failed to read resource %q response.", path), err.Error())
		return
	}

	model.Output = types.DynamicValue(output)

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
