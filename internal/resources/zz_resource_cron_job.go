// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelCronJob = &model{
	typeName:     "cron_job",
	namespace:    "cronjob",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "cronjob.create",
	updateMethod: "cronjob.update",
	getMethod:    "cronjob.get_instance",
	deleteMethod: "cronjob.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether the cron job is active and will be executed. Defaults to `true`.",
		},
		{
			name: "ignore_stderr", api: "stderr", path: "stderr",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to IGNORE standard error (if `false`, it will be added to email). Defaults to `false`.",
		},
		{
			name: "ignore_stdout", api: "stdout", path: "stdout",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether to IGNORE standard output (if `false`, it will be added to email). Defaults to `true`.",
		},
		{
			name: "schedule", api: "schedule", path: "schedule",
			kind: kindObject, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: map[string]any{"dom": "*", "dow": "*", "hour": "*", "minute": "00", "month": "*"},
			description: "Cron schedule configuration for when the job runs. Defaults to `{\"dom\":\"*\",\"dow\":\"*\",\"hour\":\"*\",\"minute\":\"00\",\"month\":\"*\"}`.",
			children: []*node{
				{
					name: "minute", api: "minute", path: "schedule.minute",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "00",
					description: "\"00\" - \"59\" Defaults to `\"00\"`.",
				},
				{
					name: "hour", api: "hour", path: "schedule.hour",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"00\" - \"23\" Defaults to `\"*\"`.",
				},
				{
					name: "dom", api: "dom", path: "schedule.dom",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"1\" - \"31\" Defaults to `\"*\"`.",
				},
				{
					name: "month", api: "month", path: "schedule.month",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"1\" (January) - \"12\" (December) Defaults to `\"*\"`.",
				},
				{
					name: "dow", api: "dow", path: "schedule.dow",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"1\" (Monday) - \"7\" (Sunday) Defaults to `\"*\"`.",
				},
			},
		},
		{
			name: "command", api: "command", path: "command",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Shell command or script to execute.",
		},
		{
			name: "description", api: "description", path: "description",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Human-readable description of what this cron job does. Defaults to `\"\"`.",
		},
		{
			name: "user", api: "user", path: "user",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "System user account to run the command as.",
		},
	},
}

// NewCronJob returns the cron_job resource.
func NewCronJob() resource.Resource {
	return &cronJobResource{crudResource{model: modelCronJob}}
}

type cronJobResource struct{ crudResource }

// NewCronJobDataSource returns the cron_job data source.
func NewCronJobDataSource() datasource.DataSource {
	return &dataSource{model: modelCronJob, attrs: dataAttrs(modelCronJob)}
}

// NewCronJobList returns the cron_job list resource, for terraform query.
func NewCronJobList() list.ListResource {
	return &listResource{crudResource{model: modelCronJob}}
}

func (r *cronJobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new cron job.\n\n`stderr` and `stdout` are boolean values which if `true`, represent that we would like to suppress\nstandard error / standard output respectively.\n\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether the cron job is active and will be executed. Defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"ignore_stderr": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to IGNORE standard error (if `false`, it will be added to email). Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"ignore_stdout": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to IGNORE standard output (if `false`, it will be added to email). Defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"schedule": schema.SingleNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Cron schedule configuration for when the job runs. Defaults to `{\"dom\":\"*\",\"dow\":\"*\",\"hour\":\"*\",\"minute\":\"00\",\"month\":\"*\"}`.",
				Attributes: map[string]schema.Attribute{
					"minute": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "\"00\" - \"59\" Defaults to `\"00\"`.",
						Default:             stringdefault.StaticString("00"),
					},
					"hour": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "\"00\" - \"23\" Defaults to `\"*\"`.",
						Default:             stringdefault.StaticString("*"),
					},
					"dom": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "\"1\" - \"31\" Defaults to `\"*\"`.",
						Default:             stringdefault.StaticString("*"),
					},
					"month": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "\"1\" (January) - \"12\" (December) Defaults to `\"*\"`.",
						Default:             stringdefault.StaticString("*"),
					},
					"dow": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "\"1\" (Monday) - \"7\" (Sunday) Defaults to `\"*\"`.",
						Default:             stringdefault.StaticString("*"),
					},
				},
				Default: objectdefault.StaticValue(types.ObjectValueMust(map[string]attr.Type{"minute": types.StringType, "hour": types.StringType, "dom": types.StringType, "month": types.StringType, "dow": types.StringType}, map[string]attr.Value{"minute": types.StringValue("00"), "hour": types.StringValue("*"), "dom": types.StringValue("*"), "month": types.StringValue("*"), "dow": types.StringValue("*")})),
			},
			"command": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Shell command or script to execute.",
			},
			"description": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Human-readable description of what this cron job does. Defaults to `\"\"`.",
				Default:             stringdefault.StaticString(""),
			},
			"user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "System user account to run the command as.",
			},
		},
	}
}
