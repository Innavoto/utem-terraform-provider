package resources

import (
	"context"
	"fmt"

	"github.com/Innavoto/terraform-provider-utem/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &scanScheduleResource{}
	_ resource.ResourceWithConfigure = &scanScheduleResource{}
)

type scanScheduleResource struct {
	client *client.Client
}

type scanScheduleResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	TargetHost     types.String `tfsdk:"target_host"`
	ScanType       types.String `tfsdk:"scan_type"`
	CronExpression types.String `tfsdk:"cron_expression"`
	IsEnabled      types.Bool   `tfsdk:"is_enabled"`
	Modules        types.List   `tfsdk:"modules"`
}

func NewScanScheduleResource() resource.Resource {
	return &scanScheduleResource{}
}

func (r *scanScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scan_schedule"
}

func (r *scanScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a UTEM scan schedule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Scan schedule ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Schedule name.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Schedule description.",
				Optional:    true,
			},
			"target_host": schema.StringAttribute{
				Description: "Target host to scan.",
				Required:    true,
			},
			"scan_type": schema.StringAttribute{
				Description: "Type of scan: full, quick, or targeted.",
				Required:    true,
			},
			"cron_expression": schema.StringAttribute{
				Description: "Cron expression for the schedule (e.g. \"0 2 * * *\").",
				Required:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the schedule is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"modules": schema.ListAttribute{
				Description: "List of scan modules to run.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *scanScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *scanScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scanScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := scanScheduleFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateScanSchedule(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating scan schedule", err.Error())
		return
	}

	resp.Diagnostics.Append(mapScanScheduleToState(ctx, result, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scanScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scanScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetScanSchedule(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading scan schedule", err.Error())
		return
	}

	resp.Diagnostics.Append(mapScanScheduleToState(ctx, result, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scanScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan scanScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state scanScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := scanScheduleFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateScanSchedule(ctx, state.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating scan schedule", err.Error())
		return
	}

	resp.Diagnostics.Append(mapScanScheduleToState(ctx, result, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scanScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scanScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteScanSchedule(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting scan schedule", err.Error())
		return
	}
}

func scanScheduleFromPlan(ctx context.Context, plan *scanScheduleResourceModel) (*client.ScanSchedule, diag.Diagnostics) {
	var diags diag.Diagnostics

	input := &client.ScanSchedule{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		TargetHost:  plan.TargetHost.ValueString(),
		ScanType:    plan.ScanType.ValueString(),
		CronExpr:    plan.CronExpression.ValueString(),
		IsEnabled:   plan.IsEnabled.ValueBool(),
	}

	if !plan.Modules.IsNull() && !plan.Modules.IsUnknown() {
		var modules []string
		diags.Append(plan.Modules.ElementsAs(ctx, &modules, false)...)
		input.Modules = modules
	}

	return input, diags
}

func mapScanScheduleToState(ctx context.Context, result *client.ScanSchedule, model *scanScheduleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ID = types.StringValue(result.ID)
	model.Name = types.StringValue(result.Name)
	model.Description = types.StringValue(result.Description)
	model.TargetHost = types.StringValue(result.TargetHost)
	model.ScanType = types.StringValue(result.ScanType)
	model.CronExpression = types.StringValue(result.CronExpr)
	model.IsEnabled = types.BoolValue(result.IsEnabled)

	if result.Modules != nil {
		listVal, d := types.ListValueFrom(ctx, types.StringType, result.Modules)
		diags.Append(d...)
		model.Modules = listVal
	}

	return diags
}
