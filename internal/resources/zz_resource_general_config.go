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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelGeneralConfig = &model{
	typeName:     "general_config",
	namespace:    "system.general",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	updateMethod: "system.general.update",
	getMethod:    "system.general.config",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "ui_certificate", api: "ui_certificate", path: "ui_certificate",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Used to enable HTTPS access to the system. If `ui_certificate` is not configured on boot, it is automatically     created by the system.",
			ref:         "id",
		},
		{
			name: "ui_httpsport", api: "ui_httpsport", path: "ui_httpsport",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "HTTPS port for the web UI.",
		},
		{
			name: "ui_httpsredirect", api: "ui_httpsredirect", path: "ui_httpsredirect",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "When set, makes sure that all HTTP requests are converted to HTTPS requests to better enhance security.",
		},
		{
			name: "ui_httpsprotocols", api: "ui_httpsprotocols", path: "ui_httpsprotocols",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Array of TLS protocol versions enabled for HTTPS connections.",
			elem: &node{
				name: "", api: "", path: "ui_httpsprotocols",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "ui_port", api: "ui_port", path: "ui_port",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "HTTP port for the web UI.",
		},
		{
			name: "ui_address", api: "ui_address", path: "ui_address",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "A list of valid IPv4 addresses which the system will listen on.",
			elem: &node{
				name: "", api: "", path: "ui_address",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "ui_v6address", api: "ui_v6address", path: "ui_v6address",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "A list of valid IPv6 addresses which the system will listen on.",
			elem: &node{
				name: "", api: "", path: "ui_v6address",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "ui_allowlist", api: "ui_allowlist", path: "ui_allowlist",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "A list of IP addresses and networks that are allow to use API and UI. If this list is empty, then all IP     addresses are allowed to use API and UI.",
			elem: &node{
				name: "", api: "", path: "ui_allowlist",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "ui_consolemsg", api: "ui_consolemsg", path: "ui_consolemsg",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to show console messages on the web UI.",
		},
		{
			name: "ui_x_frame_options", api: "ui_x_frame_options", path: "ui_x_frame_options",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "X-Frame-Options header policy for web UI security.",
		},
		{
			name: "kbdmap", api: "kbdmap", path: "kbdmap",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "System keyboard layout mapping.",
		},
		{
			name: "timezone", api: "timezone", path: "timezone",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "System timezone identifier.",
		},
		{
			name: "usage_collection", api: "usage_collection", path: "usage_collection",
			kind: kindBool, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Whether usage data collection is enabled. `null` if not set.",
		},
		{
			name: "ds_auth", api: "ds_auth", path: "ds_auth",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Controls whether configured Directory Service users that are granted with Privileges are allowed to log in to     the Web UI or use TrueNAS API.",
		},
		{
			name: "ui_restart_delay", api: "ui_restart_delay", path: "ui_restart_delay",
			kind: kindInt, role: roleOptional,
			nullable: true, updatable: true,
			description: "Seconds after the update to restart the web UI and apply UI settings. Without it, UI changes take effect at the next UI restart. The restart aborts every HTTP connection, including the provider's own; reads retry across it. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
		},
		{
			name: "rollback_timeout", api: "rollback_timeout", path: "rollback_timeout",
			kind: kindInt, role: roleOptional,
			nullable: true, updatable: true,
			description: "Timeout in seconds for automatic rollback of UI changes. `null` for no timeout. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
		},
		{
			name: "wizardshown", api: "wizardshown", path: "wizardshown",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether the initial setup wizard has been shown.",
		},
		{
			name: "usage_collection_is_set", api: "usage_collection_is_set", path: "usage_collection_is_set",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether the usage collection preference has been explicitly set.",
		},
		{
			name: "unset", api: "", path: "unset",
			kind: kindList, role: roleOptional,
			description: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
			elem: &node{
				name: "unset", api: "", path: "unset",
				kind: kindString, role: roleRequired,
			},
		},
	},
}

// NewGeneralConfig returns the general_config resource.
func NewGeneralConfig() resource.Resource {
	return &generalConfigResource{crudResource{model: modelGeneralConfig}}
}

type generalConfigResource struct{ crudResource }

// NewGeneralConfigDataSource returns the general_config data source.
func NewGeneralConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelGeneralConfig, attrs: dataAttrs(modelGeneralConfig)}
}

func (r *generalConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update System General Service Configuration.\n\nUI configuration is not applied automatically. Call `system.general.ui_restart` to apply new UI settings (all\nHTTP connections will be aborted) or specify `ui_restart_delay` (in seconds) to automatically apply them after\nsome small amount of time necessary you might need to receive the response for your settings update request.\n\nIf incorrect UI configuration is applied, you might loss API connectivity and won't be able to fix the settings.\nTo avoid that, specify `rollback_timeout` (in seconds). It will automatically roll back UI configuration to the\npreviously working settings after `rollback_timeout` passes unless you call `system.general.checkin` in case\nthe new settings were correct and no rollback is necessary.\n\nThis resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"ui_certificate": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Used to enable HTTPS access to the system. If `ui_certificate` is not configured on boot, it is automatically     created by the system.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"ui_httpsport": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "HTTPS port for the web UI.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"ui_httpsredirect": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "When set, makes sure that all HTTP requests are converted to HTTPS requests to better enhance security.",
			},
			"ui_httpsprotocols": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of TLS protocol versions enabled for HTTPS connections.",
				ElementType:         types.StringType,
			},
			"ui_port": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "HTTP port for the web UI.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"ui_address": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A list of valid IPv4 addresses which the system will listen on.",
				ElementType:         types.StringType,
			},
			"ui_v6address": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A list of valid IPv6 addresses which the system will listen on.",
				ElementType:         types.StringType,
			},
			"ui_allowlist": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A list of IP addresses and networks that are allow to use API and UI. If this list is empty, then all IP     addresses are allowed to use API and UI.",
				ElementType:         types.StringType,
			},
			"ui_consolemsg": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to show console messages on the web UI.",
			},
			"ui_x_frame_options": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "X-Frame-Options header policy for web UI security.",
				Validators:          []validator.String{stringvalidator.OneOf("SAMEORIGIN", "DENY", "ALLOW_ALL")},
			},
			"kbdmap": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "System keyboard layout mapping.",
			},
			"timezone": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "System timezone identifier.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"usage_collection": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether usage data collection is enabled. `null` if not set.",
			},
			"ds_auth": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Controls whether configured Directory Service users that are granted with Privileges are allowed to log in to     the Web UI or use TrueNAS API.",
			},
			"ui_restart_delay": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Seconds after the update to restart the web UI and apply UI settings. Without it, UI changes take effect at the next UI restart. The restart aborts every HTTP connection, including the provider's own; reads retry across it. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"rollback_timeout": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Timeout in seconds for automatic rollback of UI changes. `null` for no timeout. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"wizardshown": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the initial setup wizard has been shown.",
			},
			"usage_collection_is_set": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the usage collection preference has been explicitly set.",
			},
			"unset": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
				ElementType:         types.StringType,
			},
		},
	}
}
