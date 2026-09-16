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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelReplication = &model{
	typeName:     "replication",
	namespace:    "replication",
	primaryKey:   "id",
	idKind:       kindInt,
	updateMethod: "replication.update",
	getMethod:    "replication.get_instance",
	deleteMethod: "replication.delete",
	createMethod: "replication.create",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Name for replication task.",
		},
		{
			name: "direction", api: "direction", path: "direction",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Whether task will `PUSH` or `PULL` snapshots.",
		},
		{
			name: "transport", api: "transport", path: "transport",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Method of snapshots transfer.\n\n* `SSH` transfers snapshots via SSH connection. This method is supported everywhere but does not achieve       great performance.\n* `SSH+NETCAT` uses unencrypted connection for data transfer. This can only be used in trusted networks       and requires a port (specified by range from `netcat_active_side_port_min` to `netcat_active_side_port_max`)       to be open on `netcat_active_side`.\n* `LOCAL` replicates to or from localhost.",
		},
		{
			name: "ssh_credentials", api: "ssh_credentials", path: "ssh_credentials",
			kind: kindInt, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Keychain Credential ID of type `SSH_CREDENTIALS`.",
			ref:         "id",
		},
		{
			name: "netcat_active_side", api: "netcat_active_side", path: "netcat_active_side",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Which side actively establishes the netcat connection for `SSH+NETCAT` transport.\n\n* `LOCAL`: Local system initiates the connection\n* `REMOTE`: Remote system initiates the connection\n* `null`: Not applicable for other transport types",
		},
		{
			name: "netcat_active_side_listen_address", api: "netcat_active_side_listen_address", path: "netcat_active_side_listen_address",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "IP address for the active side to listen on for `SSH+NETCAT` transport. `null` if not applicable.",
		},
		{
			name: "netcat_active_side_port_min", api: "netcat_active_side_port_min", path: "netcat_active_side_port_min",
			kind: kindInt, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Minimum port number in the range for netcat connections. `null` if not applicable.",
		},
		{
			name: "netcat_active_side_port_max", api: "netcat_active_side_port_max", path: "netcat_active_side_port_max",
			kind: kindInt, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Maximum port number in the range for netcat connections. `null` if not applicable.",
		},
		{
			name: "netcat_passive_side_connect_address", api: "netcat_passive_side_connect_address", path: "netcat_passive_side_connect_address",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "IP address for the passive side to connect to for `SSH+NETCAT` transport. `null` if not applicable.",
		},
		{
			name: "sudo", api: "sudo", path: "sudo",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "`SSH` and `SSH+NETCAT` transports should use sudo (which is expected to be passwordless) to run `zfs`     command on the remote machine. Defaults to `false`.",
		},
		{
			name: "source_datasets", api: "source_datasets", path: "source_datasets",
			kind: kindList, role: roleRequired,
			readable: true, updatable: true,
			description: "List of datasets to replicate snapshots from.",
			elem: &node{
				name: "", api: "", path: "source_datasets",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "target_dataset", api: "target_dataset", path: "target_dataset",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Dataset to put snapshots into.",
		},
		{
			name: "recursive", api: "recursive", path: "recursive",
			kind: kindBool, role: roleRequired,
			readable: true, updatable: true,
			description: "Whether to recursively replicate child datasets.",
		},
		{
			name: "exclude", api: "exclude", path: "exclude",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "Array of dataset patterns to exclude from replication. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "exclude",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "properties", api: "properties", path: "properties",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Send dataset properties along with snapshots. Defaults to `true`.",
		},
		{
			name: "properties_exclude", api: "properties_exclude", path: "properties_exclude",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "Array of dataset property names to exclude from replication. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "properties_exclude",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "properties_override", api: "properties_override", path: "properties_override",
			kind: kindMap, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: map[string]any{},
			description: "Object mapping dataset property names to override values during replication. Defaults to `{}`.",
			elem: &node{
				name: "", api: "", path: "properties_override",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "replicate", api: "replicate", path: "replicate",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to use full ZFS replication. Defaults to `false`.",
		},
		{
			name: "encryption", api: "encryption", path: "encryption",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to enable encryption for the replicated datasets. Defaults to `false`.",
		},
		{
			name: "encryption_inherit", api: "encryption_inherit", path: "encryption_inherit",
			kind: kindBool, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Whether replicated datasets should inherit encryption from parent. `null` if encryption is disabled.",
		},
		{
			name: "encryption_key", api: "encryption_key", path: "encryption_key",
			kind: kindString, role: roleOptional,
			nullable: true, sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Encryption key for replicated datasets. `null` if not specified.",
		},
		{
			name: "encryption_key_format", api: "encryption_key_format", path: "encryption_key_format",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Format of the encryption key.\n\n* `HEX`: Hexadecimal-encoded key\n* `PASSPHRASE`: Text passphrase\n* `null`: Not applicable when encryption is disabled",
		},
		{
			name: "encryption_key_location", api: "encryption_key_location", path: "encryption_key_location",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Filesystem path where encryption key is stored. `null` if not using key file.",
		},
		{
			name: "periodic_snapshot_tasks", api: "periodic_snapshot_tasks", path: "periodic_snapshot_tasks",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "List of periodic snapshot task IDs that are sources of snapshots for this replication task. Only push     replication tasks can be bound to periodic snapshot tasks. Defaults to `[]`.",
			ref:         "id",
			elem: &node{
				name: "", api: "", path: "periodic_snapshot_tasks",
				kind: kindInt, role: roleRequired,
				readable: true,
				ref:      "id",
			},
		},
		{
			name: "naming_schema", api: "naming_schema", path: "naming_schema",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "List of naming schemas for pull replication. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "naming_schema",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "also_include_naming_schema", api: "also_include_naming_schema", path: "also_include_naming_schema",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "List of naming schemas for push replication. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "also_include_naming_schema",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "name_regex", api: "name_regex", path: "name_regex",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Replicate all snapshots which names match specified regular expression.",
		},
		{
			name: "auto", api: "auto", path: "auto",
			kind: kindBool, role: roleRequired,
			readable: true, updatable: true,
			description: "Allow replication to run automatically on schedule or after bound periodic snapshot task.",
		},
		{
			name: "schedule", api: "schedule", path: "schedule",
			kind: kindObject, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Schedule to run replication task. Only `auto` replication tasks without bound periodic snapshot tasks can have     a schedule.",
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
				{
					name: "begin", api: "begin", path: "schedule.begin",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "00:00",
					description: "Start time for the time window in HH:MM format. Defaults to `\"00:00\"`.",
				},
				{
					name: "end", api: "end", path: "schedule.end",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "23:59",
					description: "End time for the time window in HH:MM format. Defaults to `\"23:59\"`.",
				},
			},
		},
		{
			name: "restrict_schedule", api: "restrict_schedule", path: "restrict_schedule",
			kind: kindObject, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Restricts when replication task with bound periodic snapshot tasks runs. For example, you can have periodic     snapshot tasks that run every 15 minutes, but only run replication task every hour.",
			children: []*node{
				{
					name: "minute", api: "minute", path: "restrict_schedule.minute",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "00",
					description: "\"00\" - \"59\" Defaults to `\"00\"`.",
				},
				{
					name: "hour", api: "hour", path: "restrict_schedule.hour",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"00\" - \"23\" Defaults to `\"*\"`.",
				},
				{
					name: "dom", api: "dom", path: "restrict_schedule.dom",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"1\" - \"31\" Defaults to `\"*\"`.",
				},
				{
					name: "month", api: "month", path: "restrict_schedule.month",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"1\" (January) - \"12\" (December) Defaults to `\"*\"`.",
				},
				{
					name: "dow", api: "dow", path: "restrict_schedule.dow",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "*",
					description: "\"1\" (Monday) - \"7\" (Sunday) Defaults to `\"*\"`.",
				},
				{
					name: "begin", api: "begin", path: "restrict_schedule.begin",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "00:00",
					description: "Start time for the time window in HH:MM format. Defaults to `\"00:00\"`.",
				},
				{
					name: "end", api: "end", path: "restrict_schedule.end",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "23:59",
					description: "End time for the time window in HH:MM format. Defaults to `\"23:59\"`.",
				},
			},
		},
		{
			name: "only_matching_schedule", api: "only_matching_schedule", path: "only_matching_schedule",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Will only replicate snapshots that match `schedule` or `restrict_schedule`. Defaults to `false`.",
		},
		{
			name: "allow_from_scratch", api: "allow_from_scratch", path: "allow_from_scratch",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Will destroy all snapshots on target side and replicate everything from scratch if none of the snapshots on     target side matches source snapshots. Defaults to `false`.",
		},
		{
			name: "readonly", api: "readonly", path: "readonly",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "SET",
			description: "Controls destination datasets readonly property.\n\n* `SET`: Set all destination datasets to readonly=on after finishing the replication.\n* `REQUIRE`: Require all existing destination datasets to have readonly=on property.\n* `IGNORE`: Avoid this kind of behavior. Defaults to `\"SET\"`.",
		},
		{
			name: "hold_pending_snapshots", api: "hold_pending_snapshots", path: "hold_pending_snapshots",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Prevent source snapshots from being deleted by retention of replication fails for some reason. Defaults to `false`.",
		},
		{
			name: "retention_policy", api: "retention_policy", path: "retention_policy",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "How to delete old snapshots on target side:\n\n* `SOURCE`: Delete snapshots that are absent on source side.\n* `CUSTOM`: Delete snapshots that are older than `lifetime_value` and `lifetime_unit`.\n* `NONE`: Do not delete any snapshots.",
		},
		{
			name: "lifetime_value", api: "lifetime_value", path: "lifetime_value",
			kind: kindInt, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Number of time units to retain snapshots for custom retention policy. Only applies when `retention_policy` is     CUSTOM.",
		},
		{
			name: "lifetime_unit", api: "lifetime_unit", path: "lifetime_unit",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Time unit for snapshot retention for custom retention policy. Only applies when `retention_policy` is CUSTOM.",
		},
		{
			name: "lifetimes", api: "lifetimes", path: "lifetimes",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "Array of different retention schedules with their own cron schedules and lifetime settings. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "lifetimes",
				kind: kindObject, role: roleRequired,
				readable: true,
				children: []*node{
					{
						name: "schedule", api: "schedule", path: "lifetimes.schedule",
						kind: kindObject, role: roleRequired,
						readable:    true,
						description: "Cron schedule for when snapshot retention policies are applied.",
						children: []*node{
							{
								name: "minute", api: "minute", path: "lifetimes.schedule.minute",
								kind: kindString, role: roleOptionalComputed,
								readable:    true,
								description: "\"00\" - \"59\"",
							},
							{
								name: "hour", api: "hour", path: "lifetimes.schedule.hour",
								kind: kindString, role: roleOptionalComputed,
								readable:    true,
								description: "\"00\" - \"23\"",
							},
							{
								name: "dom", api: "dom", path: "lifetimes.schedule.dom",
								kind: kindString, role: roleOptionalComputed,
								readable:    true,
								description: "\"1\" - \"31\"",
							},
							{
								name: "month", api: "month", path: "lifetimes.schedule.month",
								kind: kindString, role: roleOptionalComputed,
								readable:    true,
								description: "\"1\" (January) - \"12\" (December)",
							},
							{
								name: "dow", api: "dow", path: "lifetimes.schedule.dow",
								kind: kindString, role: roleOptionalComputed,
								readable:    true,
								description: "\"1\" (Monday) - \"7\" (Sunday)",
							},
						},
					},
					{
						name: "lifetime_value", api: "lifetime_value", path: "lifetimes.lifetime_value",
						kind: kindInt, role: roleRequired,
						readable:    true,
						description: "Number of time units to retain snapshots.",
					},
					{
						name: "lifetime_unit", api: "lifetime_unit", path: "lifetimes.lifetime_unit",
						kind: kindString, role: roleRequired,
						readable:    true,
						description: "Time unit for snapshot retention.",
					},
				},
			},
		},
		{
			name: "compression", api: "compression", path: "compression",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Compresses SSH stream. Available only for SSH transport.",
		},
		{
			name: "speed_limit", api: "speed_limit", path: "speed_limit",
			kind: kindInt, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Limits speed of SSH stream. Available only for SSH transport.",
		},
		{
			name: "large_block", api: "large_block", path: "large_block",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Enable large block support for ZFS send streams. Defaults to `true`.",
		},
		{
			name: "embed", api: "embed", path: "embed",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Enable embedded block support for ZFS send streams. Defaults to `false`.",
		},
		{
			name: "compressed", api: "compressed", path: "compressed",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Enable compressed ZFS send streams. Defaults to `true`.",
		},
		{
			name: "retries", api: "retries", path: "retries",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: 5,
			description: "Number of retries before considering replication failed. Defaults to `5`.",
		},
		{
			name: "logging_level", api: "logging_level", path: "logging_level",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Log level for replication task execution. Controls verbosity of replication logs.",
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether this replication task is enabled. Defaults to `true`.",
		},
		{
			name: "has_encrypted_dataset_keys", api: "has_encrypted_dataset_keys", path: "has_encrypted_dataset_keys",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether this replication task has encrypted dataset keys available.",
		},
		{
			name: "encryption_key_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `encryption_key` again. Terraform never stores `encryption_key`.",
		},
	},
}

