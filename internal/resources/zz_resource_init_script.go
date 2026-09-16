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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelInitScript = &model{
	typeName:     "init_script",
	namespace:    "initshutdownscript",
	primaryKey:   "id",
	idKind:       kindInt,
	getMethod:    "initshutdownscript.get_instance",
	deleteMethod: "initshutdownscript.delete",
	createMethod: "initshutdownscript.create",
	updateMethod: "initshutdownscript.update",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "type", api: "type", path: "type",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Type of init/shutdown script to execute.\n\n* `COMMAND`: Execute a single command\n* `SCRIPT`: Execute a script file",
		},
		{
			name: "command", api: "command", path: "command",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Must be given if `type=\"COMMAND\"`. Defaults to `\"\"`.",
		},
		{
			name: "script", api: "script", path: "script",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Must be given if `type=\"SCRIPT\"`. Defaults to `\"\"`.",
		},
		{
			name: "when", api: "when", path: "when",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: `* "PREINIT": Early in the boot process before all services have started.
* "POSTINIT": Late in the boot process when most services have started.
* "SHUTDOWN": On shutdown.`,
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether the init/shutdown script is enabled to execute. Defaults to `true`.",
		},
		{
			name: "timeout", api: "timeout", path: "timeout",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: 10,
			description: "An integer time in seconds that the system should wait for the execution of the script/command.\n\nA hard limit for a timeout is configured by the base OS, so when a script/command is set to execute on SHUTDOWN,     the hard limit configured by the base OS is changed adding the timeout specified by script/command so it can be     ensured that it executes as desired and is not interrupted by the base OS's limit. Defaults to `10`.",
		},
		{
			name: "comment", api: "comment", path: "comment",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Optional comment describing the purpose of this script. Defaults to `\"\"`.",
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

// NewInitScript returns the init_script resource.
func NewInitScript() resource.Resource {
	return &initScriptResource{crudResource{model: modelInitScript}}
}

type initScriptResource struct{ crudResource }

// NewInitScriptDataSource returns the init_script data source.
func NewInitScriptDataSource() datasource.DataSource {
	return &dataSource{model: modelInitScript, attrs: dataAttrs(modelInitScript)}
}

// NewInitScriptList returns the init_script list resource, for terraform query.
func NewInitScriptList() list.ListResource {
	return &listResource{crudResource{model: modelInitScript}}
}

func (r *initScriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Runs a command or script at `PREINIT`, `POSTINIT` or `SHUTDOWN`. Set `command` when `type` is `COMMAND` and `script` when it is `SCRIPT`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of init/shutdown script to execute.\n\n* `COMMAND`: Execute a single command\n* `SCRIPT`: Execute a script file",
				Validators:          []validator.String{stringvalidator.OneOf("COMMAND", "SCRIPT")},
			},
			"command": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Must be given if `type=\"COMMAND\"`. Defaults to `\"\"`.",
			},
			"script": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Must be given if `type=\"SCRIPT\"`. Defaults to `\"\"`.",
			},
			"when": schema.StringAttribute{
				Required: true,
				MarkdownDescription: `* "PREINIT": Early in the boot process before all services have started.
* "POSTINIT": Late in the boot process when most services have started.
* "SHUTDOWN": On shutdown.`,
				Validators: []validator.String{stringvalidator.OneOf("PREINIT", "POSTINIT", "SHUTDOWN")},
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether the init/shutdown script is enabled to execute. Defaults to `true`.",
			},
			"timeout": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "An integer time in seconds that the system should wait for the execution of the script/command.\n\nA hard limit for a timeout is configured by the base OS, so when a script/command is set to execute on SHUTDOWN,     the hard limit configured by the base OS is changed adding the timeout specified by script/command so it can be     ensured that it executes as desired and is not interrupted by the base OS's limit. Defaults to `10`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"comment": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Optional comment describing the purpose of this script. Defaults to `\"\"`.",
				Validators:          []validator.String{stringvalidator.LengthAtMost(255)},
			},
			"unset": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
				ElementType:         types.StringType,
			},
		},
	}
}
