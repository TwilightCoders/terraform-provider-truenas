package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/resources"
)

// wrapProperties rewrites plain values into the {value, rawvalue, parsed, source} objects
// pool.dataset returns, the way middlewared reports ZFS properties.
func wrapProperties(t *testing.T) func(row map[string]any) {
	t.Helper()
	shape, ok := resources.ShapeFor("dataset")
	if !ok {
		t.Fatal("no dataset shape")
	}
	var props []string
	for _, f := range shape.Fields {
		if f.Property {
			props = append(props, f.Name)
		}
	}
	return func(row map[string]any) {
		for _, name := range props {
			if w, ok := row[name].(map[string]any); ok && (w["source"] == "" || w["source"] == nil) {
				row[name] = nil // zero value filled in by the store: never set
			}
			switch v := row[name].(type) {
			case map[string]any:
				// already wrapped
			case nil:
				row[name] = map[string]any{"value": nil, "rawvalue": "0", "parsed": nil, "source": "DEFAULT"}
			case string:
				if v == "INHERIT" {
					row[name] = map[string]any{"value": "INHERITED-VALUE", "rawvalue": "inherited-value", "parsed": "inherited-value", "source": "INHERITED"}
				} else {
					row[name] = map[string]any{"value": v, "rawvalue": strings.ToLower(v), "parsed": strings.ToLower(v), "source": "LOCAL"}
				}
			default:
				row[name] = map[string]any{"value": fmt.Sprint(v), "rawvalue": fmt.Sprint(v), "parsed": v, "source": "LOCAL"}
			}
		}
	}
}

func TestDatasetLifecycle(t *testing.T) {
	const addr = "truenas_dataset.media"
	lifecycle{
		resourceType: "dataset",
		address:      addr,
		onWrite:      wrapProperties(t),
		create: `
resource "truenas_dataset" "media" {
  name        = "tank/data/media"
  compression = "LZ4"
  quota       = 1099511627776
}`,
		update: `
resource "truenas_dataset" "media" {
  name        = "tank/data/media"
  compression = "ZSTD"
  quota       = 2199023255552
  comments    = "Plex library"
  copies      = "2"
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.StringExact("tank/data/media")),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("compression"), knownvalue.StringExact("LZ4")),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("atime"), knownvalue.StringExact("INHERIT")),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("quota"), knownvalue.Int64Exact(1099511627776)),
			statecheck.ExpectIdentityValue(addr, tfjsonpath.New("id"), knownvalue.StringExact("tank/data/media")),
		},
	}.run(t)
}

func TestZvolLifecycle(t *testing.T) {
	const addr = "truenas_zvol.disk"
	lifecycle{
		resourceType: "zvol",
		address:      addr,
		onWrite:      wrapProperties(t),
		create: `
resource "truenas_zvol" "disk" {
  name    = "tank/vm/disk0"
  volsize = 10737418240
}`,
		update: `
resource "truenas_zvol" "disk" {
  name    = "tank/vm/disk0"
  volsize = 21474836480
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("volsize"), knownvalue.Int64Exact(10737418240)),
		},
	}.run(t)
}

func TestDatasetRefusesToImportZvol(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	store := srv.ServeCRUD("dataset", "zvol")
	store.OnWrite = wrapProperties(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
resource "truenas_zvol" "disk" {
  name    = "tank/vm/disk0"
  volsize = 10737418240
}`,
			},
			{
				Config: providerConfig(srv, "") + `
resource "truenas_zvol" "disk" {
  name    = "tank/vm/disk0"
  volsize = 10737418240
}
import {
  to = truenas_dataset.wrong
  id = "tank/vm/disk0"
}
resource "truenas_dataset" "wrong" {
  name = "tank/vm/disk0"
}`,
				ExpectError: regexp.MustCompile(`(?s)has type VOLUME; this resource manages FILESYSTEM\s+only`),
			},
		},
	})
}
