package msgraph

import (
	"context"

	"github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/client"
	"github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/dynamic"
	"github.com/GoodCloudWorks/terraform-provider-msgraph/msgraph/id"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &msGraphObjectResource{}
	_ resource.ResourceWithConfigure   = &msGraphObjectResource{}
	_ resource.ResourceWithImportState = &msGraphObjectResource{}
)

type msGraphObjectResource struct {
	client client.MsGraphClient
}

type msGraphObjectResourceModel struct {
	Collection types.String  `tfsdk:"collection"`
	ID         types.String  `tfsdk:"id"`
	ApiVersion types.String  `tfsdk:"api_version"`
	Properties types.Dynamic `tfsdk:"properties"`
	Output     types.Dynamic `tfsdk:"output"`
}

func NewMsGraphObjectResource() resource.Resource {
	return &msGraphObjectResource{}
}

func (r *msGraphObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if v, ok := req.ProviderData.(client.MsGraphClient); ok {
		r.client = v
	}
}

func (*msGraphObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (*msGraphObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource provides the ability to create an object in a Microsoft Graph collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the object.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"collection": schema.StringAttribute{
				Required:    true,
				Description: "The collection of the object to retrieve.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"api_version": schema.StringAttribute{
				Optional:    true,
				Description: "Override the provider Microsoft Graph API version.",
			},

			"properties": schema.DynamicAttribute{
				Required:    true,
				Description: "The properties of the object.",
				PlanModifiers: []planmodifier.Dynamic{
					dynamic.UseStateWhen(dynamic.SemanticallyEqual),
				},
			},

			"output": schema.DynamicAttribute{
				Computed:    true,
				Description: "The object retrieved from Microsoft Graph.",
			},
		},
	}
}

func (r *msGraphObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model msGraphObjectResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	path, diags := ensureIsValidPathString(model.Collection)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	http := r.client.R(ctx, model.ApiVersion)

	if diags := ensureRequestSetBodyFromDynamic(http, model.Properties); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	response, err := post(http, path)
	if diags := ensureHttpResponseSucceeded(response, err); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	objectID, diags := ensureResponseHasObjectID(response)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	id := id.New(path, objectID)
	model.ID = id.AsString()

	content, diags := ensureGetObjectAsDynamic(http, id.Path)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	model.Output = content

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *msGraphObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model msGraphObjectResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
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

	properties, err := dynamic.Apply(content, model.Properties)
	if err != nil {
		resp.Diagnostics.Append(errorDiagnostics("Failed to apply dynamic properties.", err.Error())...)
		return
	}
	model.Properties = properties

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *msGraphObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model msGraphObjectResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := ensureParseIDString(model.ID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	http := r.client.R(ctx, model.ApiVersion)

	if diags := ensureRequestSetBodyFromDynamic(http, model.Properties); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	response, err := patch(http, id.Path)
	if diags := ensureHttpResponseSucceeded(response, err); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	content, diags := ensureGetObjectAsDynamic(http, id.Path)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	model.Output = content

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *msGraphObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model msGraphObjectResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := ensureParseIDString(model.ID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	http := r.client.R(ctx, model.ApiVersion)

	response, err := delete(http, id.Path)
	if response.StatusCode() == httpStatusNotFound {
		return
	}

	if diags := ensureHttpResponseSucceeded(response, err); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (r *msGraphObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, diags := ensureParseID(req.ID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	model := msGraphObjectResourceModel{
		ID:         id.AsString(),
		Collection: types.StringValue(id.Collection()),
		ApiVersion: types.StringNull(),
	}

	if apiVersion := id.ApiVersion(); apiVersion != "" {
		model.ApiVersion = types.StringValue(apiVersion)
	}

	http := r.client.R(ctx, model.ApiVersion)

	content, diags := ensureGetObjectAsDynamic(http, id.Path)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	model.Output = content
	model.Properties = content

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
