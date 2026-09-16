package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// asObject simulates middlewared returning a referenced object where an id was written.
func asObject(fields ...string) func(row map[string]any) {
	return func(row map[string]any) {
		for _, f := range fields {
			switch v := row[f].(type) {
			case json.Number:
				row[f] = map[string]any{"id": v, "name": "referenced"}
			case []any:
				for i, item := range v {
					if n, ok := item.(json.Number); ok {
						v[i] = map[string]any{"id": n, "name": "referenced"}
					}
				}
			}
		}
	}
}

func TestReplicationLifecycle(t *testing.T) {
	const addr = "truenas_replication.archive"
	lifecycle{
		resourceType: "replication",
		address:      addr,
		onWrite:      asObject("ssh_credentials", "periodic_snapshot_tasks"),
		create: `
resource "truenas_replication" "archive" {
  name                    = "archive to backup host"
  direction               = "PUSH"
  transport               = "SSH"
  ssh_credentials         = 4
  source_datasets         = ["tank/archive"]
  target_dataset          = "backup/truenas/archive"
  recursive               = true
  auto                    = true
  retention_policy        = "SOURCE"
  periodic_snapshot_tasks = [2]
}`,
		update: `
resource "truenas_replication" "archive" {
  name                    = "archive to backup host"
  direction               = "PUSH"
  transport               = "SSH"
  ssh_credentials         = 4
  source_datasets         = ["tank/archive"]
  target_dataset          = "backup/truenas/archive"
  recursive               = true
  auto                    = true
  retention_policy        = "SOURCE"
  periodic_snapshot_tasks = [2, 3]
  retries                 = 3
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("ssh_credentials"), knownvalue.Int64Exact(4)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("periodic_snapshot_tasks"), knownvalue.ListExact([]knownvalue.Check{knownvalue.Int64Exact(2)})),
		},
	}.run(t)
}

func TestCloudSyncCredentialsLifecycle(t *testing.T) {
	const addr = "truenas_cloudsync_credentials.b2"
	lifecycle{
		resourceType: "cloudsync_credentials",
		address:      addr,
		create: `
resource "truenas_cloudsync_credentials" "b2" {
  name = "Backblaze"
  storage = {
    b2 = {
      account = "0012345"
      key     = "K001secret"
    }
  }
}`,
		update: `
resource "truenas_cloudsync_credentials" "b2" {
  name = "Backblaze"
  storage = {
    b2 = {
      account = "0012345"
      key     = "K001rotated"
    }
  }
  storage_wo_version = 1
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("storage").AtMapKey("b2").AtMapKey("account"), knownvalue.StringExact("0012345")),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("storage").AtMapKey("b2").AtMapKey("key"), knownvalue.Null()),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("storage").AtMapKey("s3"), knownvalue.Null()),
		},
	}.run(t)
}

func TestWriteOnlySecretsAreSentOnlyOnVersionBump(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	store := srv.ServeCRUD("cloudsync_credentials")
	config := func(key, version string) string {
		return providerConfig(srv, "") + `
resource "truenas_cloudsync_credentials" "b2" {
  name = "Backblaze"
  storage = {
    b2 = {
      account = "0012345"
      key     = "` + key + `"
    }
  }
  ` + version + `
}`
	}
	storedKey := func() any {
		return store.Rows()[0]["provider"].(map[string]any)["key"]
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: config("first", ""),
				Check: func(*terraform.State) error {
					if got := storedKey(); got != "first" {
						return fmt.Errorf("created key = %v", got)
					}
					return nil
				},
			},
			{
				// Changing only a write-only value is invisible to Terraform: nothing to apply.
				Config:             config("second", ""),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: config("second", "storage_wo_version = 1"),
				Check: func(*terraform.State) error {
					if got := storedKey(); got != "second" {
						return fmt.Errorf("key after version bump = %v", got)
					}
					return nil
				},
			},
		},
	})
}

func TestCloudSyncTaskLifecycle(t *testing.T) {
	const addr = "truenas_cloudsync_task.photos"
	lifecycle{
		resourceType: "cloudsync_task",
		address:      addr,
		onWrite:      asObject("credentials"),
		create: `
resource "truenas_cloudsync_task" "photos" {
  path          = "/mnt/tank/data/cloud"
  credentials   = 1
  direction     = "PUSH"
  transfer_mode = "SYNC"
  attributes = {
    bucket = "truenas-photos"
    folder = "/"
  }
}`,
		update: `
resource "truenas_cloudsync_task" "photos" {
  path                        = "/mnt/tank/data/cloud"
  credentials                 = 1
  direction                   = "PUSH"
  transfer_mode               = "SYNC"
  encryption                  = true
  encryption_password         = "correct horse"
  encryption_password_wo_version = 1
  attributes = {
    bucket = "truenas-photos"
    folder = "/"
  }
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("credentials"), knownvalue.Int64Exact(1)),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("schedule").AtMapKey("minute"), knownvalue.StringExact("00")),
		},
	}.run(t)
}
