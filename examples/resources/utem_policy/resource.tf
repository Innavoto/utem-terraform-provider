resource "utem_policy" "require_mfa" {
  name            = "Require MFA for Admin Users"
  description     = "All admin users must have MFA enabled"
  category        = "cloud"
  severity        = "critical"
  is_enabled      = true
  resource_type   = "iam_user"
  customer_facing = true
  rego_policy     = <<-EOT
    package cloud.iam.mfa

    deny[msg] {
      input.resource_type == "iam_user"
      input.admin == true
      not input.mfa_enabled
      msg := sprintf("Admin user %s does not have MFA enabled", [input.name])
    }
  EOT
}
