package resources_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScanScheduleResource_basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the scan schedule
			{
				Config: testAccProviderConfig + `
resource "utem_scan_schedule" "test" {
  name            = "tf-acc-test-nightly"
  description     = "Nightly full scan"
  target_host     = "192.168.1.100"
  scan_type       = "full"
  cron_expression = "0 2 * * *"
  is_enabled      = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "name", "tf-acc-test-nightly"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "description", "Nightly full scan"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "target_host", "192.168.1.100"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "scan_type", "full"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "cron_expression", "0 2 * * *"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "is_enabled", "true"),
					resource.TestCheckResourceAttrSet("utem_scan_schedule.test", "id"),
				),
			},
			// Step 2: Update the scan schedule (change name + toggle is_enabled)
			{
				Config: testAccProviderConfig + `
resource "utem_scan_schedule" "test" {
  name            = "tf-acc-test-nightly-updated"
  description     = "Updated nightly scan"
  target_host     = "192.168.1.100"
  scan_type       = "quick"
  cron_expression = "0 3 * * *"
  is_enabled      = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "name", "tf-acc-test-nightly-updated"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "description", "Updated nightly scan"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "scan_type", "quick"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "cron_expression", "0 3 * * *"),
					resource.TestCheckResourceAttr("utem_scan_schedule.test", "is_enabled", "false"),
					resource.TestCheckResourceAttrSet("utem_scan_schedule.test", "id"),
				),
			},
		},
	})
}

func TestAccScanScheduleResource_withModules(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig + `
resource "utem_scan_schedule" "modules_test" {
  name            = "tf-acc-test-targeted"
  description     = "Targeted scan with modules"
  target_host     = "10.0.0.50"
  scan_type       = "targeted"
  cron_expression = "0 4 * * 1"
  is_enabled      = true
  modules         = ["port_scan", "vuln_scan", "web_scan"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("utem_scan_schedule.modules_test", "name", "tf-acc-test-targeted"),
					resource.TestCheckResourceAttr("utem_scan_schedule.modules_test", "scan_type", "targeted"),
					resource.TestCheckResourceAttr("utem_scan_schedule.modules_test", "modules.#", "3"),
					resource.TestCheckResourceAttr("utem_scan_schedule.modules_test", "modules.0", "port_scan"),
					resource.TestCheckResourceAttr("utem_scan_schedule.modules_test", "modules.1", "vuln_scan"),
					resource.TestCheckResourceAttr("utem_scan_schedule.modules_test", "modules.2", "web_scan"),
					resource.TestCheckResourceAttrSet("utem_scan_schedule.modules_test", "id"),
				),
			},
		},
	})
}
