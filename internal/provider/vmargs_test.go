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
// configuration never mentions. command_line_args is reported by the server, so the provider
// imposes no default on it and an unmentioned value is left alone.
//
// cpuset is the other half, and it is not fixed: a plain optional field plans as null when the
// configuration omits it, so pinning is erased. Adopting a hand-tuned machine means writing every
// field that matters, and this test says which half of that is the provider's doing.
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
				// Plain optional: omitting it still means null, and the pinning is lost.
				statecheck.ExpectKnownValue("truenas_vm.tuned", tfjsonpath.New("cpuset"), knownvalue.Null()),
			},
		}},
	})
}
