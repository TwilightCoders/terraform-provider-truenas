package provider

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

func TestDataSourceLookupsNeverExposeSecrets(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.ServeCRUD("acme_dns_authenticator")
	datasets := srv.ServeCRUD("dataset")
	datasets.OnWrite = wrapProperties(t)
	srv.ServeConfig("ssh_config", map[string]any{
		"id": json.Number("1"), "tcpport": json.Number("22"), "host_rsa_key": "PRIVATE", "host_rsa_key_pub": "ssh-rsa AAAA",
	})

	resources := providerConfig(srv, "") + `
resource "truenas_acme_dns_authenticator" "cloudflare" {
  name = "cloudflare"
  authenticator = {
    cloudflare = { api_token = "cf-secret" }
  }
}
resource "truenas_dataset" "a" {
  name = "tank/a"
}
resource "truenas_dataset" "b" {
  name       = "tank/b"
  depends_on = [truenas_dataset.a]
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: resources},
			{
				Config: resources + `
data "truenas_acme_dns_authenticator" "cloudflare" {
  name       = "cloudflare"
  depends_on = [truenas_acme_dns_authenticator.cloudflare]
}
data "truenas_dataset" "b" {
  id         = "tank/b"
  depends_on = [truenas_dataset.b]
}
data "truenas_ssh_config" "this" {}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.truenas_acme_dns_authenticator.cloudflare", tfjsonpath.New("id"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.truenas_dataset.b", tfjsonpath.New("name"), knownvalue.StringExact("tank/b")),
					statecheck.ExpectKnownValue("data.truenas_dataset.b", tfjsonpath.New("compression"), knownvalue.StringExact("INHERIT")),
					statecheck.ExpectKnownValue("data.truenas_ssh_config.this", tfjsonpath.New("host_rsa_key_pub"), knownvalue.StringExact("ssh-rsa AAAA")),
				},
			},
			{
				Config: resources + `
data "truenas_dataset" "all" {
  query_filters = jsonencode([["pool", "!=", "nope"]])
}`,
				ExpectError: regexp.MustCompile(`matched 2 objects`),
			},
			{
				Config: resources + `
data "truenas_dataset" "none" {
  name = "tank/missing"
}`,
				ExpectError: regexp.MustCompile(`matched 0 objects`),
			},
			{
				Config: resources + `
data "truenas_dataset" "unset" {}`,
				ExpectError: regexp.MustCompile(`Set id, name or query_filters`),
			},
		},
	})
}

func TestDataSourceSchemasOmitSecrets(t *testing.T) {
	server, err := factories["truenas"]()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := server.GetProviderSchema(t.Context(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.DataSourceSchemas) == 0 {
		t.Fatal("no data sources")
	}
	for name, s := range resp.DataSourceSchemas {
		for _, a := range s.Block.Attributes {
			if a.Sensitive || a.WriteOnly || strings.HasSuffix(a.Name, "_wo_version") {
				t.Errorf("%s exposes %s (sensitive=%v write_only=%v)", name, a.Name, a.Sensitive, a.WriteOnly)
			}
		}
	}
	ssh, ok := resp.DataSourceSchemas["truenas_ssh_config"]
	if !ok {
		t.Fatal("missing ssh config data source")
	}
	for _, a := range ssh.Block.Attributes {
		if a.Name == "host_rsa_key" {
			t.Error("ssh host private key exposed")
		}
	}
}
