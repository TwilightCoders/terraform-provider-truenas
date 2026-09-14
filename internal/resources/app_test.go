package resources

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestXPortals(t *testing.T) {
	got, err := xPortals(map[string]string{
		"Web UI": "http://plex.lan/",
		"Direct": "http://192.0.2.10:32400/web",
		"Secure": "https://nas.example",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []map[string]any{
		{"name": "Direct", "scheme": "http", "host": "192.0.2.10", "port": 32400, "path": "/web"},
		{"name": "Secure", "scheme": "https", "host": "nas.example", "port": 443, "path": "/"},
		{"name": "Web UI", "scheme": "http", "host": "plex.lan", "port": 80, "path": "/"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("xPortals = %v", got)
	}
	for _, bad := range []string{"ftp://x/", "not a url", "http:///nohost", "http://h:bad/"} {
		if _, err := xPortals(map[string]string{"x": bad}); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

func TestPortalsEqual(t *testing.T) {
	m := func(kv ...string) types.Map {
		elems := map[string]attr.Value{}
		for i := 0; i < len(kv); i += 2 {
			elems[kv[i]] = types.StringValue(kv[i+1])
		}
		v, _ := types.MapValue(types.StringType, elems)
		return v
	}
	if !portalsEqual(m("a", "http://Plex.lan:80"), m("a", "http://plex.lan/")) {
		t.Error("default port and empty path should be equal")
	}
	if portalsEqual(m("a", "http://plex.lan/"), m("a", "https://plex.lan/")) || portalsEqual(m("a", "x"), m("b", "x")) || portalsEqual(m("a", "x"), m()) {
		t.Error("different portals compared equal")
	}
	null := types.MapNull(types.StringType)
	if !portalsEqual(null, null) || portalsEqual(null, m()) || portalsEqual(types.MapUnknown(types.StringType), null) {
		t.Error("null/unknown handling")
	}
	if normalizeURL("::bad") != "::bad" {
		t.Error("unparseable URL should pass through")
	}
}

func TestComposeHelpers(t *testing.T) {
	if !yamlEqual("services:\n  a:\n    image: x\n", map[string]any{"services": map[string]any{"a": map[string]any{"image": "x"}}}) {
		t.Error("equal YAML not recognized")
	}
	if yamlEqual("", map[string]any{}) || yamlEqual(":\n-", map[string]any{}) || yamlEqual("a: 1", map[string]any{"a": 2}) {
		t.Error("unequal YAML recognized")
	}
	if inc, ok := onlyIncludes(map[string]any{"include": []any{"/a.yml", "/b.yml"}}); !ok || len(inc) != 2 {
		t.Error("includes not recognized")
	}
	for _, c := range []map[string]any{{"include": []any{1}}, {"include": "x"}, {"include": []any{"/a"}, "services": nil}} {
		if _, ok := onlyIncludes(c); ok {
			t.Errorf("%v treated as includes", c)
		}
	}

	notes := types.StringValue("hi")
	cfg, diags := composeConfig(context.Background(), appModel{
		Include: types.ListNull(types.StringType), Compose: types.StringValue("services: {}\nx-portals: []\n"),
		Portals: types.MapNull(types.StringType), Notes: notes,
	})
	if diags.HasError() || cfg["x-notes"] != "hi" || cfg["x-portals"] != nil {
		t.Errorf("composeConfig = %v %v", cfg, diags)
	}
	for _, bad := range []string{"- a list", "a: [unclosed"} {
		if _, diags := composeConfig(context.Background(), appModel{Include: types.ListNull(types.StringType), Compose: types.StringValue(bad), Portals: types.MapNull(types.StringType), Notes: types.StringNull()}); !diags.HasError() {
			t.Errorf("invalid compose %q accepted", bad)
		}
	}
}
