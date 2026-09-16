// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"math/big"
)

var modelAcmeDnsAuthenticator = &model{
	typeName:     "acme_dns_authenticator",
	namespace:    "acme.dns.authenticator",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "acme.dns.authenticator.create",
	updateMethod: "acme.dns.authenticator.update",
	getMethod:    "acme.dns.authenticator.get_instance",
	deleteMethod: "acme.dns.authenticator.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "authenticator", api: "attributes", path: "attributes",
			kind: kindUnion, role: roleRequired,
			readable: true, updatable: true,
			description:   "Authentication credentials and configuration for the DNS provider.",
			discriminator: "authenticator",
			children: []*node{
				{
					name: "cloudflare", api: "cloudflare", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `authenticator` is `cloudflare`.",
					children: []*node{
						{
							name: "cloudflare_email", api: "cloudflare_email", path: "attributes.cloudflare_email",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "Cloudflare Email.",
						},
						{
							name: "api_key", api: "api_key", path: "attributes.api_key",
							kind: kindString, role: roleOptional,
							nullable: true, sensitive: true, writeOnly: true, readable: true,
							description: "API Key.",
						},
						{
							name: "api_token", api: "api_token", path: "attributes.api_token",
							kind: kindString, role: roleOptional,
							nullable: true, sensitive: true, writeOnly: true, readable: true,
							description: "API Token.",
						},
					},
				},
				{
					name: "digitalocean", api: "digitalocean", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `authenticator` is `digitalocean`.",
					children: []*node{
						{
							name: "digitalocean_token", api: "digitalocean_token", path: "attributes.digitalocean_token",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "DigitalOcean Token.",
						},
					},
				},
				{
					name: "ovh", api: "OVH", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `authenticator` is `OVH`.",
					children: []*node{
						{
							name: "application_key", api: "application_key", path: "attributes.application_key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OVH Application Key.",
						},
						{
							name: "application_secret", api: "application_secret", path: "attributes.application_secret",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OVH Application Secret.",
						},
						{
							name: "consumer_key", api: "consumer_key", path: "attributes.consumer_key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "OVH Consumer Key.",
						},
						{
							name: "endpoint", api: "endpoint", path: "attributes.endpoint",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "OVH Endpoint.",
						},
					},
				},
				{
					name: "route53", api: "route53", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `authenticator` is `route53`.",
					children: []*node{
						{
							name: "access_key_id", api: "access_key_id", path: "attributes.access_key_id",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "AWS Access Key ID.",
						},
						{
							name: "secret_access_key", api: "secret_access_key", path: "attributes.secret_access_key",
							kind: kindString, role: roleRequired,
							sensitive: true, writeOnly: true, readable: true,
							description: "AWS Secret Access Key.",
						},
					},
				},
				{
					name: "shell", api: "shell", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `authenticator` is `shell`.",
					children: []*node{
						{
							name: "script", api: "script", path: "attributes.script",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Authentication Script.",
						},
						{
							name: "user", api: "user", path: "attributes.user",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "nobody",
							description: "Running user. Defaults to `\"nobody\"`.",
						},
						{
							name: "timeout", api: "timeout", path: "attributes.timeout",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "60",
							description: "Script Timeout. Defaults to `60`.",
						},
						{
							name: "delay", api: "delay", path: "attributes.delay",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "60",
							description: "Propagation delay. Defaults to `60`.",
						},
					},
				},
			},
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Human-readable name for the DNS authenticator.",
		},
		{
			name: "authenticator_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `authenticator` again. Terraform never stores `authenticator`.",
		},
	},
}

// NewAcmeDnsAuthenticator returns the acme_dns_authenticator resource.
func NewAcmeDnsAuthenticator() resource.Resource {
	return &acmeDnsAuthenticatorResource{crudResource{model: modelAcmeDnsAuthenticator}}
}

type acmeDnsAuthenticatorResource struct{ crudResource }

// NewAcmeDnsAuthenticatorDataSource returns the acme_dns_authenticator data source.
func NewAcmeDnsAuthenticatorDataSource() datasource.DataSource {
	return &dataSource{model: modelAcmeDnsAuthenticator, attrs: dataAttrs(modelAcmeDnsAuthenticator)}
}

// NewAcmeDnsAuthenticatorList returns the acme_dns_authenticator list resource, for terraform query.
func NewAcmeDnsAuthenticatorList() list.ListResource {
	return &listResource{crudResource{model: modelAcmeDnsAuthenticator}}
}

func (r *acmeDnsAuthenticatorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a DNS provider for ACME DNS-01 challenges. Set exactly one provider under `authenticator`. Credentials are write-only: they never enter Terraform state. Bump `authenticator_wo_version` to resend them.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"authenticator": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Authentication credentials and configuration for the DNS provider.",
				Attributes: map[string]schema.Attribute{
					"cloudflare": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `authenticator` is `cloudflare`.",
						Attributes: map[string]schema.Attribute{
							"cloudflare_email": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Cloudflare Email.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"api_key": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "API Key.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"api_token": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "API Token.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"digitalocean": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `authenticator` is `digitalocean`.",
						Attributes: map[string]schema.Attribute{
							"digitalocean_token": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "DigitalOcean Token.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"ovh": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `authenticator` is `OVH`.",
						Attributes: map[string]schema.Attribute{
							"application_key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OVH Application Key.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"application_secret": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OVH Application Secret.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"consumer_key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "OVH Consumer Key.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"endpoint": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "OVH Endpoint.",
								Validators:          []validator.String{stringvalidator.OneOf("ovh-eu", "ovh-ca", "ovh-us", "kimsufi-eu", "kimsufi-ca", "soyoustart-eu", "soyoustart-ca")},
							},
						},
					},
					"route53": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `authenticator` is `route53`.",
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "AWS Access Key ID.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"secret_access_key": schema.StringAttribute{
								Required: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "AWS Secret Access Key.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"shell": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `authenticator` is `shell`.",
						Attributes: map[string]schema.Attribute{
							"script": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Authentication Script.",
							},
							"user": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Running user. Defaults to `\"nobody\"`.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
								Default:             stringdefault.StaticString("nobody"),
							},
							"timeout": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Script Timeout. Defaults to `60`.",
								Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 5}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(60)),
							},
							"delay": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Propagation delay. Defaults to `60`.",
								Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 10}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(60)),
							},
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable name for the DNS authenticator.",
			},
			"authenticator_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `authenticator` again. Terraform never stores `authenticator`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