// NewReplication returns the replication resource.
func NewReplication() resource.Resource {
	return &replicationResource{crudResource{model: modelReplication}}
}

type replicationResource struct{ crudResource }

// NewReplicationDataSource returns the replication data source.
func NewReplicationDataSource() datasource.DataSource {
	return &dataSource{model: modelReplication, attrs: dataAttrs(modelReplication)}
}

// NewReplicationList returns the replication list resource, for terraform query.
func NewReplicationList() list.ListResource {
	return &listResource{crudResource{model: modelReplication}}
}

func (r *replicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Replication Task that will push or pull ZFS snapshots to or from remote host.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name for replication task.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"direction": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Whether task will `PUSH` or `PULL` snapshots.",
				Validators:          []validator.String{stringvalidator.OneOf("PUSH", "PULL")},
			},
			"transport": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Method of snapshots transfer.\n\n* `SSH` transfers snapshots via SSH connection. This method is supported everywhere but does not achieve       great performance.\n* `SSH+NETCAT` uses unencrypted connection for data transfer. This can only be used in trusted networks       and requires a port (specified by range from `netcat_active_side_port_min` to `netcat_active_side_port_max`)       to be open on `netcat_active_side`.\n* `LOCAL` replicates to or from localhost.",
				Validators:          []validator.String{stringvalidator.OneOf("SSH", "SSH+NETCAT", "LOCAL")},
			},
			"ssh_credentials": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Keychain Credential ID of type `SSH_CREDENTIALS`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"netcat_active_side": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Which side actively establishes the netcat connection for `SSH+NETCAT` transport.\n\n* `LOCAL`: Local system initiates the connection\n* `REMOTE`: Remote system initiates the connection\n* `null`: Not applicable for other transport types",
				Validators:          []validator.String{stringvalidator.OneOf("LOCAL", "REMOTE")},
			},
			"netcat_active_side_listen_address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "IP address for the active side to listen on for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"netcat_active_side_port_min": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum port number in the range for netcat connections. `null` if not applicable.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"netcat_active_side_port_max": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum port number in the range for netcat connections. `null` if not applicable.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"netcat_passive_side_connect_address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "IP address for the passive side to connect to for `SSH+NETCAT` transport. `null` if not applicable.",
			},
			"sudo": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "`SSH` and `SSH+NETCAT` transports should use sudo (which is expected to be passwordless) to run `zfs`     command on the remote machine. Defaults to `false`.",
			},
			"source_datasets": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "List of datasets to replicate snapshots from.",
				ElementType:         types.StringType,
			},
			"target_dataset": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Dataset to put snapshots into.",
			},
			"recursive": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether to recursively replicate child datasets.",
			},
			"exclude": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of dataset patterns to exclude from replication. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"properties": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Send dataset properties along with snapshots. Defaults to `true`.",
			},
			"properties_exclude": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of dataset property names to exclude from replication. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"properties_override": schema.MapAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Object mapping dataset property names to override values during replication. Defaults to `{}`.",
				ElementType:         types.StringType,
			},
			"replicate": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to use full ZFS replication. Defaults to `false`.",
			},
			"encryption": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable encryption for the replicated datasets. Defaults to `false`.",
			},
			"encryption_inherit": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether replicated datasets should inherit encryption from parent. `null` if encryption is disabled.",
			},
			"encryption_key": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Encryption key for replicated datasets. `null` if not specified.",
			},
			"encryption_key_format": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Format of the encryption key.\n\n* `HEX`: Hexadecimal-encoded key\n* `PASSPHRASE`: Text passphrase\n* `null`: Not applicable when encryption is disabled",
				Validators:          []validator.String{stringvalidator.OneOf("HEX", "PASSPHRASE")},
			},
			"encryption_key_location": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filesystem path where encryption key is stored. `null` if not using key file.",
			},
			"periodic_snapshot_tasks": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "List of periodic snapshot task IDs that are sources of snapshots for this replication task. Only push     replication tasks can be bound to periodic snapshot tasks. Defaults to `[]`.",
				ElementType:         types.NumberType,
			},
			"naming_schema": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "List of naming schemas for pull replication. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"also_include_naming_schema": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "List of naming schemas for push replication. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"name_regex": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Replicate all snapshots which names match specified regular expression.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"auto": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Allow replication to run automatically on schedule or after bound periodic snapshot task.",
			},
			"schedule": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Schedule to run replication task. Only `auto` replication tasks without bound periodic snapshot tasks can have     a schedule.",
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
					"begin": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Start time for the time window in HH:MM format. Defaults to `\"00:00\"`.",
						Default:             stringdefault.StaticString("00:00"),
					},
					"end": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "End time for the time window in HH:MM format. Defaults to `\"23:59\"`.",
						Default:             stringdefault.StaticString("23:59"),
					},
				},
			},
			"restrict_schedule": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Restricts when replication task with bound periodic snapshot tasks runs. For example, you can have periodic     snapshot tasks that run every 15 minutes, but only run replication task every hour.",
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
					"begin": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Start time for the time window in HH:MM format. Defaults to `\"00:00\"`.",
						Default:             stringdefault.StaticString("00:00"),
					},
					"end": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "End time for the time window in HH:MM format. Defaults to `\"23:59\"`.",
						Default:             stringdefault.StaticString("23:59"),
					},
				},
			},
			"only_matching_schedule": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Will only replicate snapshots that match `schedule` or `restrict_schedule`. Defaults to `false`.",
			},
			"allow_from_scratch": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Will destroy all snapshots on target side and replicate everything from scratch if none of the snapshots on     target side matches source snapshots. Defaults to `false`.",
			},
			"readonly": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Controls destination datasets readonly property.\n\n* `SET`: Set all destination datasets to readonly=on after finishing the replication.\n* `REQUIRE`: Require all existing destination datasets to have readonly=on property.\n* `IGNORE`: Avoid this kind of behavior. Defaults to `\"SET\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("SET", "REQUIRE", "IGNORE")},
			},
			"hold_pending_snapshots": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Prevent source snapshots from being deleted by retention of replication fails for some reason. Defaults to `false`.",
			},
			"retention_policy": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "How to delete old snapshots on target side:\n\n* `SOURCE`: Delete snapshots that are absent on source side.\n* `CUSTOM`: Delete snapshots that are older than `lifetime_value` and `lifetime_unit`.\n* `NONE`: Do not delete any snapshots.",
				Validators:          []validator.String{stringvalidator.OneOf("SOURCE", "CUSTOM", "NONE")},
			},
			"lifetime_value": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Number of time units to retain snapshots for custom retention policy. Only applies when `retention_policy` is     CUSTOM.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}},
			},
			"lifetime_unit": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Time unit for snapshot retention for custom retention policy. Only applies when `retention_policy` is CUSTOM.",
				Validators:          []validator.String{stringvalidator.OneOf("HOUR", "DAY", "WEEK", "MONTH", "YEAR")},
			},
			"lifetimes": schema.ListNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of different retention schedules with their own cron schedules and lifetime settings. Defaults to `[]`.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"schedule": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "Cron schedule for when snapshot retention policies are applied.",
							Attributes: map[string]schema.Attribute{
								"minute": schema.StringAttribute{
									Optional: true, Computed: true,
									MarkdownDescription: "\"00\" - \"59\"",
								},
								"hour": schema.StringAttribute{
									Optional: true, Computed: true,
									MarkdownDescription: "\"00\" - \"23\"",
								},
								"dom": schema.StringAttribute{
									Optional: true, Computed: true,
									MarkdownDescription: "\"1\" - \"31\"",
								},
								"month": schema.StringAttribute{
									Optional: true, Computed: true,
									MarkdownDescription: "\"1\" (January) - \"12\" (December)",
								},
								"dow": schema.StringAttribute{
									Optional: true, Computed: true,
									MarkdownDescription: "\"1\" (Monday) - \"7\" (Sunday)",
								},
							},
						},
						"lifetime_value": schema.NumberAttribute{
							Required:            true,
							MarkdownDescription: "Number of time units to retain snapshots.",
							Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}},
						},
						"lifetime_unit": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Time unit for snapshot retention.",
							Validators:          []validator.String{stringvalidator.OneOf("HOUR", "DAY", "WEEK", "MONTH", "YEAR")},
						},
					},
				},
			},
			"compression": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Compresses SSH stream. Available only for SSH transport.",
				Validators:          []validator.String{stringvalidator.OneOf("LZ4", "PIGZ", "PLZIP")},
			},
			"speed_limit": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Limits speed of SSH stream. Available only for SSH transport.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}},
			},
			"large_block": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable large block support for ZFS send streams. Defaults to `true`.",
			},
			"embed": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable embedded block support for ZFS send streams. Defaults to `false`.",
			},
			"compressed": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable compressed ZFS send streams. Defaults to `true`.",
			},
			"retries": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Number of retries before considering replication failed. Defaults to `5`.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}},
			},
			"logging_level": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Log level for replication task execution. Controls verbosity of replication logs.",
				Validators:          []validator.String{stringvalidator.OneOf("DEBUG", "INFO", "WARNING", "ERROR")},
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether this replication task is enabled. Defaults to `true`.",
			},
			"has_encrypted_dataset_keys": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether this replication task has encrypted dataset keys available.",
			},
			"encryption_key_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `encryption_key` again. Terraform never stores `encryption_key`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
