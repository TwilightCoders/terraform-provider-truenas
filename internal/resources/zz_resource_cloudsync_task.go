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

var modelCloudsyncTask = &model{
	typeName:     "cloudsync_task",
	namespace:    "cloudsync",
	primaryKey:   "id",
	idKind:       kindInt,
	deleteMethod: "cloudsync.delete",
	createMethod: "cloudsync.create",
	updateMethod: "cloudsync.update",
	getMethod:    "cloudsync.get_instance",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "description", api: "description", path: "description",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "The name of the task to display in the UI. Defaults to `\"\"`.",
		},
		{
			name: "path", api: "path", path: "path",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "The local path to back up beginning with `/mnt` or `/dev/zvol`.",
		},
		{
			name: "credentials", api: "credentials", path: "credentials",
			kind: kindInt, role: roleRequired,
			readable: true, updatable: true,
			description: "ID of the cloud credential.",
			ref:         "id",
		},
		{
			name: "attributes", api: "attributes", path: "attributes",
			kind: kindObject, role: roleRequired,
			readable: true, updatable: true,
			description: "Additional information for each backup, e.g. bucket name.",
			children: []*node{
				{
					name: "bucket", api: "bucket", path: "attributes.bucket",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Name of the cloud storage bucket or container.",
				},
				{
					name: "folder", api: "folder", path: "attributes.folder",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Path within the cloud storage bucket to use as the root directory for operations.",
				},
				{
					name: "fast_list", api: "fast_list", path: "attributes.fast_list",
					kind: kindBool, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Valid only for some providers. Use fewer transactions in exchange for more RAM. This may also speed up or slow     down your transfer. See https://rclone.org/docs/#fast-list for more details.",
				},
				{
					name: "bucket_policy_only", api: "bucket_policy_only", path: "attributes.bucket_policy_only",
					kind: kindBool, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Valid only for GOOGLE_CLOUD_STORAGE provider. Access checks should use bucket-level IAM policies. If you want     to upload objects to a bucket with Bucket Policy Only set then you will need to set this.",
				},
				{
					name: "chunk_size", api: "chunk_size", path: "attributes.chunk_size",
					kind: kindInt, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Valid only for DROPBOX provider. Upload chunk size in MiB. Must fit in memory. Note that these chunks are     buffered in memory and there might be a maximum of `--transfers` chunks in progress at once. Dropbox Business     accounts can have monthly data transfer limits per team per month. By using larger chunk sizes you will decrease     the number of data transfer calls used and you'll be able to transfer more data to your Dropbox Business account.",
				},
				{
					name: "acknowledge_abuse", api: "acknowledge_abuse", path: "attributes.acknowledge_abuse",
					kind: kindBool, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Valid only for GOOGLE_DRIVER provider. Allow files which return cannotDownloadAbusiveFile to be downloaded. If     downloading a file returns the error \"This file has been identified as malware or spam and cannot be downloaded\"     with the error code \"cannotDownloadAbusiveFile\" then enable this flag to indicate you acknowledge the risks of     downloading the file and TrueNAS will download it anyway.",
				},
				{
					name: "region", api: "region", path: "attributes.region",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Valid only for S3 provider. S3 Region.",
				},
				{
					name: "encryption", api: "encryption", path: "attributes.encryption",
					kind: kindString, role: roleOptionalComputed,
					nullable: true, readable: true, updatable: true,
					description: "Valid only for S3 provider. Server-Side Encryption.",
				},
				{
					name: "storage_class", api: "storage_class", path: "attributes.storage_class",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Valid only for S3 provider. The storage class to use.",
				},
			},
		},
		{
			name: "schedule", api: "schedule", path: "schedule",
			kind: kindObject, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: map[string]any{"dom": "*", "dow": "*", "hour": "*", "minute": "00", "month": "*"},
			description: "Cron schedule dictating when the task should run. Defaults to `{\"dom\":\"*\",\"dow\":\"*\",\"hour\":\"*\",\"minute\":\"00\",\"month\":\"*\"}`.",
			children: []*node{
				{
					name: "minute", api: "minute", path: "schedule.minute",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: "00",
					description: "Defaults to `\"00\"`.",
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
			name: "pre_script", api: "pre_script", path: "pre_script",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "A Bash script to run immediately before every backup. Defaults to `\"\"`.",
		},
		{
			name: "post_script", api: "post_script", path: "post_script",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "A Bash script to run immediately after every backup if it succeeds. Defaults to `\"\"`.",
		},
		{
			name: "snapshot", api: "snapshot", path: "snapshot",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to create a temporary snapshot of the dataset before every backup. Defaults to `false`.",
		},
		{
			name: "include", api: "include", path: "include",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Paths to pass to `restic backup --include`.",
			elem: &node{
				name: "", api: "", path: "include",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "exclude", api: "exclude", path: "exclude",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Paths to pass to `restic backup --exclude`.",
			elem: &node{
				name: "", api: "", path: "exclude",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "args", api: "args", path: "args",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "(Slated for removal). Defaults to `\"\"`.",
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Can enable/disable the task. Defaults to `true`.",
		},
		{
			name: "bwlimit", api: "bwlimit", path: "bwlimit",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Schedule of bandwidth limits.",
			elem: &node{
				name: "", api: "", path: "bwlimit",
				kind: kindObject, role: roleRequired,
				readable: true,
				children: []*node{
					{
						name: "time", api: "time", path: "bwlimit.time",
						kind: kindString, role: roleRequired,
						readable:    true,
						description: "Time at which the bandwidth limit takes effect in 24-hour format.",
					},
					{
						name: "bandwidth", api: "bandwidth", path: "bwlimit.bandwidth",
						kind: kindInt, role: roleRequired,
						nullable: true, readable: true,
						description: "Bandwidth limit in bytes per second (upload and download).",
					},
				},
			},
		},
		{
			name: "transfers", api: "transfers", path: "transfers",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Maximum number of parallel file transfers. `null` for default.",
		},
		{
			name: "direction", api: "direction", path: "direction",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Direction of the cloud sync operation.\n\n* `PUSH`: Upload local files to cloud storage\n* `PULL`: Download files from cloud storage to local storage",
		},
		{
			name: "transfer_mode", api: "transfer_mode", path: "transfer_mode",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "How files are transferred between local and cloud storage.\n\n* `SYNC`: Synchronize directories (add new, update changed, remove deleted)\n* `COPY`: Copy files without removing any existing files\n* `MOVE`: Move files (copy then delete from source)",
		},
		{
			name: "encryption", api: "encryption", path: "encryption",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to encrypt files before uploading to cloud storage. Defaults to `false`.",
		},
		{
			name: "filename_encryption", api: "filename_encryption", path: "filename_encryption",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to encrypt filenames in addition to file contents. Defaults to `false`.",
		},
		{
			name: "encryption_password", api: "encryption_password", path: "encryption_password",
			kind: kindString, role: roleOptional,
			sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Password for client-side encryption. Empty string if encryption is disabled.",
		},
		{
			name: "encryption_salt", api: "encryption_salt", path: "encryption_salt",
			kind: kindString, role: roleOptional,
			sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Salt value for encryption key derivation. Empty string if encryption is disabled.",
		},
		{
			name: "create_empty_src_dirs", api: "create_empty_src_dirs", path: "create_empty_src_dirs",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to create empty directories in the destination that exist in the source. Defaults to `false`.",
		},
		{
			name: "follow_symlinks", api: "follow_symlinks", path: "follow_symlinks",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to follow symbolic links and sync the files they point to. Defaults to `false`.",
		},
		{
			name: "locked", api: "locked", path: "locked",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "A locked task cannot run.",
		},
		{
			name: "encryption_password_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `encryption_password` again. Terraform never stores `encryption_password`.",
		},
		{
			name: "encryption_salt_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `encryption_salt` again. Terraform never stores `encryption_salt`.",
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

// NewCloudsyncTask returns the cloudsync_task resource.
func NewCloudsyncTask() resource.Resource {
	return &cloudsyncTaskResource{crudResource{model: modelCloudsyncTask}}
}

type cloudsyncTaskResource struct{ crudResource }

// NewCloudsyncTaskDataSource returns the cloudsync_task data source.
func NewCloudsyncTaskDataSource() datasource.DataSource {
	return &dataSource{model: modelCloudsyncTask, attrs: dataAttrs(modelCloudsyncTask)}
}

// NewCloudsyncTaskList returns the cloudsync_task list resource, for terraform query.
func NewCloudsyncTaskList() list.ListResource {
	return &listResource{crudResource{model: modelCloudsyncTask}}
}

func (r *cloudsyncTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a new cloud_sync entry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The name of the task to display in the UI. Defaults to `\"\"`.",
			},
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The local path to back up beginning with `/mnt` or `/dev/zvol`.",
			},
			"credentials": schema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "ID of the cloud credential.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"attributes": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Additional information for each backup, e.g. bucket name.",
				Attributes: map[string]schema.Attribute{
					"bucket": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Name of the cloud storage bucket or container.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"folder": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Path within the cloud storage bucket to use as the root directory for operations.",
					},
					"fast_list": schema.BoolAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Valid only for some providers. Use fewer transactions in exchange for more RAM. This may also speed up or slow     down your transfer. See https://rclone.org/docs/#fast-list for more details.",
					},
					"bucket_policy_only": schema.BoolAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Valid only for GOOGLE_CLOUD_STORAGE provider. Access checks should use bucket-level IAM policies. If you want     to upload objects to a bucket with Bucket Policy Only set then you will need to set this.",
					},
					"chunk_size": schema.NumberAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Valid only for DROPBOX provider. Upload chunk size in MiB. Must fit in memory. Note that these chunks are     buffered in memory and there might be a maximum of `--transfers` chunks in progress at once. Dropbox Business     accounts can have monthly data transfer limits per team per month. By using larger chunk sizes you will decrease     the number of data transfer calls used and you'll be able to transfer more data to your Dropbox Business account.",
						Validators:          []validator.Number{numberIsInteger{}},
					},
					"acknowledge_abuse": schema.BoolAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Valid only for GOOGLE_DRIVER provider. Allow files which return cannotDownloadAbusiveFile to be downloaded. If     downloading a file returns the error \"This file has been identified as malware or spam and cannot be downloaded\"     with the error code \"cannotDownloadAbusiveFile\" then enable this flag to indicate you acknowledge the risks of     downloading the file and TrueNAS will download it anyway.",
					},
					"region": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Valid only for S3 provider. S3 Region.",
					},
					"encryption": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Valid only for S3 provider. Server-Side Encryption.",
						Validators:          []validator.String{stringvalidator.OneOf("AES256")},
					},
					"storage_class": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Valid only for S3 provider. The storage class to use.",
						Validators:          []validator.String{stringvalidator.OneOf("", "STANDARD", "REDUCED_REDUNDANCY", "STANDARD_IA", "ONEZONE_IA", "INTELLIGENT_TIERING", "GLACIER", "GLACIER_IR", "DEEP_ARCHIVE")},
					},
				},
			},
			"schedule": schema.SingleNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Cron schedule dictating when the task should run. Defaults to `{\"dom\":\"*\",\"dow\":\"*\",\"hour\":\"*\",\"minute\":\"00\",\"month\":\"*\"}`.",
				Attributes: map[string]schema.Attribute{
					"minute": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Defaults to `\"00\"`.",
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
			},
			"pre_script": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A Bash script to run immediately before every backup. Defaults to `\"\"`.",
			},
			"post_script": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A Bash script to run immediately after every backup if it succeeds. Defaults to `\"\"`.",
			},
			"snapshot": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to create a temporary snapshot of the dataset before every backup. Defaults to `false`.",
			},
			"include": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Paths to pass to `restic backup --include`.",
				ElementType:         types.StringType,
			},
			"exclude": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Paths to pass to `restic backup --exclude`.",
				ElementType:         types.StringType,
			},
			"args": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "(Slated for removal). Defaults to `\"\"`.",
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Can enable/disable the task. Defaults to `true`.",
			},
			"bwlimit": schema.ListNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Schedule of bandwidth limits.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"time": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Time at which the bandwidth limit takes effect in 24-hour format.",
						},
						"bandwidth": schema.NumberAttribute{
							Required:            true,
							MarkdownDescription: "Bandwidth limit in bytes per second (upload and download).",
							Validators:          []validator.Number{numberIsInteger{}},
						},
					},
				},
			},
			"transfers": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum number of parallel file transfers. `null` for default.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"direction": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Direction of the cloud sync operation.\n\n* `PUSH`: Upload local files to cloud storage\n* `PULL`: Download files from cloud storage to local storage",
				Validators:          []validator.String{stringvalidator.OneOf("PUSH", "PULL")},
			},
			"transfer_mode": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "How files are transferred between local and cloud storage.\n\n* `SYNC`: Synchronize directories (add new, update changed, remove deleted)\n* `COPY`: Copy files without removing any existing files\n* `MOVE`: Move files (copy then delete from source)",
				Validators:          []validator.String{stringvalidator.OneOf("SYNC", "COPY", "MOVE")},
			},
			"encryption": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to encrypt files before uploading to cloud storage. Defaults to `false`.",
			},
			"filename_encryption": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to encrypt filenames in addition to file contents. Defaults to `false`.",
			},
			"encryption_password": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Password for client-side encryption. Empty string if encryption is disabled.",
			},
			"encryption_salt": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Salt value for encryption key derivation. Empty string if encryption is disabled.",
			},
			"create_empty_src_dirs": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to create empty directories in the destination that exist in the source. Defaults to `false`.",
			},
			"follow_symlinks": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to follow symbolic links and sync the files they point to. Defaults to `false`.",
			},
			"locked": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "A locked task cannot run.",
			},
			"encryption_password_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `encryption_password` again. Terraform never stores `encryption_password`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"encryption_salt_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `encryption_salt` again. Terraform never stores `encryption_salt`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"unset": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
				ElementType:         types.StringType,
			},
		},
	}
}
