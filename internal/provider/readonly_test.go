package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

// read_only is the guard that has stopped more than one apply from doing something irreversible to
// a live system, so each hand-written resource is checked to refuse in that mode. That it still
// allows changes which only touch Terraform state is covered by
// TestReadOnlyRecordsCreateOnlyValuesAfterImport.

func TestReadOnlyRefusesFileWrite(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)
	srv.Files.Put("/mnt/pool/f", []byte("before"))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "read_only = true") + `
resource "truenas_file" "f" {
  path    = "/mnt/pool/f"
  content = "after"
}`,
			ExpectError: regexp.MustCompile(`read-only`),
		}},
	})
	if got, _ := srv.Files.Get("/mnt/pool/f"); string(got) != "before" {
		t.Errorf("read-only mode still wrote the file: %q", got)
	}
}

func TestReadOnlyRefusesACLChange(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	acls := serveACLs(srv, "/mnt/pool/share")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config: providerConfig(srv, "read_only = true") + `
resource "truenas_filesystem_acl" "share" {
  path = "/mnt/pool/share"
  entries = [
    { tag = "owner@", type = "ALLOW", perms = "FULL_CONTROL", flags = "INHERIT" },
  ]
}`,
			ExpectError: regexp.MustCompile(`read-only`),
		}},
	})
	if n := len(acls.acls["/mnt/pool/share"]["acl"].([]any)); n != 0 {
		t.Errorf("read-only mode still set %d ACL entries", n)
	}
}
