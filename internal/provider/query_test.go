package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

func TestSnapshotTaskQuery(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.ServeCRUD("snapshot_task")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_14_0)},
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
resource "truenas_snapshot_task" "games" {
  dataset = "tank/games"
}
resource "truenas_snapshot_task" "home" {
  dataset    = "tank/home"
  depends_on = [truenas_snapshot_task.games] # fixes id order
}`,
			},
			{
				Query: true,
				Config: `
list "truenas_snapshot_task" "all" {
  provider         = truenas
  include_resource = true
}

list "truenas_snapshot_task" "home" {
  provider = truenas
  config {
    query_filters = jsonencode([["dataset", "=", "tank/home"]])
  }
}`,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("truenas_snapshot_task.all", 2),
					querycheck.ExpectIdentity("truenas_snapshot_task.all", map[string]knownvalue.Check{"id": knownvalue.Int64Exact(1)}),
					querycheck.ExpectIdentity("truenas_snapshot_task.all", map[string]knownvalue.Check{"id": knownvalue.Int64Exact(2)}),
					querycheck.ExpectResourceDisplayName("truenas_snapshot_task.all",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{"id": knownvalue.Int64Exact(2)}),
						knownvalue.StringExact("tank/home")),
					querycheck.ExpectResourceKnownValues("truenas_snapshot_task.all",
						queryfilter.ByResourceIdentity(map[string]knownvalue.Check{"id": knownvalue.Int64Exact(1)}),
						[]querycheck.KnownValueCheck{{Path: tfjsonpath.New("lifetime_unit"), KnownValue: knownvalue.StringExact("WEEK")}}),
					querycheck.ExpectLength("truenas_snapshot_task.home", 1),
				},
			},
		},
	})
}
