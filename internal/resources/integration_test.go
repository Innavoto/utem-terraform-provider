package resources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationResource_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the integration
			{
				Config: testAccProviderConfig + `
resource "utem_integration" "test" {
  name             = "tf-acc-test-slack"
  integration_type = "slack"
  is_enabled       = true
  description      = "Acceptance test integration"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_integration.test", "name", "tf-acc-test-slack"),
					resource.TestCheckResourceAttr("utem_integration.test", "integration_type", "slack"),
					resource.TestCheckResourceAttr("utem_integration.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("utem_integration.test", "description", "Acceptance test integration"),
					resource.TestCheckResourceAttrSet("utem_integration.test", "id"),
				),
			},
			// Step 2: Update the integration (change name + toggle is_enabled)
			{
				Config: testAccProviderConfig + `
resource "utem_integration" "test" {
  name             = "tf-acc-test-slack-updated"
  integration_type = "slack"
  is_enabled       = false
  description      = "Updated acceptance test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_integration.test", "name", "tf-acc-test-slack-updated"),
					resource.TestCheckResourceAttr("utem_integration.test", "integration_type", "slack"),
					resource.TestCheckResourceAttr("utem_integration.test", "is_enabled", "false"),
					resource.TestCheckResourceAttr("utem_integration.test", "description", "Updated acceptance test"),
					resource.TestCheckResourceAttrSet("utem_integration.test", "id"),
				),
			},
		},
	})
}

func TestAccIntegrationResource_withWebhookURL(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + `
resource "utem_integration" "webhook_test" {
  name             = "tf-acc-test-discord"
  integration_type = "discord"
  is_enabled       = true
  webhook_url      = "https://discord.com/api/webhooks/test/token"
  description      = "Discord integration test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_integration.webhook_test", "name", "tf-acc-test-discord"),
					resource.TestCheckResourceAttr("utem_integration.webhook_test", "integration_type", "discord"),
					resource.TestCheckResourceAttr("utem_integration.webhook_test", "webhook_url", "https://discord.com/api/webhooks/test/token"),
					resource.TestCheckResourceAttrSet("utem_integration.webhook_test", "id"),
				),
			},
		},
	})
}
