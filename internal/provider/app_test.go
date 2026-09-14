package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

// fakeApps behaves like middlewared's app.* for custom apps: portals and notes are read from the
// stored compose only on create, portals are rendered as URLs with explicit ports, and every
// call is recorded.
type fakeApps struct {
	mu    sync.Mutex
	apps  map[string]map[string]any
	calls []string
}

func serveApps(srv *middlewaretest.Server) *fakeApps {
	f := &fakeApps{apps: map[string]map[string]any{}}
	arg := func(params []json.RawMessage, i int, v any) { _ = json.Unmarshal(params[i], v) }
	record := func(call string) { f.calls = append(f.calls, call) }

	srv.Handle("app.create", func(_ context.Context, params []json.RawMessage) (any, error) {
		var req struct {
			Name   string         `json:"app_name"`
			Custom bool           `json:"custom_app"`
			Config map[string]any `json:"custom_compose_config"`
		}
		arg(params, 0, &req)
		f.mu.Lock()
		defer f.mu.Unlock()
		record("create " + req.Name)
		portals := map[string]string{}
		if list, ok := req.Config["x-portals"].([]any); ok {
			for _, p := range list {
				m := p.(map[string]any)
				portals[m["name"].(string)] = fmt.Sprintf("%s://%s:%v%s", m["scheme"], m["host"], m["port"], m["path"])
			}
		}
		entry := map[string]any{"name": req.Name, "state": "RUNNING", "custom_app": req.Custom, "portals": portals, "notes": nil, "config": req.Config}
		if n, ok := req.Config["x-notes"].(string); ok {
			entry["notes"] = n + "\n"
		}
		f.apps[req.Name] = entry
		return entry, nil
	})
	srv.Handle("app.update", func(_ context.Context, params []json.RawMessage) (any, error) {
		var name string
		var req struct {
			Config map[string]any `json:"custom_compose_config"`
		}
		arg(params, 0, &name)
		arg(params, 1, &req)
		f.mu.Lock()
		defer f.mu.Unlock()
		record("update " + name)
		f.apps[name]["config"] = req.Config // portals and notes deliberately not refreshed
		return f.apps[name], nil
	})
	simple := func(method, state string) {
		srv.Handle(method, func(_ context.Context, params []json.RawMessage) (any, error) {
			var name string
			arg(params, 0, &name)
			f.mu.Lock()
			defer f.mu.Unlock()
			record(method + " " + name)
			if state != "" {
				f.apps[name]["state"] = state
			}
			return true, nil
		})
	}
	simple("app.redeploy", "")
	simple("app.start", "RUNNING")
	simple("app.stop", "STOPPED")
	srv.Handle("app.delete", func(_ context.Context, params []json.RawMessage) (any, error) {
		var name string
		var opts map[string]any
		arg(params, 0, &name)
		arg(params, 1, &opts)
		f.mu.Lock()
		defer f.mu.Unlock()
		record(fmt.Sprintf("delete %s images=%v volumes=%v", name, opts["remove_images"], opts["remove_ix_volumes"]))
		delete(f.apps, name)
		return true, nil
	})
	srv.Handle("app.query", func(_ context.Context, params []json.RawMessage) (any, error) {
		var filters [][]any
		arg(params, 0, &filters)
		f.mu.Lock()
		defer f.mu.Unlock()
		out := []any{}
		for name, a := range f.apps {
			if len(filters) == 0 || filters[0][2] == name {
				out = append(out, a)
			}
		}
		return out, nil
	})
	return f
}

func (f *fakeApps) take() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := f.calls
	f.calls = nil
	sort.Strings(out)
	return out
}

func expectCalls(f *fakeApps, want ...string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		got := f.take()
		sort.Strings(want)
		if fmt.Sprint(got) != fmt.Sprint(want) {
			return fmt.Errorf("calls = %v, want %v", got, want)
		}
		return nil
	}
}

