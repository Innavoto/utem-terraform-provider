package resources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWebhookResource_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the webhook
			{
				Config: testAccProviderConfig + `
resource "utem_webhook" "test" {
  name       = "tf-acc-test-findings-hook"
  url        = "https://hooks.example.com/utem/findings"
  is_enabled = true
  secret     = "whsec_test_secret_value"
  events     = ["finding.created", "scan.completed"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_webhook.test", "name", "tf-acc-test-findings-hook"),
					resource.TestCheckResourceAttr("utem_webhook.test", "url", "https://hooks.example.com/utem/findings"),
					resource.TestCheckResourceAttr("utem_webhook.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("utem_webhook.test", "events.#", "2"),
					resource.TestCheckResourceAttr("utem_webhook.test", "events.0", "finding.created"),
					resource.TestCheckResourceAttr("utem_webhook.test", "events.1", "scan.completed"),
					resource.TestCheckResourceAttrSet("utem_webhook.test", "id"),
				),
			},
			// Step 2: Update the webhook (change name + toggle is_enabled)
			{
				Config: testAccProviderConfig + `
resource "utem_webhook" "test" {
  name       = "tf-acc-test-findings-hook-updated"
  url        = "https://hooks.example.com/utem/v2/findings"
  is_enabled = false
  secret     = "whsec_updated_secret"
  events     = ["finding.created", "finding.resolved", "scan.completed"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_webhook.test", "name", "tf-acc-test-findings-hook-updated"),
					resource.TestCheckResourceAttr("utem_webhook.test", "url", "https://hooks.example.com/utem/v2/findings"),
					resource.TestCheckResourceAttr("utem_webhook.test", "is_enabled", "false"),
					resource.TestCheckResourceAttr("utem_webhook.test", "events.#", "3"),
					resource.TestCheckResourceAttrSet("utem_webhook.test", "id"),
				),
			},
		},
	})
}

func TestAccWebhookResource_minimal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + `
resource "utem_webhook" "minimal_test" {
  name       = "tf-acc-test-minimal-hook"
  url        = "https://hooks.example.com/minimal"
  is_enabled = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_webhook.minimal_test", "name", "tf-acc-test-minimal-hook"),
					resource.TestCheckResourceAttr("utem_webhook.minimal_test", "url", "https://hooks.example.com/minimal"),
					resource.TestCheckResourceAttr("utem_webhook.minimal_test", "is_enabled", "true"),
					resource.TestCheckResourceAttrSet("utem_webhook.minimal_test", "id"),
				),
			},
		},
	})
}
