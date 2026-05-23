resource "utem_notification_rule" "critical_alerts" {
  name         = "Critical Finding Alerts"
  description  = "Route critical and high findings to Slack"
  is_enabled   = true
  severities   = ["critical", "high"]
  categories   = ["vulnerability", "misconfiguration"]
  channel_type = "slack"
}
