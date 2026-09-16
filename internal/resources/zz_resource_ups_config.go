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

var modelUPSConfig = &model{
	typeName:     "ups_config",
	namespace:    "ups",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	updateMethod: "ups.update",
	getMethod:    "ups.config",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "powerdown", api: "powerdown", path: "powerdown",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether the UPS should power down after completing the shutdown sequence.",
		},
		{
			name: "rmonitor", api: "rmonitor", path: "rmonitor",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to enable remote monitoring of the UPS status over the network.",
		},
		{
			name: "nocommwarntime", api: "nocommwarntime", path: "nocommwarntime",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Seconds to wait before warning about communication loss with UPS. `null` for default.",
		},
		{
			name: "remoteport", api: "remoteport", path: "remoteport",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Network port for communicating with remote UPS monitoring systems.",
		},
		{
			name: "shutdowntimer", api: "shutdowntimer", path: "shutdowntimer",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Seconds to wait after initiating shutdown before forcing power off.",
		},
		{
			name: "hostsync", api: "hostsync", path: "hostsync",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum seconds to wait for other systems to shutdown before continuing.",
		},
		{
			name: "description", api: "description", path: "description",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Human-readable description of this UPS configuration.",
		},
		{
			name: "driver", api: "driver", path: "driver",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "UPS driver name that handles communication with the specific UPS hardware model.",
		},
		{
			name: "extrausers", api: "extrausers", path: "extrausers",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional user configurations for UPS monitoring access.",
		},
		{
			name: "identifier", api: "identifier", path: "identifier",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Unique identifier name for this UPS device within the monitoring system.",
		},
		{
			name: "mode", api: "mode", path: "mode",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Operating mode.\n* `MASTER` controls the UPS directly\n* `SLAVE` monitors remotely",
		},
		{
			name: "monpwd", api: "monpwd", path: "monpwd",
			kind: kindString, role: roleOptional,
			sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Password for UPS monitoring authentication (required for updates).",
		},
		{
			name: "monuser", api: "monuser", path: "monuser",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Username for UPS monitoring authentication.",
		},
		{
			name: "options", api: "options", path: "options",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional configuration options passed to the UPS driver.",
		},
		{
			name: "optionsupsd", api: "optionsupsd", path: "optionsupsd",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional configuration options for the UPS daemon.",
		},
		{
			name: "port", api: "port", path: "port",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Serial port or device path for UPS communication.",
		},
		{
			name: "remotehost", api: "remotehost", path: "remotehost",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Hostname or IP address of remote UPS server when operating in SLAVE mode.",
		},
		{
			name: "shutdown", api: "shutdown", path: "shutdown",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Shutdown trigger condition: LOWBATT on low battery, BATT when on battery power.",
		},
		{
			name: "shutdowncmd", api: "shutdowncmd", path: "shutdowncmd",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Custom command to execute during UPS shutdown sequence. `null` for default.",
		},
		{
			name: "complete_identifier", api: "complete_identifier", path: "complete_identifier",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Complete UPS identifier including hostname for network monitoring.",
		},
		{
			name: "monpwd_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `monpwd` again. Terraform never stores `monpwd`.",
		},
	},
}

// NewUPSConfig returns the ups_config resource.
func NewUPSConfig() resource.Resource {
	return &uPSConfigResource{crudResource{model: modelUPSConfig}}
}

type uPSConfigResource struct{ crudResource }

// NewUPSConfigDataSource returns the ups_config data source.
func NewUPSConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelUPSConfig, attrs: dataAttrs(modelUPSConfig)}
}

func (r *uPSConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update UPS Service Configuration.\n\n`powerdown` when enabled, sets UPS to power off after shutting down the system.\n\n`nocommwarntime` is a value in seconds which makes UPS Service wait the specified seconds before alerting that\nthe Service cannot reach configured UPS.\n\n`shutdowntimer` is a value in seconds which tells the Service to wait specified seconds for the UPS before\ninitiating a shutdown. This only applies when `shutdown` is set to \"BATT\".\n\n`shutdowncmd` is the command which is executed to initiate a shutdown. It defaults to \"poweroff\".\n\nThis resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"powerdown": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether the UPS should power down after completing the shutdown sequence.",
			},
			"rmonitor": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable remote monitoring of the UPS status over the network.",
			},
			"nocommwarntime": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Seconds to wait before warning about communication loss with UPS. `null` for default.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"remoteport": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Network port for communicating with remote UPS monitoring systems.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"shutdowntimer": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Seconds to wait after initiating shutdown before forcing power off.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"hostsync": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum seconds to wait for other systems to shutdown before continuing.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}},
			},
			"description": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Human-readable description of this UPS configuration.",
			},
			"driver": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "UPS driver name that handles communication with the specific UPS hardware model.",
			},
			"extrausers": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional user configurations for UPS monitoring access.",
			},
			"identifier": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Unique identifier name for this UPS device within the monitoring system.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"mode": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Operating mode.\n* `MASTER` controls the UPS directly\n* `SLAVE` monitors remotely",
				Validators:          []validator.String{stringvalidator.OneOf("MASTER", "SLAVE")},
			},
			"monpwd": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Password for UPS monitoring authentication (required for updates).",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"monuser": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Username for UPS monitoring authentication.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"options": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional configuration options passed to the UPS driver.",
			},
			"optionsupsd": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional configuration options for the UPS daemon.",
			},
			"port": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Serial port or device path for UPS communication.",
			},
			"remotehost": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Hostname or IP address of remote UPS server when operating in SLAVE mode.",
			},
			"shutdown": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Shutdown trigger condition: LOWBATT on low battery, BATT when on battery power.",
				Validators:          []validator.String{stringvalidator.OneOf("LOWBATT", "BATT")},
			},
			"shutdowncmd": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Custom command to execute during UPS shutdown sequence. `null` for default.",
			},
			"complete_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Complete UPS identifier including hostname for network monitoring.",
			},
			"monpwd_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `monpwd` again. Terraform never stores `monpwd`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
