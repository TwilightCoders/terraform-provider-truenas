package resources

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sample_media_acl.json is filesystem.getacl output for a real NFS4 dataset.
func TestDecodeRealNFS4ACL(t *testing.T) {
	raw, err := os.ReadFile("testdata/sample_media_acl.json")
	if err != nil {
		t.Fatal(err)
	}
	var acl wireACL
	if err := json.Unmarshal(raw, &acl); err != nil {
		t.Fatal(err)
	}
	if acl.ACLType != "NFS4" || len(acl.ACL) == 0 {
		t.Fatalf("unexpected fixture: %+v", acl)
	}
	var sawBasic, sawAdvanced bool
	for _, w := range acl.ACL {
		e, err := entryFromAPI(w, acl.ACLType)
		if err != nil {
			t.Fatalf("%s: %v", w.Tag, err)
		}
		if !e.Perms.IsNull() {
			sawBasic = true
		}
		if !e.AdvancedPerms.IsNull() {
			sawAdvanced = true
		}
		back, err := entryToAPI(e, acl.ACLType)
		if err != nil {
			t.Fatalf("%s: re-encode: %v", w.Tag, err)
		}
		var original map[string]any
		_ = json.Unmarshal(w.Perms, &original)
		gotPerms, _ := json.Marshal(back["perms"])
		wantPerms, _ := json.Marshal(original)
		if string(gotPerms) != string(wantPerms) {
			t.Errorf("%s perms round trip: %s != %s", w.Tag, gotPerms, wantPerms)
		}
	}
	if !sawBasic || !sawAdvanced {
		t.Errorf("fixture should exercise basic and advanced perms (basic=%v advanced=%v)", sawBasic, sawAdvanced)
	}
}

func TestPOSIXEntries(t *testing.T) {
	e := aclEntry{
		Tag: types.StringValue("USER_OBJ"), ID: types.Int64Null(),
		Read: types.BoolValue(true), Write: types.BoolValue(true), Execute: types.BoolValue(false), Default: types.BoolValue(false),
	}
	wire, err := entryToAPI(e, "POSIX1E")
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(wire)
	var w wireEntry
	_ = json.Unmarshal(raw, &w)
	back, err := entryFromAPI(w, "POSIX1E")
	if err != nil || !back.Read.ValueBool() || back.Execute.ValueBool() {
		t.Fatalf("POSIX round trip: %+v %v", back, err)
	}
	e.Read = types.BoolNull()
	if _, err := entryToAPI(e, "POSIX1E"); err == nil {
		t.Error("incomplete POSIX entry accepted")
	}
}
