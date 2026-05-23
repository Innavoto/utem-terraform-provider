resource "utem_scan_schedule" "nightly_full" {
  name            = "Nightly Full Scan"
  description     = "Comprehensive scan of all external assets"
  target_host     = "*.innavoto.com"
  scan_type       = "full"
  cron_expression = "0 2 * * *"
  is_enabled      = true
  modules         = ["subdomain", "port_scan", "vuln_scan", "ssl_check", "waf_detect"]
}

resource "utem_scan_schedule" "hourly_quick" {
  name            = "Hourly Quick Check"
  description     = "Fast port and SSL check every hour"
  target_host     = "api.innavoto.com"
  scan_type       = "quick"
  cron_expression = "0 * * * *"
  is_enabled      = true
}
