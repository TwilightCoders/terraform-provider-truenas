package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
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

// fakeACLs stores ACLs per path the way filesystem.setacl/getacl do.
type fakeACLs struct {
	mu   sync.Mutex
	acls map[string]map[string]any
}

func serveACLs(srv *middlewaretest.Server, paths ...string) *fakeACLs {
	f := &fakeACLs{acls: map[string]map[string]any{}}
	for _, p := range paths {
		f.acls[p] = map[string]any{"path": p, "acltype": "NFS4", "uid": 0, "gid": 0, "acl": []any{}}
	}
	srv.Handle("filesystem.getacl", func(_ context.Context, params []json.RawMessage) (any, error) {
		var p string
		_ = json.Unmarshal(params[0], &p)
		f.mu.Lock()
		defer f.mu.Unlock()
		acl, ok := f.acls[p]
		if !ok {
			return nil, &middleware.Error{Errno: 2, Errname: "ENOENT", Reason: p + " not found"}
		}
		return acl, nil
	})
	srv.Handle("filesystem.setacl", func(_ context.Context, params []json.RawMessage) (any, error) {
		var req struct {
			Path    string `json:"path"`
			DACL    []any  `json:"dacl"`
			ACLType string `json:"acltype"`
			UID     *int64 `json:"uid"`
			GID     *int64 `json:"gid"`
		}
		if err := json.Unmarshal(params[0], &req); err != nil {
			return nil, err
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		acl, ok := f.acls[req.Path]
		if !ok {
			return nil, &middleware.Error{Errno: 2, Errname: "ENOENT", Reason: req.Path + " not found"}
		}
		acl["acl"] = req.DACL
		if req.UID != nil {
			acl["uid"] = *req.UID
		}
		if req.GID != nil {
			acl["gid"] = *req.GID
		}
		return acl, nil
	})
	return f
}

func (f *fakeACLs) entries(path string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.acls[path]["acl"].([]any))
}

func TestFilesystemACLLifecycle(t *testing.T) {
	const addr = "truenas_filesystem_acl.media"
	srv := middlewaretest.NewServer(t)
	acls := serveACLs(srv, "/mnt/tank/data/media")

	base := `
resource "truenas_filesystem_acl" "media" {
  path = "/mnt/tank/data/media"
  uid  = 3000
  gid  = 3000
  entries = [
    { tag = "owner@", type = "ALLOW", perms = "FULL_CONTROL", flags = "INHERIT" },
    { tag = "group@", type = "ALLOW", perms = "FULL_CONTROL", flags = "INHERIT" },
    { tag = "USER", id = 3001, type = "ALLOW", perms = "READ", flags = "INHERIT" },
    { tag = "everyone@", type = "ALLOW", advanced_perms = [], advanced_flags = [] },
    %s
  ]
}`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		CheckDestroy: func(*terraform.State) error {
			if n := acls.entries("/mnt/tank/data/media"); n == 0 {
				return fmt.Errorf("destroy stripped the ACL")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + fmt.Sprintf(base, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("acltype"), knownvalue.StringExact("NFS4")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("entries"), knownvalue.SetSizeExact(4)),
				},
			},
			{
				Config:           providerConfig(srv, "") + fmt.Sprintf(base, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				ResourceName:                         addr,
				ImportState:                          true,
				ImportStateId:                        "/mnt/tank/data/media",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "path",
			},
			{
				Config: providerConfig(srv, "") + fmt.Sprintf(base, `{ tag = "GROUP", id = 3002, type = "ALLOW", advanced_perms = ["READ_DATA", "EXECUTE"], advanced_flags = ["FILE_INHERIT"] },`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate)},
				},
			},
			{
				Config:           providerConfig(srv, "") + fmt.Sprintf(base, `{ tag = "GROUP", id = 3002, type = "ALLOW", advanced_perms = ["READ_DATA", "EXECUTE"], advanced_flags = ["FILE_INHERIT"] },`),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				Config:      providerConfig(srv, "") + fmt.Sprintf(base, `{ tag = "GROUP", id = 3003, type = "ALLOW", perms = "READ", advanced_perms = ["READ_DATA"], flags = "INHERIT" },`),
				ExpectError: regexp.MustCompile(`set perms or advanced_perms, not both`),
			},
		},
	})
}
