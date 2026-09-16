package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

// serveStat answers filesystem.stat from the fake server's in-memory files. Middleware reports
// mtime as a {"$date": milliseconds} wrapper, so the fake does too.
func serveStat(srv *middlewaretest.Server) {
	srv.Handle("filesystem.stat", func(_ context.Context, params []json.RawMessage) (any, error) {
		var p string
		_ = json.Unmarshal(params[0], &p)
		content, ok := srv.Files.Get(p)
		if !ok {
			return nil, &middleware.Error{Errno: 2, Errname: "ENOENT", Reason: p + " does not exist"}
		}
		return map[string]any{
			"size": len(content), "mode": 0o100644,
			"mtime": map[string]any{"$date": int64(1700000000000)},
		}, nil
	})
}

func fileConfig(srv *middlewaretest.Server, content, extra string) string {
	return providerConfig(srv, "") + fmt.Sprintf(`
resource "truenas_file" "compose" {
  path    = "/mnt/pool/apps/compose.yml"
  content = %q
  %s
}
`, content, extra)
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func TestFileCreatesAndKeepsContentOutOfState(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config: fileConfig(srv, "services: {}\n", ""),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue("truenas_file.compose", tfjsonpath.New("content_sha256"),
					knownvalue.StringExact(sha256Hex("services: {}\n"))),
				// Write-only: the contents must never be persisted.
				statecheck.ExpectKnownValue("truenas_file.compose", tfjsonpath.New("content"),
					knownvalue.Null()),
				statecheck.ExpectKnownValue("truenas_file.compose", tfjsonpath.New("id"),
					knownvalue.StringExact("/mnt/pool/apps/compose.yml")),
			},
			Check: func(*terraform.State) error {
				got, ok := srv.Files.Get("/mnt/pool/apps/compose.yml")
				if !ok || string(got) != "services: {}\n" {
					return fmt.Errorf("file on host = %q (present=%v)", got, ok)
				}
				return nil
			},
		}},
	})
}

// TestFileDetectsRemoteDrift is the point of the resource: an edit made on the host, which a local
// filesha256() cannot see, must show up as a change.
func TestFileDetectsRemoteDrift(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)
	cfg := fileConfig(srv, "services: {}\n", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				PreConfig: func() {
					srv.Files.Put("/mnt/pool/apps/compose.yml", []byte("services: {edited: true}\n"))
				},
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("truenas_file.compose", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				// The apply put the configured contents back.
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestFileDriftDetectionNone covers the mode meant for secrets: nothing derived from the contents
// is refreshed, so an edit on the host is deliberately invisible.
func TestFileDriftDetectionNone(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)
	cfg := fileConfig(srv, "TOKEN=abc\n", `drift_detection = "none"`)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				PreConfig: func() {
					srv.Files.Put("/mnt/pool/apps/compose.yml", []byte("TOKEN=tampered\n"))
				},
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestFileDriftDetectionStat documents the lossy fast path: an edit that preserves size and mtime
// is not noticed, which is exactly why it is not the default.
func TestFileDriftDetectionStat(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)
	cfg := fileConfig(srv, "aaaa\n", `drift_detection = "stat"`)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				PreConfig: func() {
					// Same length, same mtime: stat cannot tell.
					srv.Files.Put("/mnt/pool/apps/compose.yml", []byte("bbbb\n"))
				},
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestFileOnDestroyLeavesFileByDefault(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: fileConfig(srv, "keep me\n", "")},
			{Config: providerConfig(srv, "")}, // resource removed: destroy
		},
	})
	if got, ok := srv.Files.Get("/mnt/pool/apps/compose.yml"); !ok || string(got) != "keep me\n" {
		t.Errorf("destroy did not leave the file intact: %q (present=%v)", got, ok)
	}
}

func TestFileOnDestroyTruncate(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: fileConfig(srv, "goodbye\n", `on_destroy = "truncate"`)},
			{Config: providerConfig(srv, "")},
		},
	})
	got, ok := srv.Files.Get("/mnt/pool/apps/compose.yml")
	if !ok || len(got) != 0 {
		t.Errorf("destroy did not truncate: %q (present=%v)", got, ok)
	}
}

// TestFileReadOnlyRefusesWriteButAllowsStateOnlyChange mirrors the rule the rest of the provider
// follows: read_only refuses calls to TrueNAS, not changes that only touch state.
func TestFileReadOnlyRefusesWrite(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "read_only = true") + `
resource "truenas_file" "compose" {
  path    = "/mnt/pool/apps/compose.yml"
  content = "nope\n"
}
`,
			ExpectError: regexp.MustCompile(`read-only`),
		}},
	})
}
