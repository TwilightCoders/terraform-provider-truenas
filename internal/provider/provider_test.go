package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

var factories = map[string]func() (tfprotov6.ProviderServer, error){
	"truenas": providerserver.NewProtocol6WithError(New("test")()),
}

// providerConfig points the provider at a fake middleware, pinning its certificate.
func providerConfig(srv *middlewaretest.Server, extra string) string {
	sum := sha256.Sum256(srv.Certificate().Raw)
	return fmt.Sprintf(`
provider "truenas" {
  host     = %q
  username = %q
  api_key  = %q
  tls = {
    fingerprint = %q
  }
  %s
}
`, srv.Listener.Addr().String(), middlewaretest.Username, middlewaretest.APIKey, hex.EncodeToString(sum[:]), extra)
}

func TestProviderServesSchema(t *testing.T) {
	server, err := factories["truenas"]()
	if err != nil {
		t.Fatalf("creating provider server: %v", err)
	}
	resp, err := server.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema: %v", err)
	}
	for _, d := range resp.Diagnostics {
		t.Errorf("schema diagnostic: %s: %s", d.Summary, d.Detail)
	}
	if _, ok := resp.ResourceSchemas["truenas_cron_job"]; !ok {
		t.Error("truenas_cron_job missing from provider schema")
	}
}

func TestProviderMetadata(t *testing.T) {
	server, err := providerserver.NewProtocol6WithError(New("1.2.3")())()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := server.GetMetadata(context.Background(), &tfprotov6.GetMetadataRequest{})
	if err != nil || len(resp.Diagnostics) > 0 {
		t.Fatalf("GetMetadata: %v %v", err, resp.Diagnostics)
	}
}

func TestProviderConfigurationErrors(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	t.Setenv(EnvHost, "")
	t.Setenv(EnvUsername, "")
	t.Setenv(EnvAPIKey, "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "truenas" {}
resource "truenas_cron_job" "x" {
  command = "true"
  user    = "root"
}`,
				ExpectError: regexp.MustCompile(`Missing provider configuration`),
			},
			{
				Config: fmt.Sprintf(`
provider "truenas" {
  host     = %q
  username = "root"
  api_key  = "wrong"
  tls = { insecure_skip_verify = true }
}
resource "truenas_cron_job" "x" {
  command = "true"
  user    = "root"
}`, srv.Listener.Addr().String()),
				ExpectError: regexp.MustCompile(`invalid username or API key`),
			},
		},
	})
}
