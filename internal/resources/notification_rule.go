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
	_ resource.Resource              = &notificationRuleResource{}
	_ resource.ResourceWithConfigure = &notificationRuleResource{}
)

type notificationRuleResource struct {
	client *client.Client
}

type notificationRuleResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsEnabled   types.Bool   `tfsdk:"is_enabled"`
	Severities  types.List   `tfsdk:"severities"`
	Categories  types.List   `tfsdk:"categories"`
	ChannelType types.String `tfsdk:"channel_type"`
	ChannelID   types.String `tfsdk:"channel_id"`
}

func NewNotificationRuleResource() resource.Resource {
	return &notificationRuleResource{}
}

func (r *notificationRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rule"
}

func (r *notificationRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a UTEM notification rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Notification rule ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Rule name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Rule description.",
				Optional:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the rule is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"severities": schema.ListAttribute{
				Description: "Severity levels to match (e.g. critical, high).",
				Optional:    true,
				ElementType: types.StringType,
			},
			"categories": schema.ListAttribute{
				Description: "Categories to match.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"channel_type": schema.StringAttribute{
				Description: "Notification channel type: slack, email, webhook, or discord.",
				Required:    true,
			},
			"channel_id": schema.StringAttribute{
				Description: "Channel identifier (e.g. Slack channel ID, email address).",
				Optional:    true,
			},
		},
	}
}

func (r *notificationRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *notificationRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := notificationRuleFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateNotificationRule(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating notification rule", err.Error())
		return
	}

	resp.Diagnostics.Append(mapNotificationRuleToState(ctx, result, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *notificationRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetNotificationRule(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading notification rule", err.Error())
		return
	}

	resp.Diagnostics.Append(mapNotificationRuleToState(ctx, result, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *notificationRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state notificationRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := notificationRuleFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(state.ID.ValueInt64())
	result, err := r.client.UpdateNotificationRule(ctx, id, input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating notification rule", err.Error())
		return
	}

	resp.Diagnostics.Append(mapNotificationRuleToState(ctx, result, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *notificationRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNotificationRule(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting notification rule", err.Error())
		return
	}
}

func notificationRuleFromPlan(ctx context.Context, plan *notificationRuleResourceModel) (*client.NotificationRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	input := &client.NotificationRule{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsEnabled:   plan.IsEnabled.ValueBool(),
		ChannelType: plan.ChannelType.ValueString(),
		ChannelID:   plan.ChannelID.ValueString(),
	}

	if !plan.Severities.IsNull() && !plan.Severities.IsUnknown() {
		var severities []string
		diags.Append(plan.Severities.ElementsAs(ctx, &severities, false)...)
		input.Severities = severities
	}

	if !plan.Categories.IsNull() && !plan.Categories.IsUnknown() {
		var categories []string
		diags.Append(plan.Categories.ElementsAs(ctx, &categories, false)...)
		input.Categories = categories
	}

	return input, diags
}

func mapNotificationRuleToState(ctx context.Context, result *client.NotificationRule, model *notificationRuleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.Int64Value(int64(result.ID))
	model.Name = types.StringValue(result.Name)
	model.Description = types.StringValue(result.Description)
	model.IsEnabled = types.BoolValue(result.IsEnabled)
	model.ChannelType = types.StringValue(result.ChannelType)
	model.ChannelID = types.StringValue(result.ChannelID)

	if result.Severities != nil {
		listVal, d := types.ListValueFrom(ctx, types.StringType, result.Severities)
		diags.Append(d...)
		model.Severities = listVal
	}

	if result.Categories != nil {
		listVal, d := types.ListValueFrom(ctx, types.StringType, result.Categories)
		diags.Append(d...)
		model.Categories = listVal
	}

	return diags
}
