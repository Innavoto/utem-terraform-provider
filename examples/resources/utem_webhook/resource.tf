resource "utem_webhook" "ci_pipeline" {
  name       = "CI/CD Pipeline Hook"
  url        = "https://ci.innavoto.com/hooks/utem-scan-complete"
  is_enabled = true
  secret     = var.webhook_secret
  events     = ["scan.completed", "finding.created", "finding.resolved"]
}
