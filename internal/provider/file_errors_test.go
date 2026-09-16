package provider

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

// TestFileVanishedIsRecreated covers a file deleted outside Terraform. Refresh has to notice it is
// gone and plan to put it back, rather than failing or reporting it as present.
func TestFileVanishedIsRecreated(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)
	cfg := fileConfig(srv, "keep\n", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				PreConfig: func() { srv.Files.Delete("/mnt/pool/apps/compose.yml") },
				Config:    cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("truenas_file.compose", plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

// TestFileReportsWriteFailure covers the path taken when the server refuses the upload: the error
// has to reach the operator rather than the apply reporting success.
func TestFileReportsWriteFailure(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	serveStat(srv)
	srv.Handle("auth.generate_token", func(context.Context, []json.RawMessage) (any, error) {
		return nil, &middleware.Error{Errno: 13, Errname: "EACCES", Reason: "not permitted"}
	})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{{
			Config:      fileConfig(srv, "nope\n", ""),
			ExpectError: regexp.MustCompile(`Cannot write file`),
		}},
	})
}
