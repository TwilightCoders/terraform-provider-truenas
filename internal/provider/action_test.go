package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

func TestActionsCallMethodsWithPositionalArguments(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	var mu sync.Mutex
	got := map[string]string{}
	record := func(method string) {
		srv.Handle(method, func(_ context.Context, params []json.RawMessage) (any, error) {
			b, _ := json.Marshal(params)
			mu.Lock()
			got[method] = string(b)
			mu.Unlock()
			return true, nil
		})
	}
	record("replication.run")
	record("cloudsync.sync")
	record("app.redeploy")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_14_0)},
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
action "truenas_replication_run" "archive" {
  config {
    id = 7
  }
}
action "truenas_cloudsync_sync" "photos" {
  config {
    id      = 3
    options = { dry_run = true }
  }
}
action "truenas_app_redeploy" "plex" {
  config {
    name = "plex"
  }
}
resource "terraform_data" "deploy" {
  input = "v1"
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.truenas_replication_run.archive, action.truenas_cloudsync_sync.photos, action.truenas_app_redeploy.plex]
    }
  }
}`,
				Check: func(*terraform.State) error {
					mu.Lock()
					defer mu.Unlock()
					for method, want := range map[string]string{
						"replication.run": `[7]`,
						"cloudsync.sync":  `[3,{"dry_run":true}]`,
						"app.redeploy":    `["plex"]`,
					} {
						if got[method] != want {
							return fmt.Errorf("%s called with %s, want %s", method, got[method], want)
						}
					}
					return nil
				},
			},
		},
	})
}
