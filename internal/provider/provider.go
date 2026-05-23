package provider

import (
	"context"
	"os"

	"github.com/Innavoto/utem-terraform-provider/internal/client"
	"github.com/Innavoto/utem-terraform-provider/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &utemProvider{}

type utemProvider struct {
	version string
}

type utemProviderModel struct {
	BaseURL  types.String `tfsdk:"base_url"`
	APIKey   types.String `tfsdk:"api_key"`
	TenantID types.String `tfsdk:"tenant_id"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &utemProvider{version: version}
	}
}

func (p *utemProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utem"
	resp.Version = p.version
}

func (p *utemProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage UTEM platform resources — integrations, scan schedules, policies, notification rules, and webhooks.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "UTEM API base URL. Defaults to UTEM_BASE_URL env var.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "UTEM API key. Defaults to UTEM_API_KEY env var.",
				Optional:    true,
				Sensitive:   true,
			},
			"tenant_id": schema.StringAttribute{
				Description: "Tenant ID. Defaults to UTEM_TENANT_ID env var.",
				Optional:    true,
			},
		},
	}
}

func (p *utemProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config utemProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseURL := os.Getenv("UTEM_BASE_URL")
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}
	if baseURL == "" {
		baseURL = "https://utem.innavoto.com"
	}

	apiKey := os.Getenv("UTEM_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	tenantID := os.Getenv("UTEM_TENANT_ID")
	if !config.TenantID.IsNull() {
		tenantID = config.TenantID.ValueString()
	}
	if tenantID == "" {
		tenantID = "1"
	}

	c := client.New(baseURL, apiKey, tenantID)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *utemProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewIntegrationResource,
		resources.NewScanScheduleResource,
		resources.NewPolicyResource,
		resources.NewNotificationRuleResource,
		resources.NewWebhookResource,
	}
}

func (p *utemProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
