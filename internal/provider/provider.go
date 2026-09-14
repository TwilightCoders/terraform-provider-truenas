// Package provider wires the TrueNAS provider into terraform-plugin-framework.
package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/TwilightCoders/terraform-provider-truenas/api"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/apischema"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/engine"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/functions"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/resources"
)

// Address is the registry address Terraform uses to locate this provider.
const Address = "registry.terraform.io/twilightcoders/truenas"

// Environment variables read when the corresponding attribute is not set.
const (
	EnvHost     = "TRUENAS_HOST"
	EnvUsername = "TRUENAS_USERNAME"
	EnvAPIKey   = "TRUENAS_API_KEY"
)

var (
	_ provider.Provider                  = (*Provider)(nil)
	_ provider.ProviderWithListResources = (*Provider)(nil)
	_ provider.ProviderWithFunctions     = (*Provider)(nil)
)

// Provider is the TrueNAS provider.
type Provider struct {
	version  string
	snapshot *apischema.Snapshot
}

// New returns a constructor for the provider at the given version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{version: version, snapshot: apischema.MustLoad(api.Latest)}
	}
}

type config struct {
	Host              types.String `tfsdk:"host"`
	Username          types.String `tfsdk:"username"`
	APIKey            types.String `tfsdk:"api_key"`
	TLS               *tlsConfig   `tfsdk:"tls"`
	ReadOnly          types.Bool   `tfsdk:"read_only"`
	AllowReverseProxy types.Bool   `tfsdk:"allow_reverse_proxy"`
	InsecureLoopback  types.Bool   `tfsdk:"insecure_loopback"`
}

