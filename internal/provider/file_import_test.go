package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

// TestFileImport covers adopting a file that already exists. Import carries only the path, so the
// resource has to decide the rest: content drift detection on, and a destroy that leaves the file
// alone. Getting either default wrong on an adopted file is destructive.
func TestFileImport(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)
	srv.Files.Put("/mnt/pool/apps/existing.yml", []byte("services: {}\n"))

	cfg := providerConfig(srv, "") + `
resource "truenas_file" "existing" {
  path    = "/mnt/pool/apps/existing.yml"
  content = "services: {}\n"
}`

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config:             cfg,
				ResourceName:       "truenas_file.existing",
				ImportState:        true,
				ImportStateId:      "/mnt/pool/apps/existing.yml",
				ImportStatePersist: true,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("truenas_file.existing", tfjsonpath.New("drift_detection"),
						knownvalue.StringExact("content")),
					statecheck.ExpectKnownValue("truenas_file.existing", tfjsonpath.New("on_destroy"),
						knownvalue.StringExact("leave")),
				},
			},
			{
				// The imported file already matches the configuration, so adopting it changes nothing.
				Config: cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})

	if got, ok := srv.Files.Get("/mnt/pool/apps/existing.yml"); !ok || string(got) != "services: {}\n" {
		t.Errorf("import rewrote the file: %q", got)
	}
}

// TestFileRejectsBadMode covers the one configuration error the resource has to catch itself:
// permission bits that are not octal.
func TestFileRejectsBadMode(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "") + `
resource "truenas_file" "bad" {
  path    = "/mnt/pool/f"
  content = "x"
  mode    = "0999"
}`,
			ExpectError: regexp.MustCompile(`octal permission bits`),
		}},
	})
}

// TestFileRelativePathRejected covers the other: a path that is not absolute would be resolved
// against whatever directory middleware happens to be in.
func TestFileRelativePathRejected(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "") + `
resource "truenas_file" "rel" {
  path    = "apps/compose.yml"
  content = "x"
}`,
			ExpectError: regexp.MustCompile(`absolute path`),
		}},
	})
}
