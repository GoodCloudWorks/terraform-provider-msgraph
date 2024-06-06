package msgraph

import (
	"context"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &msGraphObjectDataSource{}
	_ datasource.DataSourceWithConfigure = &msGraphObjectDataSource{}
)

type msGraphObjectDataSource struct {
	client client.MsGraphClient
}

type msGraphObjectDataSourceModel struct {
	ID         types.String  `tfsdk:"id"`
	ApiVersion types.String  `tfsdk:"api_version"`
	Collection types.String  `tfsdk:"collection"`
	Output     types.Dynamic `tfsdk:"output"`
}

func NewMsGraphObjectDataSource() datasource.DataSource {
	return &msGraphObjectDataSource{}
}

func (r *msGraphObjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if v, ok := req.ProviderData.(client.MsGraphClient); ok {
		r.client = v
	}
}

func (r *msGraphObjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (r *msGraphObjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides access to Microsoft Graph objects.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the object to retrieve.",
			},

			"collection": schema.StringAttribute{
				Computed:    true,
				Description: "The collection of the object to retrieve.",
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

func (r *msGraphObjectDataSource) Read(ctx context.Context, request datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model msGraphObjectDataSourceModel
	resp.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := ensureParseIDString(model.ID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	http := r.client.R(ctx, model.ApiVersion)

	content, diags := ensureGetObjectAsDynamic(http, id.Path)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	model.Output = content
	model.Collection = types.StringValue(id.Collection())

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