type tlsConfig struct {
	CAPEM              types.String `tfsdk:"ca_pem"`
	Fingerprint        types.String `tfsdk:"fingerprint"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
	ServerName         types.String `tfsdk:"server_name"`
}

func (p *Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "truenas"
	resp.Version = p.version
}

func (p *Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage TrueNAS SCALE through its versioned JSON-RPC API. Requires TrueNAS 25.10 or later.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "TrueNAS HTTPS address with optional port, e.g. `nas.lan` or `nas.lan:8443`. Defaults to `$" + EnvHost + "`. " +
					"Point it at TrueNAS itself: TrueNAS permanently revokes an API key that reaches it over plaintext, " +
					"which happens behind a reverse proxy that terminates TLS and forwards over HTTP.",
			},
			"allow_reverse_proxy": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Connect even when `host` is answered by a web server other than TrueNAS's own. " +
					"The provider refuses by default because a TLS-terminating proxy that forwards over HTTP gets the API key revoked. " +
					"Only set this for a proxy that re-encrypts to TrueNAS.",
			},
			"insecure_loopback": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Speak plaintext HTTP and WebSocket to a loopback address, typically the local end of an SSH tunnel " +
					"to TrueNAS's HTTP port (`ssh -L 18080:127.0.0.1:80 nas`, then `host = \"127.0.0.1:18080\"`). " +
					"TrueNAS accepts API keys over plaintext only from loopback, so the provider refuses any host that resolves elsewhere. " +
					"Cannot be combined with `tls`.",
			},
			"username": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "User that owns the API key. Defaults to `$" + EnvUsername + "`.",
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "API key. Defaults to `$" + EnvAPIKey + "`.",
			},
			"read_only": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Refuse any plan that would create, update or destroy a resource. " +
					"Use it to import and plan against a production box with a guarantee of no writes.",
			},
			"tls": schema.SingleNestedAttribute{
				Optional: true,
				MarkdownDescription: "Server certificate verification. With a publicly trusted certificate (for example an ACME certificate managed by `truenas_certificate`), leave this unset. " +
					"For a private CA, prefer `ca_pem`: it survives certificate renewal. `fingerprint` pins one certificate and breaks when it is renewed.",
				Attributes: map[string]schema.Attribute{
					"ca_pem": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "PEM-encoded CA certificates to trust.",
					},
					"fingerprint": schema.StringAttribute{
						Optional: true,
						MarkdownDescription: "SHA-256 fingerprint of the server certificate, hex with optional colons. Replaces chain and hostname verification. " +
							"Breaks whenever the certificate is renewed; use it for short-lived bootstrap only.",
					},
					"insecure_skip_verify": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Disable certificate verification. Not recommended.",
					},
					"server_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Name to verify the certificate against, when it differs from `host`.",
					},
				},
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg config
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := stringOrEnv(cfg.Host, EnvHost)
	username := stringOrEnv(cfg.Username, EnvUsername)
	apiKey := stringOrEnv(cfg.APIKey, EnvAPIKey)
	for name, value := range map[string]types.String{"host": cfg.Host, "username": cfg.Username, "api_key": cfg.APIKey} {
		if value.IsUnknown() {
			resp.Diagnostics.AddAttributeError(path.Root(name), "Unknown provider configuration",
				"The provider needs "+name+" to be known during planning.")
		}
	}
	for name, value := range map[string]string{"host": host, "username": username, "api_key": apiKey} {
		if value == "" {
			resp.Diagnostics.AddAttributeError(path.Root(name), "Missing provider configuration",
				"Set "+name+" in the provider block or the corresponding TRUENAS_* environment variable.")
		}
	}
	if cfg.InsecureLoopback.ValueBool() && cfg.TLS != nil {
		resp.Diagnostics.AddAttributeError(path.Root("tls"), "Conflicting provider configuration",
			"tls has no effect with insecure_loopback, which never uses TLS. Remove one of them.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	mwCfg := middleware.Config{
		Host:       host,
		APIVersion: p.snapshot.Version,
		Username:   username,
		APIKey:     apiKey,
		Logger:     tflogAdapter{},

		AllowReverseProxy: cfg.AllowReverseProxy.ValueBool(),
		InsecureLoopback:  cfg.InsecureLoopback.ValueBool(),
	}
	if cfg.TLS != nil {
		mwCfg.TLS = middleware.TLSConfig{
			CAPEM:              cfg.TLS.CAPEM.ValueString(),
			Fingerprint:        cfg.TLS.Fingerprint.ValueString(),
			InsecureSkipVerify: cfg.TLS.InsecureSkipVerify.ValueBool(),
			ServerName:         cfg.TLS.ServerName.ValueString(),
		}
	}

	client, err := middleware.Dial(ctx, mwCfg)
	if err != nil {
		resp.Diagnostics.AddError("Unable to connect to TrueNAS", err.Error())
		return
	}
	tflog.Info(ctx, "connected to TrueNAS", map[string]any{"host": host, "api_version": client.APIVersion()})

	data := &engine.ProviderData{Client: client, ReadOnly: cfg.ReadOnly.ValueBool()}
	resp.ResourceData = data
	resp.DataSourceData = data
	resp.ListResourceData = data
}

func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	out := make([]func() resource.Resource, 0, len(resources.All)+len(resources.Custom))
	for _, spec := range resources.All {
		out = append(out, engine.NewResource(p.snapshot, spec))
	}
	return append(out, resources.Custom...)
}

func (p *Provider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{functions.NewSizeBytes}
}

func (p *Provider) ListResources(_ context.Context) []func() list.ListResource {
	var out []func() list.ListResource
	for _, spec := range resources.All {
		if engine.Listable(p.snapshot, spec) {
			out = append(out, engine.NewListResource(p.snapshot, spec))
		}
	}
	return out
}

func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	out := make([]func() datasource.DataSource, 0, len(resources.All))
	for _, spec := range resources.All {
		out = append(out, engine.NewDataSource(p.snapshot, spec))
	}
	return out
}

func stringOrEnv(v types.String, env string) string {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueString()
	}
	return os.Getenv(env)
}

type tflogAdapter struct{}

func (tflogAdapter) Debug(ctx context.Context, msg string, fields map[string]any) {
	tflog.Debug(ctx, msg, fields)
}
