package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

const cronJobAddr = "truenas_cron_job.backup"

func TestCronJobLifecycle(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	store := srv.ServeCRUD("cron_job")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_12_0)},
		CheckDestroy: func(*terraform.State) error {
			if rows := store.Rows(); len(rows) != 0 {
				return fmt.Errorf("%d cron jobs remain after destroy", len(rows))
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				// Create with defaults filled in from the API schema.
				Config: providerConfig(srv, "") + `
resource "truenas_cron_job" "backup" {
  command = "/usr/local/bin/backup.sh"
  user    = "root"
  schedule = {
    minute = "30"
    hour   = "2"
  }
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(cronJobAddr, tfjsonpath.New("id"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(cronJobAddr, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(cronJobAddr, tfjsonpath.New("ignore_stdout"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(cronJobAddr, tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(cronJobAddr, tfjsonpath.New("schedule"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"minute": knownvalue.StringExact("30"),
						"hour":   knownvalue.StringExact("2"),
						"dom":    knownvalue.StringExact("*"),
						"month":  knownvalue.StringExact("*"),
						"dow":    knownvalue.StringExact("*"),
					})),
					statecheck.ExpectIdentityValue(cronJobAddr, tfjsonpath.New("id"), knownvalue.Int64Exact(1)),
				},
			},
			{
				// Update in place sends only what changed.
				Config: providerConfig(srv, "") + `
resource "truenas_cron_job" "backup" {
  command     = "/usr/local/bin/backup.sh"
  user        = "root"
  description = "nightly backup"
  enabled     = false
  schedule = {
    minute = "30"
    hour   = "2"
  }
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(cronJobAddr, plancheck.ResourceActionUpdate)},
				},
				Check: func(*terraform.State) error {
					row := store.Rows()[0]
					if row["description"] != "nightly backup" || row["enabled"] != false {
						return fmt.Errorf("store not updated: %v", row)
					}
					return lastUpdateFields(srv, "description", "enabled")
				},
			},
			{
				ResourceName:      cronJobAddr,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:    cronJobAddr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
			{
				// A change made outside Terraform shows up as drift and is corrected.
				PreConfig: func() {
					store.Mutate(1, func(row map[string]any) { row["command"] = "rm -rf /tmp/x" })
				},
				Config: providerConfig(srv, "") + `
resource "truenas_cron_job" "backup" {
  command     = "/usr/local/bin/backup.sh"
  user        = "root"
  description = "nightly backup"
  enabled     = false
  schedule = {
    minute = "30"
    hour   = "2"
  }
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(cronJobAddr, plancheck.ResourceActionUpdate)},
				},
				Check: func(*terraform.State) error {
					if got := store.Rows()[0]["command"]; got != "/usr/local/bin/backup.sh" {
						return fmt.Errorf("drift not corrected: command = %v", got)
					}
					return nil
				},
			},
			{
				// Read-only mode refuses a plan with changes.
				Config: providerConfig(srv, "read_only = true") + `
resource "truenas_cron_job" "backup" {
  command     = "/usr/local/bin/backup.sh"
  user        = "root"
  description = "changed while read-only"
  enabled     = false
  schedule = {
    minute = "30"
    hour   = "2"
  }
}`,
				ExpectError: regexp.MustCompile(`Provider is read-only`),
			},
			{
				// Deletion outside Terraform plans a re-create.
				PreConfig: func() { store.Remove(1) },
				Config: providerConfig(srv, "") + `
resource "truenas_cron_job" "backup" {
  command     = "/usr/local/bin/backup.sh"
  user        = "root"
  description = "nightly backup"
  enabled     = false
  schedule = {
    minute = "30"
    hour   = "2"
  }
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(cronJobAddr, plancheck.ResourceActionCreate)},
				},
			},
		},
	})
}

func TestCronJobValidationErrorsPointAtAttributes(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.ServeCRUD("cron_job")
	srv.Handle("cronjob.create", func(context.Context, []json.RawMessage) (any, error) {
		return nil, &middleware.Error{Errno: 22, Errname: "EINVAL", Fields: []middleware.FieldError{
			{Attribute: "cronjob_create.schedule.minute", Message: "Invalid minute value 61", Errno: 22},
		}}
	})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
resource "truenas_cron_job" "bad" {
  command  = "true"
  user     = "root"
  schedule = { minute = "61" }
}`,
				ExpectError: regexp.MustCompile(`(?s)schedule.*Invalid minute value 61`),
			},
		},
	})
}

// lastUpdateFields asserts that the most recent cronjob.update sent exactly the named fields.
func lastUpdateFields(srv *middlewaretest.Server, want ...string) error {
	calls := srv.Calls()
	for i := len(calls) - 1; i >= 0; i-- {
		if calls[i].Method != "cronjob.update" {
			continue
		}
		var patch map[string]any
		if err := json.Unmarshal(calls[i].Params[1], &patch); err != nil {
			return err
		}
		if len(patch) != len(want) {
			return fmt.Errorf("update sent %v, want only %v", patch, want)
		}
		for _, f := range want {
			if _, ok := patch[f]; !ok {
				return fmt.Errorf("update sent %v, missing %s", patch, f)
			}
		}
		return nil
	}
	return fmt.Errorf("no cronjob.update call")
}
