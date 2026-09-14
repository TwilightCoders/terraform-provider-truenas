package provider

import (
	"encoding/json"
	"fmt"
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

// lifecycle runs create, both import styles, an in-place update and destroy for one resource.
type lifecycle struct {
	namespace string
	address   string
	create    string
	update    string
	checks    []statecheck.StateCheck
	onWrite   func(row map[string]any)
	// importIgnore lists attributes that cannot survive import, such as write-only companions.
	importIgnore []string
	// importUpdates is set when an import block plans an update that records create-only values
	// TrueNAS never reports; such an update sends nothing to TrueNAS.
	importUpdates bool
}

func (lc lifecycle) run(t *testing.T) {
	t.Helper()
	srv := middlewaretest.NewServer(t)
	store := srv.ServeCRUD(apischema.MustLoad(api.Latest), lc.namespace)
	store.OnWrite = lc.onWrite

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		CheckDestroy: func(*terraform.State) error {
			if n := len(store.Rows()); n != 0 {
				return fmt.Errorf("%d %s rows remain after destroy", n, lc.namespace)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config:            providerConfig(srv, "") + lc.create,
				ConfigStateChecks: lc.checks,
			},
			{
				// Re-planning the same configuration must be a no-op.
				Config: providerConfig(srv, "") + lc.create,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ResourceName:            lc.address,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: lc.importIgnore,
			},
			{
				ResourceName:    lc.address,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
				SkipFunc:        func() (bool, error) { return lc.importUpdates, nil },
			},
			{
				Config: providerConfig(srv, "") + lc.update,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(lc.address, plancheck.ResourceActionUpdate)},
				},
			},
			{
				Config: providerConfig(srv, "") + lc.update,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestSnapshotTaskLifecycle(t *testing.T) {
	const addr = "truenas_snapshot_task.home"
	lifecycle{
		namespace: "pool.snapshottask",
		address:   addr,
		create: `
resource "truenas_snapshot_task" "home" {
  dataset        = "tank/home"
  recursive      = true
  lifetime_value = 14
  lifetime_unit  = "DAY"
  naming_schema  = "auto-%Y%m%d-%H%M"
  schedule = {
    minute = "45"
    hour   = "2"
  }
}`,
		update: `
resource "truenas_snapshot_task" "home" {
  dataset        = "tank/home"
  recursive      = true
  lifetime_value = 30
  lifetime_unit  = "DAY"
  naming_schema  = "auto-%Y%m%d-%H%M"
  schedule = {
    minute = "45"
    hour   = "2"
  }
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule").AtMapKey("end"), knownvalue.StringExact("23:59")),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("allow_empty"), knownvalue.Bool(true)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("exclude"), knownvalue.ListExact(nil)),
		},
	}.run(t)
}

func TestSMBShareLifecycle(t *testing.T) {
	const addr = "truenas_smb_share.media"
	lifecycle{
		namespace: "sharing.smb",
		address:   addr,
		// Like middlewared, return options without the purpose that selects them.
		onWrite: func(row map[string]any) {
			if opts, ok := row["options"].(map[string]any); ok {
				delete(opts, "purpose")
			}
		},
		create: `
resource "truenas_smb_share" "media" {
  name = "Media"
  path = "/mnt/tank/data/media"
}`,
		update: `
resource "truenas_smb_share" "media" {
  name    = "Media"
  path    = "/mnt/tank/data/media"
  purpose = "TIMEMACHINE_SHARE"
  comment = "Backups"
  options = {
    timemachine_share = {
      timemachine_quota = 500
    }
  }
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("purpose"), knownvalue.StringExact("DEFAULT_SHARE")),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("options"), knownvalue.Null()),
		},
	}.run(t)
}

func TestInitScriptLifecycle(t *testing.T) {
	lifecycle{
		namespace: "initshutdownscript",
		address:   "truenas_init_script.governor",
		create: `
resource "truenas_init_script" "governor" {
  type    = "COMMAND"
  when    = "POSTINIT"
  command = "echo performance > /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor"
  comment = "CPU performance governor"
}`,
		update: `
resource "truenas_init_script" "governor" {
  type    = "COMMAND"
  when    = "POSTINIT"
  command = "echo performance > /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor"
  comment = "CPU performance governor"
  timeout = 30
}`,
	}.run(t)
}

func TestUserLifecycle(t *testing.T) {
	const addr = "truenas_user.media"
	lifecycle{
		namespace: "user",
		address:   addr,
		// Like middlewared: assign a uid, and return the primary group as an object.
		onWrite: func(row map[string]any) {
			if row["uid"] == nil {
				row["uid"] = json.Number("3000")
			}
			if id, ok := row["group"].(json.Number); ok {
				row["group"] = map[string]any{"id": id, "bsdgrp_gid": json.Number("3000"), "bsdgrp_group": "media"}
			}
		},
		create: `
resource "truenas_user" "media" {
  username  = "media"
  full_name = "Media services"
  group     = 116
}`,
		update: `
resource "truenas_user" "media" {
  username            = "media"
  full_name           = "Media services"
  group               = 116
  password            = "correct horse battery staple, rotated"
  password_wo_version = 2
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("password"), knownvalue.Null()),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("uid"), knownvalue.Int64Exact(3000)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("group"), knownvalue.Int64Exact(116)),
		},
	}.run(t)
}

func TestGroupLifecycle(t *testing.T) {
	lifecycle{
		namespace: "group",
		address:   "truenas_group.media",
		onWrite: func(row map[string]any) {
			if row["gid"] == nil {
				row["gid"] = json.Number("3001")
			}
		},
		create: `
resource "truenas_group" "media" {
  name = "media"
}`,
		update: `
resource "truenas_group" "media" {
  name          = "media"
  sudo_commands = ["/usr/bin/systemctl restart plex"]
}`,
	}.run(t)
}

func TestVMLifecycle(t *testing.T) {
	const addr = "truenas_vm.vm1"
	lifecycle{
		namespace: "vm",
		address:   addr,
		onWrite: func(row map[string]any) {
			if row["uuid"] == nil {
				row["uuid"] = "1b4e28ba-2fa1-11d2-883f-0016d3cca427"
			}
			if row["machine_type"] == nil {
				row["machine_type"] = "q35"
			}
			if row["arch_type"] == nil {
				row["arch_type"] = "x86_64"
			}
		},
		create: `
resource "truenas_vm" "vm1" {
  name   = "vm1"
  memory = 8192
  vcpus  = 1
  cores  = 4
}`,
		update: `
resource "truenas_vm" "vm1" {
  name      = "vm1"
  memory    = 16384
  vcpus     = 1
  cores     = 4
  autostart = false
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("uuid"), knownvalue.StringExact("1b4e28ba-2fa1-11d2-883f-0016d3cca427")),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("bootloader"), knownvalue.StringExact("UEFI")),
		},
	}.run(t)
}

func TestVMDeviceLifecycle(t *testing.T) {
	const addr = "truenas_vm_device.display"
	lifecycle{
		namespace: "vm.device",
		address:   addr,
		create: `
resource "truenas_vm_device" "display" {
  vm = 4
  attributes = {
    display = {
      type     = "SPICE"
      password = "view-only"
    }
  }
}`,
		update: `
resource "truenas_vm_device" "display" {
  vm    = 4
  order = 1002
  attributes = {
    display = {
      type     = "SPICE"
      password = "view-only"
    }
  }
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("attributes").AtMapKey("display").AtMapKey("password"), knownvalue.Null()),
		},
	}.run(t)
}
