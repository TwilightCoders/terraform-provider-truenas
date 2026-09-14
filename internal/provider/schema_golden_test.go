package provider

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var updateGolden = flag.Bool("update", false, "rewrite testdata/schema.golden.json")

const goldenPath = "testdata/schema.golden.json"

// TestSchemaGolden locks every resource schema into a reviewable file. Schema derivation changes
// show up as a diff; run `make schema-golden` to accept them.
func TestSchemaGolden(t *testing.T) {
	server, err := factories["truenas"]()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	identities, err := server.GetResourceIdentitySchemas(context.Background(), &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		t.Fatal(err)
	}

	out := map[string]any{}
	for name, s := range resp.ResourceSchemas {
		entry := map[string]any{"attributes": attributes(s.Block.Attributes)}
		if id, ok := identities.IdentitySchemas[name]; ok {
			ids := map[string]string{}
			for _, a := range id.IdentityAttributes {
				ids[a.Name] = a.Type.String()
			}
			entry["identity"] = ids
		}
		out[name] = entry
	}
	for name, s := range resp.DataSourceSchemas {
		out["data."+name] = map[string]any{"attributes": attributes(s.Block.Attributes)}
	}
	got, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file (run `make schema-golden`): %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("resource schemas changed; review the diff and run `make schema-golden`.\n%s", firstDiff(string(want), string(got)))
	}
}

func attributes(attrs []*tfprotov6.SchemaAttribute) map[string]any {
	out := make(map[string]any, len(attrs))
	for _, a := range attrs {
		entry := map[string]any{}
		switch {
		case a.Required:
			entry["mode"] = "required"
		case a.Optional && a.Computed:
			entry["mode"] = "optional+computed"
		case a.Optional:
			entry["mode"] = "optional"
		default:
			entry["mode"] = "computed"
		}
		if a.Sensitive {
			entry["sensitive"] = true
		}
		if a.WriteOnly {
			entry["write_only"] = true
		}
		if a.NestedType != nil {
			entry["nesting"] = a.NestedType.Nesting.String()
			entry["attributes"] = attributes(a.NestedType.Attributes)
		} else {
			entry["type"] = a.Type.String()
		}
		out[a.Name] = entry
	}
	return out
}

func firstDiff(want, got string) string {
	wl, gl := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := range max(len(wl), len(gl)) {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return "line " + strconv.Itoa(i+1) + ":\n  golden: " + w + "\n  actual: " + g
		}
	}
	return ""
}
