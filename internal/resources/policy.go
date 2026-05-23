package resources

import (
	"context"
	"fmt"

	"github.com/Innavoto/utem-terraform-provider/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &policyResource{}
	_ resource.ResourceWithConfigure = &policyResource{}
)

type policyResource struct {
	client *client.Client
}

type policyResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Category       types.String `tfsdk:"category"`
	Severity       types.String `tfsdk:"severity"`
	IsEnabled      types.Bool   `tfsdk:"is_enabled"`
	ResourceType   types.String `tfsdk:"resource_type"`
	RegoPolicy     types.String `tfsdk:"rego_policy"`
	CustomerFacing types.Bool   `tfsdk:"customer_facing"`
}

func NewPolicyResource() resource.Resource {
	return &policyResource{}
}

func (r *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (r *policyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a UTEM security policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Policy ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Policy name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Policy description.",
				Optional:    true,
			},
			"category": schema.StringAttribute{
				Description: "Policy category: cloud, containers, ai-agents, compliance, or private-cloud.",
				Required:    true,
			},
			"severity": schema.StringAttribute{
				Description: "Severity level: critical, high, medium, low, or info.",
				Required:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the policy is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"resource_type": schema.StringAttribute{
				Description: "Resource type this policy targets.",
				Optional:    true,
			},
			"rego_policy": schema.StringAttribute{
				Description: "Rego policy source code.",
				Optional:    true,
			},
			"customer_facing": schema.BoolAttribute{
				Description: "Whether this policy is visible to customers.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *policyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.Policy{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Category:       plan.Category.ValueString(),
		Severity:       plan.Severity.ValueString(),
		IsEnabled:      plan.IsEnabled.ValueBool(),
		ResourceType:   plan.ResourceType.ValueString(),
		RegoPolicy:     plan.RegoPolicy.ValueString(),
		CustomerFacing: plan.CustomerFacing.ValueBool(),
	}

	result, err := r.client.CreatePolicy(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating policy", err.Error())
		return
	}

	mapPolicyToState(result, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetPolicy(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading policy", err.Error())
		return
	}

	mapPolicyToState(result, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan policyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state policyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.Policy{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Category:       plan.Category.ValueString(),
		Severity:       plan.Severity.ValueString(),
		IsEnabled:      plan.IsEnabled.ValueBool(),
		ResourceType:   plan.ResourceType.ValueString(),
		RegoPolicy:     plan.RegoPolicy.ValueString(),
		CustomerFacing: plan.CustomerFacing.ValueBool(),
	}

	result, err := r.client.UpdatePolicy(ctx, state.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating policy", err.Error())
		return
	}

	mapPolicyToState(result, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePolicy(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting policy", err.Error())
		return
	}
}

func mapPolicyToState(result *client.Policy, model *policyResourceModel) {
	model.ID = types.StringValue(result.ID)
	model.Name = types.StringValue(result.Name)
	model.Description = types.StringValue(result.Description)
	model.Category = types.StringValue(result.Category)
	model.Severity = types.StringValue(result.Severity)
	model.IsEnabled = types.BoolValue(result.IsEnabled)
	model.ResourceType = types.StringValue(result.ResourceType)
	model.RegoPolicy = types.StringValue(result.RegoPolicy)
	model.CustomerFacing = types.BoolValue(result.CustomerFacing)
}
