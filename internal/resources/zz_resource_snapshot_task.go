// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
)

var modelSnapshotTask = &model{
	typeName:     "snapshot_task",
	namespace:    "pool.snapshottask",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "pool.snapshottask.create",
	updateMethod: "pool.snapshottask.update",
	getMethod:    "pool.snapshottask.get_instance",
	deleteMethod: "pool.snapshottask.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "dataset", api: "dataset", path: "dataset",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "The dataset to take snapshots of.",
		},
		{
			name: "recursive", api: "recursive", path: "recursive",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to recursively snapshot child datasets. Defaults to `false`.",
		},
		{
			name: "lifetime_value", api: "lifetime_value", path: "lifetime_value",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "2",
			description: "Number of time units to retain snapshots. `lifetime_unit` gives the time unit. Defaults to `2`.",
		},
		{
			name: "lifetime_unit", api: "lifetime_unit", path: "lifetime_unit",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "WEEK",
			description: "Unit of time for snapshot retention. Defaults to `\"WEEK\"`.",
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether this periodic snapshot task is enabled. Defaults to `true`.",
		},
		{
			name: "exclude", api: "exclude", path: "exclude",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []interface{}{},
			description: "Array of dataset patterns to exclude from recursive snapshots. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "exclude",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "naming_schema", api: "naming_schema", path: "naming_schema",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "auto-%Y-%m-%d_%H-%M",
			description: "Naming pattern for generated snapshots using strftime format. Defaults to `\"auto-%Y-%m-%d_%H-%M\"`.",
		},
		{
			name: "allow_empty", api: "allow_empty", path: "allow_empty",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether to take snapshots even if no data has changed. Defaults to `true`.",
		},
		{
			name: "schedule", api: "schedule", path: "schedule",
			kind: kindObject, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: map[string]interface{}{"begin": "00:00", "dom": "*", "dow": "*", "end": "23:59", "hour": "*", "minute": "00", "month": "*"},
			description: "Cron schedule for when snapshots should be taken. Defaults to `{\"begin\":\"00:00\",\"dom\":\"*\",\"dow\":\"*\",\"end\":\"23:59\",\"hour\":\"*\",\"minute\":\"00\",\"month\":\"*\"}`.",
			children: []*node{
				{
					name: "minute", api: "minute", path: "schedule.minute",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "00",
					description: "Minute when snapshots should be taken (cron format). Defaults to `\"00\"`.",
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
				{
					name: "begin", api: "begin", path: "schedule.begin",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "00:00",
					description: "Start time of the window when snapshots can be taken. Defaults to `\"00:00\"`.",
				},
				{
					name: "end", api: "end", path: "schedule.end",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "23:59",
					description: "End time of the window when snapshots can be taken. Defaults to `\"23:59\"`.",
				},
			},
		},
		{
			name: "vmware_sync", api: "vmware_sync", path: "vmware_sync",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether VMware VMs are synced before taking snapshots.",
		},
	},
}

// NewSnapshotTask returns the snapshot_task resource.
func NewSnapshotTask() resource.Resource {
	return &snapshotTaskResource{crudResource{model: modelSnapshotTask}}
}

type snapshotTaskResource struct{ crudResource }

// NewSnapshotTaskDataSource returns the snapshot_task data source.
func NewSnapshotTaskDataSource() datasource.DataSource {
	return &dataSource{model: modelSnapshotTask, attrs: dataAttrs(modelSnapshotTask)}
}

// NewSnapshotTaskList returns the snapshot_task list resource, for terraform query.
func NewSnapshotTaskList() list.ListResource {
	return &listResource{crudResource{model: modelSnapshotTask}}
}

func (r *snapshotTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Periodic Snapshot Task\n\nCreate a Periodic Snapshot Task that will take snapshots of specified `dataset` at specified `schedule`.\nRecursive snapshots can be created if `recursive` flag is enabled. You can `exclude` specific child datasets\nor zvols from the snapshot.\nSnapshots will be automatically destroyed after a certain amount of time, specified by\n`lifetime_value` and `lifetime_unit`.\nIf multiple periodic tasks create snapshots at the same time (for example hourly and daily at 00:00) the snapshot\nwill be kept until the last of these tasks reaches its expiry time.\nSnapshots will be named according to `naming_schema` which is a `strftime`-like template for snapshot name\nand must contain `%Y`, `%m`, `%d`, `%H` and `%M`.\n\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"dataset": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The dataset to take snapshots of.",
			},
			"recursive": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to recursively snapshot child datasets. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"lifetime_value": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Number of time units to retain snapshots. `lifetime_unit` gives the time unit. Defaults to `2`.",
				Validators:          []validator.Number{numberIsInteger{}},
				Default:             numberdefault.StaticBigFloat(big.NewFloat(2)),
			},
			"lifetime_unit": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Unit of time for snapshot retention. Defaults to `\"WEEK\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("HOUR", "DAY", "WEEK", "MONTH", "YEAR")},
				Default:             stringdefault.StaticString("WEEK"),
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether this periodic snapshot task is enabled. Defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"exclude": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of dataset patterns to exclude from recursive snapshots. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"naming_schema": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Naming pattern for generated snapshots using strftime format. Defaults to `\"auto-%Y-%m-%d_%H-%M\"`.",
				Default:             stringdefault.StaticString("auto-%Y-%m-%d_%H-%M"),
			},
			"allow_empty": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to take snapshots even if no data has changed. Defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"schedule": schema.SingleNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Cron schedule for when snapshots should be taken. Defaults to `{\"begin\":\"00:00\",\"dom\":\"*\",\"dow\":\"*\",\"end\":\"23:59\",\"hour\":\"*\",\"minute\":\"00\",\"month\":\"*\"}`.",
				Attributes: map[string]schema.Attribute{
					"minute": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Minute when snapshots should be taken (cron format). Defaults to `\"00\"`.",
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
					"begin": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Start time of the window when snapshots can be taken. Defaults to `\"00:00\"`.",
						Default:             stringdefault.StaticString("00:00"),
					},
					"end": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "End time of the window when snapshots can be taken. Defaults to `\"23:59\"`.",
						Default:             stringdefault.StaticString("23:59"),
					},
				},
				Default: objectdefault.StaticValue(types.ObjectValueMust(map[string]attr.Type{"minute": types.StringType, "hour": types.StringType, "dom": types.StringType, "month": types.StringType, "dow": types.StringType, "begin": types.StringType, "end": types.StringType}, map[string]attr.Value{"minute": types.StringValue("00"), "hour": types.StringValue("*"), "dom": types.StringValue("*"), "month": types.StringValue("*"), "dow": types.StringValue("*"), "begin": types.StringValue("00:00"), "end": types.StringValue("23:59")})),
			},
			"vmware_sync": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether VMware VMs are synced before taking snapshots.",
			},
		},
	}
}
