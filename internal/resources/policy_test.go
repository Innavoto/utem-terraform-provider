package resources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyResource_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the policy
			{
				Config: testAccProviderConfig + `
resource "utem_policy" "test" {
  name            = "tf-acc-test-s3-public"
  description     = "Detect publicly accessible S3 buckets"
  category        = "cloud"
  severity        = "critical"
  is_enabled      = true
  resource_type   = "aws_s3_bucket"
  customer_facing = false
  rego_policy     = <<-EOT
    package cloud.s3_public
    default allow = false
    allow { input.public_access == false }
  EOT
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_policy.test", "name", "tf-acc-test-s3-public"),
					resource.TestCheckResourceAttr("utem_policy.test", "category", "cloud"),
					resource.TestCheckResourceAttr("utem_policy.test", "severity", "critical"),
					resource.TestCheckResourceAttr("utem_policy.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("utem_policy.test", "resource_type", "aws_s3_bucket"),
					resource.TestCheckResourceAttr("utem_policy.test", "customer_facing", "false"),
					resource.TestCheckResourceAttrSet("utem_policy.test", "id"),
				),
			},
			// Step 2: Update the policy (change name + toggle is_enabled)
			{
				Config: testAccProviderConfig + `
resource "utem_policy" "test" {
  name            = "tf-acc-test-s3-public-updated"
  description     = "Updated: detect publicly accessible S3 buckets"
  category        = "cloud"
  severity        = "high"
  is_enabled      = false
  resource_type   = "aws_s3_bucket"
  customer_facing = true
  rego_policy     = <<-EOT
    package cloud.s3_public
    default allow = false
    allow { input.public_access == false }
  EOT
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_policy.test", "name", "tf-acc-test-s3-public-updated"),
					resource.TestCheckResourceAttr("utem_policy.test", "severity", "high"),
					resource.TestCheckResourceAttr("utem_policy.test", "is_enabled", "false"),
					resource.TestCheckResourceAttr("utem_policy.test", "customer_facing", "true"),
					resource.TestCheckResourceAttrSet("utem_policy.test", "id"),
				),
			},
		},
	})
}

func TestAccPolicyResource_aiAgents(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + `
resource "utem_policy" "ai_test" {
  name            = "tf-acc-test-ai-agent-scope"
  description     = "Detect AI agents with excessive scope"
  category        = "ai-agents"
  severity        = "high"
  is_enabled      = true
  resource_type   = "ai_agent"
  customer_facing = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_policy.ai_test", "name", "tf-acc-test-ai-agent-scope"),
					resource.TestCheckResourceAttr("utem_policy.ai_test", "category", "ai-agents"),
					resource.TestCheckResourceAttr("utem_policy.ai_test", "severity", "high"),
					resource.TestCheckResourceAttr("utem_policy.ai_test", "customer_facing", "true"),
					resource.TestCheckResourceAttrSet("utem_policy.ai_test", "id"),
				),
			},
		},
	})
}