func TestAppLifecycle(t *testing.T) {
	const addr = "truenas_app.plex"
	srv := middlewaretest.NewServer(t)
	apps := serveApps(srv)
	config := func(extra string) string {
		return providerConfig(srv, "") + `
resource "truenas_app" "plex" {
  name    = "plex"
  include = ["/mnt/tank/apps/docker/plex.yml"]
  portals = {
    "Web UI" = "http://plex.example.lan/"
    "Direct" = "http://192.0.2.10:32400/web"
  }
  notes = "# Plex"
  ` + extra + `
}`
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		CheckDestroy:             expectCalls(apps, "delete plex images=false volumes=false"),
		Steps: []resource.TestStep{
			{
				Config: config(""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("state"), knownvalue.StringExact("RUNNING")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("portals").AtMapKey("Web UI"), knownvalue.StringExact("http://plex.example.lan/")),
				},
				Check: expectCalls(apps, "create plex"),
			},
			{
				// Portal URLs rendered with explicit ports and notes with a trailing newline are not drift.
				Config:           config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				ResourceName:                         addr,
				ImportState:                          true,
				ImportStateId:                        "plex",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"portals", "notes"},
			},
			{
				Config: config(`redeploy_trigger = "sha-1"`),
				Check:  expectCalls(apps, "app.redeploy plex"),
			},
			{
				Config: config(`redeploy_trigger = "sha-1"
  desired_state    = "STOPPED"`),
				Check: expectCalls(apps, "app.stop plex"),
			},
			{
				Config: config(`redeploy_trigger = "sha-1"`),
				Check:  expectCalls(apps, "app.start plex"),
			},
			{
				// A changed portal cannot be applied in place: TrueNAS reads portals only on create.
				Config: providerConfig(srv, "") + `
resource "truenas_app" "plex" {
  name    = "plex"
  include = ["/mnt/tank/apps/docker/plex.yml"]
  portals = { "Web UI" = "https://plex.example.com/" }
  notes   = "# Plex"
  redeploy_trigger = "sha-1"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionDestroyBeforeCreate)},
				},
				Check: expectCalls(apps, "delete plex images=false volumes=false", "create plex"),
			},
			{
				Config: providerConfig(srv, "") + `
resource "truenas_app" "plex" {
  name    = "plex"
  compose = "services:\n  plex:\n    image: plexinc/pms-docker\n"
  portals = { "Web UI" = "https://plex.example.com/" }
  notes   = "# Plex"
  redeploy_trigger = "sha-1"
  hold    = true
}`,
				Check: expectCalls(apps, "update plex"),
			},
			{
				Config: providerConfig(srv, "") + `
resource "truenas_app" "plex" {
  name    = "plex"
  compose = "services:\n  plex:\n    image: plexinc/pms-docker:latest\n"
  portals = { "Web UI" = "https://plex.example.com/" }
  notes   = "# Plex"
  redeploy_trigger = "sha-1"
  hold    = true
}`,
				ExpectError: regexp.MustCompile(`App is held`),
			},
			{
				// Releasing the hold is always allowed.
				Config: providerConfig(srv, "") + `
resource "truenas_app" "plex" {
  name    = "plex"
  compose = "services:\n  plex:\n    image: plexinc/pms-docker\n"
  portals = { "Web UI" = "https://plex.example.com/" }
  notes   = "# Plex"
  redeploy_trigger = "sha-1"
}`,
			},
		},
	})
}

func TestAppRefusesCatalogApps(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	apps := serveApps(srv)
	apps.apps["nextcloud"] = map[string]any{"name": "nextcloud", "state": "RUNNING", "custom_app": false, "config": map[string]any{}}
	_ = middleware.CodeCallError

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
import {
  to = truenas_app.nextcloud
  id = "nextcloud"
}
resource "truenas_app" "nextcloud" {
  name    = "nextcloud"
  include = ["/mnt/tank/apps/docker/nextcloud.yml"]
}`,
				ExpectError: regexp.MustCompile(`is a catalog app`),
			},
		},
	})
}

