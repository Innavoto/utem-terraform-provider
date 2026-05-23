# UTEM Terraform Provider

Manage your [UTEM](https://utem.innavoto.com) security platform as infrastructure-as-code.

## Requirements

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22 (for building from source)
- UTEM API key ([Settings → API Tokens](https://utem.innavoto.com/settings/api-tokens))

## Installation

```hcl
terraform {
  required_providers {
    utem = {
      source  = "Innavoto/utem"
      version = "~> 0.1"
    }
  }
}
```

## Authentication

```hcl
provider "utem" {
  base_url  = "https://utem.innavoto.com"  # or UTEM_BASE_URL env var
  api_key   = var.utem_api_key              # or UTEM_API_KEY env var (recommended)
  tenant_id = "1"                           # or UTEM_TENANT_ID env var
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `utem_integration` | Manage integrations (Slack, Jira, ServiceNow, Discord, Syslog) |
| `utem_scan_schedule` | Create and manage automated scan schedules |
| `utem_policy` | Create custom security policies (Rego-based) |
| `utem_notification_rule` | Configure severity-based notification routing |
| `utem_webhook` | Register webhook endpoints for event delivery |

## Quick Start

```hcl
# Configure the provider
provider "utem" {}

# Create a Slack integration
resource "utem_integration" "slack_alerts" {
  name             = "Security Alerts"
  integration_type = "slack"
  is_enabled       = true
  webhook_url      = var.slack_webhook_url
  description      = "Critical and high severity finding alerts"
}

# Schedule a nightly full scan
resource "utem_scan_schedule" "nightly" {
  name            = "Nightly Full Scan"
  target_host     = "*.innavoto.com"
  scan_type       = "full"
  cron_expression = "0 2 * * *"
  is_enabled      = true
  modules         = ["subdomain", "port_scan", "vuln_scan", "ssl_check"]
}

# Route critical findings to Slack
resource "utem_notification_rule" "critical_slack" {
  name         = "Critical to Slack"
  is_enabled   = true
  severities   = ["critical", "high"]
  channel_type = "slack"
  channel_id   = utem_integration.slack_alerts.id
}

# Register a webhook for scan completions
resource "utem_webhook" "ci_webhook" {
  name      = "CI Pipeline Webhook"
  url       = "https://ci.innavoto.com/hooks/utem"
  is_enabled = true
  events    = ["scan.completed", "finding.created"]
}
```

## Building from Source

```bash
git clone https://github.com/Innavoto/terraform-provider-utem.git
cd terraform-provider-utem
go build -o terraform-provider-utem
```

## Development

```bash
# Run tests
go test ./...

# Run acceptance tests (requires UTEM_BASE_URL + UTEM_API_KEY)
TF_ACC=1 go test ./...

# Install locally for testing
go install .
```

## Publisher

Innavoto India Pvt Ltd — [utem.innavoto.com](https://utem.innavoto.com)

## License

MPL-2.0
