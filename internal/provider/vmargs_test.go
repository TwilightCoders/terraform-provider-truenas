package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

// TestVMOmittedArgsSurvive covers what an imported VM does with hypervisor settings a
// configuration never mentions: it leaves them alone. A hand-tuned command line and CPU pinning
// are exactly what a partial configuration must not quietly erase.
//
// Removing one deliberately is what unset is for, which the second step exercises.
func TestVMOmittedArgsSurvive(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	store := srv.ServeCRUD("vm")
	store.Seed(map[string]any{
		"id": 1, "name": "tuned", "memory": 8192, "cpuset": "0,16,1,17",
		"command_line_args": "-smbios 'type=0,vendor=X,, LLC.'",
	})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "") + `
resource "truenas_vm" "tuned" {
  name   = "tuned"
  memory = 8192
}`,
			ResourceName:       "truenas_vm.tuned",
			ImportState:        true,
			ImportStateId:      "1",
			ImportStatePersist: true,
		}, {
			Config: providerConfig(srv, "") + `
resource "truenas_vm" "tuned" {
  name   = "tuned"
  memory = 8192
}`,
			ConfigStateChecks: []statecheck.StateCheck{
				// Reported by the server, so no default is imposed and it survives being unmentioned.
				statecheck.ExpectKnownValue("truenas_vm.tuned", tfjsonpath.New("command_line_args"),
					knownvalue.StringExact("-smbios 'type=0,vendor=X,, LLC.'")),
				// Also reported by the server, so pinning survives being unmentioned.
				statecheck.ExpectKnownValue("truenas_vm.tuned", tfjsonpath.New("cpuset"),
					knownvalue.StringExact("0,16,1,17")),
			},
		}, {
			// Clearing is possible, and has to be said out loud.
			Config: providerConfig(srv, "") + `
resource "truenas_vm" "tuned" {
  name   = "tuned"
  memory = 8192
  unset  = ["cpuset"]
}`,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue("truenas_vm.tuned", tfjsonpath.New("cpuset"), knownvalue.Null()),
				statecheck.ExpectKnownValue("truenas_vm.tuned", tfjsonpath.New("command_line_args"),
					knownvalue.StringExact("-smbios 'type=0,vendor=X,, LLC.'")),
			},
		}},
	})
}