// sample_app_plex.json is plex's real app.query entry: an include-wrapped custom app whose
// portals and notes were written by an external metadata script.
func TestAppImportsRealIncludeWrappedAppWithoutChanges(t *testing.T) {
	raw, err := os.ReadFile("testdata/sample_app_plex.json")
	if err != nil {
		t.Fatal(err)
	}
	var plex map[string]any
	if err := json.Unmarshal(raw, &plex); err != nil {
		t.Fatal(err)
	}
	srv := middlewaretest.NewServer(t)
	apps := serveApps(srv)
	apps.apps["plex"] = plex

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
import {
  to = truenas_app.plex
  id = "plex"
}
resource "truenas_app" "plex" {
  name    = "plex"
  include = ["/mnt/tank/apps/docker/plex.yml"]
  portals = {
    "Direct" = "http://192.0.2.10:32400/web"
    "Web UI" = "http://plex.example.lan/"
  }
  hold = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction("truenas_app.plex", plancheck.ResourceActionUpdate)},
				},
				Check: expectCalls(apps),
			},
			{
				Config: providerConfig(srv, "") + `
resource "truenas_app" "plex" {
  name    = "plex"
  include = ["/mnt/tank/apps/docker/plex.yml"]
  portals = {
    "Direct" = "http://192.0.2.10:32400/web"
    "Web UI" = "http://plex.example.lan/"
  }
  hold = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				// Destroy is refused while held.
				Config:      providerConfig(srv, "") + `# empty`,
				ExpectError: regexp.MustCompile(`App is held`),
			},
			{
				Config: providerConfig(srv, "") + `
resource "truenas_app" "plex" {
  name    = "plex"
  include = ["/mnt/tank/apps/docker/plex.yml"]
  portals = {
    "Direct" = "http://192.0.2.10:32400/web"
    "Web UI" = "http://plex.example.lan/"
  }
}`,
			},
		},
	})
}

// Recording provider-side settings changes nothing on TrueNAS, so read-only mode must allow it;
// anything that would call TrueNAS must still be refused at plan time.
func TestReadOnlyAllowsStateOnlyChanges(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	apps := serveApps(srv)
	apps.apps["postgres"] = map[string]any{
		"name": "postgres", "state": "RUNNING", "custom_app": true, "portals": map[string]any{}, "notes": nil,
		"config": map[string]any{"include": []any{"/mnt/tank/apps/docker/postgres.yml"}},
	}
	app := func(extra string) string {
		return providerConfig(srv, "read_only = true") + `
resource "truenas_app" "postgres" {
  name    = "postgres"
  include = ["/mnt/tank/apps/docker/postgres.yml"]
  ` + extra + `
}`
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: app(`hold = true
  delete_images = false`) + `
import {
  to = truenas_app.postgres
  id = "postgres"
}`,
				Check: expectCalls(apps),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("truenas_app.postgres", tfjsonpath.New("hold"), knownvalue.Bool(true)),
				},
			},
			{
				Config: app(`hold = false
  delete_ix_volumes = true`),
				Check: expectCalls(apps),
			},
			{
				Config:      app(`redeploy_trigger = "new"`),
				ExpectError: regexp.MustCompile(`would\s+call\s+app.redeploy\s+on\s+postgres,\s+but\s+the\s+provider\s+has\s+read_only`),
			},
			{
				Config:      app(`desired_state = "STOPPED"`),
				ExpectError: regexp.MustCompile(`would\s+call\s+app.stop`),
			},
			{
				Config: app(`portals = { "UI" = "http://pg.lan/" }
  hold = false`),
				ExpectError: regexp.MustCompile(`would\s+replace\s+postgres`),
			},
			{
				Config:      providerConfig(srv, "read_only = true") + `# destroy`,
				ExpectError: regexp.MustCompile(`would\s+destroy\s+postgres`),
			},
			{
				// Leave the fixture removable.
				Config: providerConfig(srv, "") + `
resource "truenas_app" "postgres" {
  name    = "postgres"
  include = ["/mnt/tank/apps/docker/postgres.yml"]
}`,
			},
		},
	})
}
