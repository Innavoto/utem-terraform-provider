package resources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationRuleResource_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the notification rule
			{
				Config: testAccProviderConfig + `
resource "utem_notification_rule" "test" {
  name         = "tf-acc-test-critical-alerts"
  description  = "Alert on critical findings"
  is_enabled   = true
  channel_type = "slack"
  channel_id   = "C0123456789"
  severities   = ["critical", "high"]
  categories   = ["vulnerability", "misconfiguration"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_notification_rule.test", "name", "tf-acc-test-critical-alerts"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "description", "Alert on critical findings"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "channel_type", "slack"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "channel_id", "C0123456789"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "severities.#", "2"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "severities.0", "critical"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "severities.1", "high"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "categories.#", "2"),
					resource.TestCheckResourceAttrSet("utem_notification_rule.test", "id"),
				),
			},
			// Step 2: Update the notification rule (change name + toggle is_enabled)
			{
				Config: testAccProviderConfig + `
resource "utem_notification_rule" "test" {
  name         = "tf-acc-test-critical-alerts-updated"
  description  = "Updated alert rule"
  is_enabled   = false
  channel_type = "email"
  channel_id   = "security@innavoto.com"
  severities   = ["critical"]
  categories   = ["vulnerability"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_notification_rule.test", "name", "tf-acc-test-critical-alerts-updated"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "is_enabled", "false"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "channel_type", "email"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "channel_id", "security@innavoto.com"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "severities.#", "1"),
					resource.TestCheckResourceAttr("utem_notification_rule.test", "severities.0", "critical"),
					resource.TestCheckResourceAttrSet("utem_notification_rule.test", "id"),
				),
			},
		},
	})
}

func TestAccNotificationRuleResource_minimal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + `
resource "utem_notification_rule" "minimal_test" {
  name         = "tf-acc-test-minimal-rule"
  channel_type = "webhook"
  is_enabled   = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_notification_rule.minimal_test", "name", "tf-acc-test-minimal-rule"),
					resource.TestCheckResourceAttr("utem_notification_rule.minimal_test", "channel_type", "webhook"),
					resource.TestCheckResourceAttr("utem_notification_rule.minimal_test", "is_enabled", "true"),
					resource.TestCheckResourceAttrSet("utem_notification_rule.minimal_test", "id"),
				),
			},
		},
	})
}
