package resources

import (
	"context"
	"fmt"

	"github.com/Innavoto/terraform-provider-utem/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &webhookResource{}
	_ resource.ResourceWithConfigure = &webhookResource{}
)

type webhookResource struct {
	client *client.Client
}

type webhookResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	URL       types.String `tfsdk:"url"`
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
	Secret    types.String `tfsdk:"secret"`
	Events    types.List   `tfsdk:"events"`
}

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a UTEM webhook endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Webhook ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Webhook name.",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "Webhook endpoint URL.",
				Required:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the webhook is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"secret": schema.StringAttribute{
				Description: "Webhook signing secret.",
				Optional:    true,
				Sensitive:   true,
			},
			"events": schema.ListAttribute{
				Description: "Event types to deliver (e.g. finding.created, scan.completed).",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := webhookFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateWebhook(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating webhook", err.Error())
		return
	}

	resp.Diagnostics.Append(mapWebhookToState(ctx, result, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetWebhook(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading webhook", err.Error())
		return
	}

	resp.Diagnostics.Append(mapWebhookToState(ctx, result, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := webhookFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(state.ID.ValueInt64())
	result, err := r.client.UpdateWebhook(ctx, id, input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating webhook", err.Error())
		return
	}

	resp.Diagnostics.Append(mapWebhookToState(ctx, result, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWebhook(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting webhook", err.Error())
		return
	}
}

func webhookFromPlan(ctx context.Context, plan *webhookResourceModel) (*client.Webhook, diag.Diagnostics) {
	var diags diag.Diagnostics

	input := &client.Webhook{
		Name:      plan.Name.ValueString(),
		URL:       plan.URL.ValueString(),
		IsEnabled: plan.IsEnabled.ValueBool(),
		Secret:    plan.Secret.ValueString(),
	}

	if !plan.Events.IsNull() && !plan.Events.IsUnknown() {
		var events []string
		diags.Append(plan.Events.ElementsAs(ctx, &events, false)...)
		input.Events = events
	}

	return input, diags
}

func mapWebhookToState(ctx context.Context, result *client.Webhook, model *webhookResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.Int64Value(int64(result.ID))
	model.Name = types.StringValue(result.Name)
	model.URL = types.StringValue(result.URL)
	model.IsEnabled = types.BoolValue(result.IsEnabled)
	model.Secret = types.StringValue(result.Secret)

	if result.Events != nil {
		listVal, d := types.ListValueFrom(ctx, types.StringType, result.Events)
		diags.Append(d...)
		model.Events = listVal
	}

	return diags
}
