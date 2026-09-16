// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var modelSnmpConfig = &model{
	typeName:     "snmp_config",
	namespace:    "snmp",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	updateMethod: "snmp.update",
	getMethod:    "snmp.config",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "location", api: "location", path: "location",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "A comment describing the physical location of the server.",
		},
		{
			name: "contact", api: "contact", path: "contact",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Contact information for the system administrator (email or name).",
		},
		{
			name: "traps", api: "traps", path: "traps",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether SNMP traps are enabled.",
		},
		{
			name: "v3", api: "v3", path: "v3",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether SNMP version 3 is enabled.  Enabling version 3 also requires username, authtype and password.",
		},
		{
			name: "community", api: "community", path: "community",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "SNMP community string for v1/v2c access. Allows         letters and numbers: a-zA-Z0-9          special characters: !$%&()+-_={}[]<>,.?          and spaces. Notable excluded characters: # / \\ @",
		},
		{
			name: "v3_username", api: "v3_username", path: "v3_username",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Username for SNMP version 3 authentication.",
		},
		{
			name: "v3_authtype", api: "v3_authtype", path: "v3_authtype",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Authentication type for SNMP version 3 (empty string means no authentication).",
		},
		{
			name: "v3_password", api: "v3_password", path: "v3_password",
			kind: kindString, role: roleOptional,
			sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Password for SNMP version 3 authentication.",
		},
		{
			name: "v3_privproto", api: "v3_privproto", path: "v3_privproto",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Privacy protocol for SNMP version 3 encryption. `null` means no encryption.      If set, ['AES'|'DES'], a `privpassphrase` must be supplied.",
		},
		{
			name: "v3_privpassphrase", api: "v3_privpassphrase", path: "v3_privpassphrase",
			kind: kindString, role: roleOptional,
			nullable: true, sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Privacy passphrase for SNMP version 3 encryption. This field is required when `privproto` is set.",
		},
		{
			name: "loglevel", api: "loglevel", path: "loglevel",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Logging level for SNMP daemon (0=emergency to 7=debug).",
		},
		{
			name: "options", api: "options", path: "options",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional SNMP daemon configuration options.     Manual settings should be used with caution as they may render the SNMP service non-functional.",
		},
		{
			name: "zilstat", api: "zilstat", path: "zilstat",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to enable ZFS dataset statistics collection for SNMP.",
		},
		{
			name: "v3_password_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `v3_password` again. Terraform never stores `v3_password`.",
		},
		{
			name: "v3_privpassphrase_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `v3_privpassphrase` again. Terraform never stores `v3_privpassphrase`.",
		},
	},
}

// NewSnmpConfig returns the snmp_config resource.
func NewSnmpConfig() resource.Resource {
	return &snmpConfigResource{crudResource{model: modelSnmpConfig}}
}

type snmpConfigResource struct{ crudResource }

// NewSnmpConfigDataSource returns the snmp_config data source.
func NewSnmpConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelSnmpConfig, attrs: dataAttrs(modelSnmpConfig)}
}

func (r *snmpConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Update SNMP Service Configuration.

--- Rules ---
Enabling v3:
    requires v3_username, v3_authtype and v3_password
Disabling v3:
    By itself will retain the v3 user settings and config in the 'private' config,
    but remove the entry in the public config to block v3 access by that user.
Disabling v3 and clearing the v3_username:
    This will do the actions described in 'Disabling v3' and take the extra step to
    remove the user from the 'private' config.

The 'v3_*' settings are valid and enforced only when 'v3' is enabled

This resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"location": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A comment describing the physical location of the server.",
			},
			"contact": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Contact information for the system administrator (email or name).",
			},
			"traps": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether SNMP traps are enabled.",
			},
			"v3": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether SNMP version 3 is enabled.  Enabling version 3 also requires username, authtype and password.",
			},
			"community": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "SNMP community string for v1/v2c access. Allows         letters and numbers: a-zA-Z0-9          special characters: !$%&()+-_={}[]<>,.?          and spaces. Notable excluded characters: # / \\ @",
			},
			"v3_username": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Username for SNMP version 3 authentication.",
				Validators:          []validator.String{stringvalidator.LengthAtMost(20)},
			},
			"v3_authtype": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Authentication type for SNMP version 3 (empty string means no authentication).",
				Validators:          []validator.String{stringvalidator.OneOf("", "MD5", "SHA")},
			},
			"v3_password": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Password for SNMP version 3 authentication.",
			},
			"v3_privproto": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Privacy protocol for SNMP version 3 encryption. `null` means no encryption.      If set, ['AES'|'DES'], a `privpassphrase` must be supplied.",
				Validators:          []validator.String{stringvalidator.OneOf("AES", "DES")},
			},
			"v3_privpassphrase": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Privacy passphrase for SNMP version 3 encryption. This field is required when `privproto` is set.",
			},
			"loglevel": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Logging level for SNMP daemon (0=emergency to 7=debug).",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 7}},
			},
			"options": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional SNMP daemon configuration options.     Manual settings should be used with caution as they may render the SNMP service non-functional.",
			},
			"zilstat": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable ZFS dataset statistics collection for SNMP.",
			},
			"v3_password_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `v3_password` again. Terraform never stores `v3_password`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"v3_privpassphrase_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `v3_privpassphrase` again. Terraform never stores `v3_privpassphrase`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
