package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/TwilightCoders/terraform-provider-truenas/api"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/apischema"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

func TestCertificateLifecycle(t *testing.T) {
	const addr = "truenas_certificate.lan"
	lifecycle{
		namespace: "certificate",
		address:   addr,
		// Like middlewared, report the private key back on read.
		onWrite: func(row map[string]any) {
			row["privatekey"] = "-----BEGIN PRIVATE KEY-----generated"
		},
		create: `
resource "truenas_certificate" "lan" {
  name        = "truenas-lan"
  create_type = "CERTIFICATE_CREATE_CSR"
  common      = "nas.example.com"
  san         = ["nas.example.com"]
}`,
		update: `
resource "truenas_certificate" "lan" {
  name        = "truenas-lan"
  create_type = "CERTIFICATE_CREATE_CSR"
  common      = "nas.example.com"
  san         = ["nas.example.com"]
  renew_days  = 30
}`,
		checks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("privatekey"), knownvalue.Null()),
			statecheck.ExpectKnownValue(addr, tfjsonpath.New("renew_days"), knownvalue.Int64Exact(10)),
		},
		// TrueNAS never reports how a certificate was created.
		importIgnore:  []string{"create_type"},
		importUpdates: true,
	}.run(t)
}

func TestACMEDNSAuthenticatorLifecycle(t *testing.T) {
	lifecycle{
		namespace: "acme.dns.authenticator",
		address:   "truenas_acme_dns_authenticator.cloudflare",
		create: `
resource "truenas_acme_dns_authenticator" "cloudflare" {
  name = "cloudflare"
  authenticator = {
    cloudflare = {
      api_token = "cf-token"
    }
  }
}`,
		update: `
resource "truenas_acme_dns_authenticator" "cloudflare" {
  name = "cloudflare-dns"
  authenticator = {
    cloudflare = {
      api_token = "cf-token-rotated"
    }
  }
}`,
	}.run(t)
}

func TestGeneralConfigUICertificate(t *testing.T) {
	const addr = "truenas_general_config.this"
	srv := middlewaretest.NewServer(t)
	srv.ServeConfig(apischema.MustLoad(api.Latest), "system.general", map[string]any{
		"id": json.Number("1"), "ui_httpsport": json.Number("8443"), "ui_port": json.Number("8080"),
		"ui_address": []any{"0.0.0.0"}, "timezone": "UTC",
		"ui_certificate": map[string]any{"id": json.Number("1"), "name": "truenas_default"},
	})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(srv, "") + `
resource "truenas_general_config" "this" {
  ui_httpsport = 8444
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("ui_certificate"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("ui_httpsport"), knownvalue.Int64Exact(8444)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("timezone"), knownvalue.StringExact("UTC")),
				},
			},
		},
	})
}

func TestCreateOnlyChangeReplaces(t *testing.T) {
	const addr = "truenas_certificate.lan"
	srv := middlewaretest.NewServer(t)
	srv.ServeCRUD(apischema.MustLoad(api.Latest), "certificate")
	config := func(createType string) string {
		return providerConfig(srv, "") + `
resource "truenas_certificate" "lan" {
  name        = "truenas-lan"
  create_type = "` + createType + `"
  common      = "nas.example.com"
}`
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{Config: config("CERTIFICATE_CREATE_CSR")},
			{
				Config: config("CERTIFICATE_CREATE_IMPORTED_CSR"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionDestroyBeforeCreate)},
				},
			},
		},
	})
}
