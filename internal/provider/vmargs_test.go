package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

// TestVMOmittedArgsAreCleared pins a hazard on an imported VM. A configuration that leaves
// command_line_args and cpuset out does not leave them alone: the first is optional and computed
// with an empty default, the second plain optional, so an apply erases both. On a hand-tuned
// hypervisor configuration that is destructive, and the plan is the only warning.
func TestVMOmittedArgsAreCleared(t *testing.T) {
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
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("truenas_vm.tuned", plancheck.ResourceActionUpdate),
				},
			},
			// The apply goes through, and the tuning is gone.
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue("truenas_vm.tuned", tfjsonpath.New("cpuset"), knownvalue.Null()),
				statecheck.ExpectKnownValue("truenas_vm.tuned", tfjsonpath.New("command_line_args"),
					knownvalue.StringExact("")),
			},
		}},
	})
}
