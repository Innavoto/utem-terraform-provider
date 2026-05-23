resource "utem_integration" "slack" {
  name             = "Security Alerts Channel"
  integration_type = "slack"
  is_enabled       = true
  webhook_url      = var.slack_webhook_url  # Never hardcode — use a variable or env var
  description      = "Route critical findings to #security-alerts"
}

resource "utem_integration" "jira" {
  name             = "Jira Ticket Creation"
  integration_type = "jira"
  is_enabled       = true
  description      = "Auto-create Jira tickets for high-severity findings"
  config = {
    project_key = "SEC"
    issue_type  = "Bug"
    priority    = "High"
  }
}
