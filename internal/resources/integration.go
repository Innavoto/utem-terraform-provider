package resources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Innavoto/terraform-provider-utem/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &integrationResource{}
	_ resource.ResourceWithConfigure = &integrationResource{}
)

type integrationResource struct {
	client *client.Client
}

type integrationResourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	IntegrationType types.String `tfsdk:"integration_type"`
	IsEnabled       types.Bool   `tfsdk:"is_enabled"`
	Config          types.Map    `tfsdk:"config"`
	WebhookURL      types.String `tfsdk:"webhook_url"`
	Description     types.String `tfsdk:"description"`
}

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *integrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a UTEM integration (Slack, Jira, ServiceNow, Discord, Syslog, etc.).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Integration ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Integration name.",
				Required:    true,
			},
			"integration_type": schema.StringAttribute{
				Description: "Type of integration (e.g. slack, jira, servicenow, discord, syslog).",
				Required:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the integration is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"config": schema.MapAttribute{
				Description: "Arbitrary key-value configuration for the integration.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"webhook_url": schema.StringAttribute{
				Description: "Webhook URL for the integration.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the integration.",
				Optional:    true,
			},
		},
	}
}

func (r *integrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan integrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.Integration{
		Name:            plan.Name.ValueString(),
		IntegrationType: plan.IntegrationType.ValueString(),
		IsEnabled:       plan.IsEnabled.ValueBool(),
		WebhookURL:      plan.WebhookURL.ValueString(),
		Description:     plan.Description.ValueString(),
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configMap := make(map[string]string)
		resp.Diagnostics.Append(plan.Config.ElementsAs(ctx, &configMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Config = make(map[string]interface{}, len(configMap))
		for k, v := range configMap {
			input.Config[k] = v
		}
	}

	result, err := r.client.CreateIntegration(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating integration", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(result.ID))
	plan.Name = types.StringValue(result.Name)
	plan.IntegrationType = types.StringValue(result.IntegrationType)
	plan.IsEnabled = types.BoolValue(result.IsEnabled)
	plan.WebhookURL = types.StringValue(result.WebhookURL)
	plan.Description = types.StringValue(result.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state integrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetIntegration(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading integration", err.Error())
		return
	}

	state.ID = types.Int64Value(int64(result.ID))
	state.Name = types.StringValue(result.Name)
	state.IntegrationType = types.StringValue(result.IntegrationType)
	state.IsEnabled = types.BoolValue(result.IsEnabled)
	state.WebhookURL = types.StringValue(result.WebhookURL)
	state.Description = types.StringValue(result.Description)

	if result.Config != nil {
		configMap := make(map[string]string, len(result.Config))
		for k, v := range result.Config {
			configMap[k] = fmt.Sprintf("%v", v)
		}
		mapVal, diags := types.MapValueFrom(ctx, types.StringType, configMap)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Config = mapVal
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan integrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state integrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &client.Integration{
		Name:            plan.Name.ValueString(),
		IntegrationType: plan.IntegrationType.ValueString(),
		IsEnabled:       plan.IsEnabled.ValueBool(),
		WebhookURL:      plan.WebhookURL.ValueString(),
		Description:     plan.Description.ValueString(),
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configMap := make(map[string]string)
		resp.Diagnostics.Append(plan.Config.ElementsAs(ctx, &configMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Config = make(map[string]interface{}, len(configMap))
		for k, v := range configMap {
			input.Config[k] = v
		}
	}

	id := int(state.ID.ValueInt64())
	result, err := r.client.UpdateIntegration(ctx, id, input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating integration", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(result.ID))
	plan.Name = types.StringValue(result.Name)
	plan.IntegrationType = types.StringValue(result.IntegrationType)
	plan.IsEnabled = types.BoolValue(result.IsEnabled)
	plan.WebhookURL = types.StringValue(result.WebhookURL)
	plan.Description = types.StringValue(result.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIntegration(ctx, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting integration", err.Error())
		return
	}
}

// importStateID parses the integration ID from the import string.
func integrationImportID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid integration ID %q: %w", s, err)
	}
	return id, nil
}
