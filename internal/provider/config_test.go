package provider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/api"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/apischema"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

func TestSSHConfigSingleton(t *testing.T) {
	const addr = "truenas_ssh_config.this"
	srv := middlewaretest.NewServer(t)
	initial := map[string]any{
		"id": json.Number("1"), "tcpport": json.Number("22"), "passwordauth": true, "options": "",
		"host_rsa_key": "-----BEGIN OPENSSH PRIVATE KEY-----secret", "host_rsa_key_pub": "ssh-rsa AAAA",
	}
	cfg := srv.ServeConfig(apischema.MustLoad(api.Latest), "ssh", initial)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		CheckDestroy: func(*terraform.State) error {
			// Destroy only forgets the settings; TrueNAS keeps them.
			if got := cfg.Data()["passwordauth"]; got != false {
				return fmt.Errorf("destroy changed settings: passwordauth = %v", got)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
resource "truenas_ssh_config" "this" {
  passwordauth = false
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("passwordauth"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("tcpport"), knownvalue.Int64Exact(22)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("host_rsa_key_pub"), knownvalue.StringExact("ssh-rsa AAAA")),
				},
				Check: func(*terraform.State) error {
					// Adopting the settings sends only what is configured.
					for _, c := range srv.Calls() {
						if c.Method != "ssh.update" {
							continue
						}
						var patch map[string]any
						_ = json.Unmarshal(c.Params[0], &patch)
						if !reflect.DeepEqual(patch, map[string]any{"passwordauth": false}) {
							return fmt.Errorf("ssh.update sent %v", patch)
						}
					}
					return nil
				},
			},
			{
				Config: providerConfig(srv, "") + `
resource "truenas_ssh_config" "this" {
  passwordauth = false
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				ResourceName:      addr,
				ImportState:       true,
				ImportStateId:     "ssh",
				ImportStateVerify: true,
			},
			{
				Config: providerConfig(srv, "") + `
resource "truenas_ssh_config" "this" {
  passwordauth = false
  tcpport      = 2222
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate)},
				},
				Check: func(*terraform.State) error {
					if got := cfg.Data()["tcpport"]; fmt.Sprint(got) != "2222" {
						return fmt.Errorf("tcpport = %v", got)
					}
					return nil
				},
			},
		},
	})
}
