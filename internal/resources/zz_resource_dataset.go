// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var modelDataset = &model{
	typeName:      "dataset",
	namespace:     "pool.dataset",
	primaryKey:    "id",
	idKind:        kindString,
	discriminator: "type",
	variant:       "FILESYSTEM",
	createMethod:  "pool.dataset.create",
	updateMethod:  "pool.dataset.update",
	getMethod:     "pool.dataset.get_instance",
	deleteMethod:  "pool.dataset.delete",
	listFilters:   [][]interface{}(nil),
	listOptions:   map[string]interface{}{"extra": map[string]interface{}{"flat": true}},
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindString, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			replace: true, readable: true,
			description: "The name of the dataset to create. Changing this forces a new resource.",
		},
		{
			name: "comments", api: "comments", path: "comments",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Comments or description for the dataset. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "sync", api: "sync", path: "sync",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Synchronous write behavior for the dataset. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "snapdev", api: "snapdev", path: "snapdev",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			description: "Controls visibility of volume snapshots under /dev/zvol/.",
			dialect:     apiDialect,
		},
		{
			name: "compression", api: "compression", path: "compression",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Compression algorithm to use for the dataset. Higher numbered variants provide better compression     but use more CPU. 'INHERIT' uses the parent dataset's setting. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "exec", api: "exec", path: "exec",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Whether files in this dataset can be executed. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "managedby", api: "managedby", path: "managedby",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Identifies which service or system manages this dataset. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "quota_warning", api: "quota_warning", path: "quota_warning",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true, intOrString: true,
			hasDefault: true, def: "INHERIT",
			description: "Percentage of dataset quota at which to issue a warning. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "quota_critical", api: "quota_critical", path: "quota_critical",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true, intOrString: true,
			hasDefault: true, def: "INHERIT",
			description: "Percentage of dataset quota at which to issue a critical alert. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "refquota_warning", api: "refquota_warning", path: "refquota_warning",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true, intOrString: true,
			hasDefault: true, def: "INHERIT",
			description: "Percentage of reference quota at which to issue a warning. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "refquota_critical", api: "refquota_critical", path: "refquota_critical",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true, intOrString: true,
			hasDefault: true, def: "INHERIT",
			description: "Percentage of reference quota at which to issue a critical alert. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "reservation", api: "reservation", path: "reservation",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true, property: true,
			description: "Minimum disk space guaranteed to this dataset and its children in bytes.",
			dialect:     apiDialect,
		},
		{
			name: "refreservation", api: "refreservation", path: "refreservation",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true, property: true,
			description: "Minimum disk space guaranteed to this dataset itself in bytes.",
			dialect:     apiDialect,
		},
		{
			name: "special_small_block_size", api: "special_small_block_size", path: "special_small_block_size",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true, intOrString: true,
			description: "Size threshold below which blocks are stored on special vdevs.",
			dialect:     apiDialect,
		},
		{
			name: "copies", api: "copies", path: "copies",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true, intOrString: true,
			hasDefault: true, def: "INHERIT",
			description: "Number of copies of data blocks to maintain for redundancy. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "snapdir", api: "snapdir", path: "snapdir",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Controls visibility of the `.zfs/snapshot` directory. 'DISABLED' hides snapshots, 'VISIBLE' shows them,     'HIDDEN' makes them accessible but not listed. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "deduplication", api: "deduplication", path: "deduplication",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Deduplication setting. 'ON' enables dedup, 'VERIFY' enables with checksum verification, 'OFF' disables. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "checksum", api: "checksum", path: "checksum",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Checksum algorithm to verify data integrity. Higher security algorithms like SHA256 provide better     protection but use more CPU. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "readonly", api: "readonly", path: "readonly",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			hasDefault: true, def: "INHERIT",
			description: "Whether the dataset is read-only. Defaults to `\"INHERIT\"`.",
			dialect:     apiDialect,
		},
		{
			name: "share_type", api: "share_type", path: "share_type",
			kind: kindString, role: roleOptional,
			createOnly:  true,
			description: "Optimization type for the dataset based on its intended use. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "encryption_options", api: "encryption_options", path: "encryption_options",
			kind: kindObject, role: roleOptional,
			createOnly:  true,
			description: "Configuration for encryption of dataset for `name` pool. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
			children: []*node{
				{
					name: "generate_key", api: "generate_key", path: "encryption_options.generate_key",
					kind: kindBool, role: roleOptional,
					description: "Automatically generate the key to be used for dataset encryption. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				},
				{
					name: "pbkdf2iters", api: "pbkdf2iters", path: "encryption_options.pbkdf2iters",
					kind: kindInt, role: roleOptional,
					description: "Number of PBKDF2 iterations for key derivation from passphrase. Higher iterations improve security     against brute force attacks but increase unlock time. Default 350,000 balances security and performance. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				},
				{
					name: "algorithm", api: "algorithm", path: "encryption_options.algorithm",
					kind: kindString, role: roleOptional,
					description: "Encryption algorithm to use for dataset encryption. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				},
				{
					name: "passphrase", api: "passphrase", path: "encryption_options.passphrase",
					kind: kindString, role: roleOptional,
					nullable: true, sensitive: true, writeOnly: true,
					description: "Must be specified if encryption for root dataset is desired with a passphrase as a key.",
				},
				{
					name: "key", api: "key", path: "encryption_options.key",
					kind: kindString, role: roleOptional,
					nullable: true, sensitive: true, writeOnly: true,
					description: "A hex-encoded key specified as an alternative to using `passphrase`.",
				},
			},
		},
		{
			name: "encryption", api: "encryption", path: "encryption",
			kind: kindBool, role: roleOptional,
			createOnly:  true,
			description: "Create a ZFS encrypted root dataset for `name` pool.\nThere is 1 case where ZFS encryption is not allowed for a dataset:\n1) If the parent dataset is encrypted with a passphrase and `name` is being created with a key for encrypting the        dataset. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "inherit_encryption", api: "inherit_encryption", path: "inherit_encryption",
			kind: kindBool, role: roleOptional,
			createOnly:  true,
			description: "Whether to inherit encryption settings from the parent dataset. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "create_ancestors", api: "create_ancestors", path: "create_ancestors",
			kind: kindBool, role: roleOptional,
			updatable:   true,
			description: "Whether to create any missing parent datasets. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
		},
		{
			name: "aclmode", api: "aclmode", path: "aclmode",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			description: "How Access Control Lists are handled when chmod is used.",
			dialect:     apiDialect,
		},
		{
			name: "acltype", api: "acltype", path: "acltype",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			description: "The type of Access Control List system to use.",
			dialect:     apiDialect,
		},
		{
			name: "atime", api: "atime", path: "atime",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true, inherit: true,
			description: "Whether file access times are updated when files are accessed.",
			dialect:     apiDialect,
		},
		{
			name: "casesensitivity", api: "casesensitivity", path: "casesensitivity",
			kind: kindString, role: roleOptionalComputed,
			replace: true, stable: true, readable: true, property: true, inherit: true,
			description: "File name case sensitivity setting. Changing this forces a new resource.",
			dialect:     apiDialect,
		},
		{
			name: "quota", api: "quota", path: "quota",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true, property: true,
			description: "Maximum disk space this dataset and its children can consume in bytes.",
			dialect:     apiDialect,
		},
		{
			name: "refquota", api: "refquota", path: "refquota",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true, property: true,
			description: "Maximum disk space this dataset itself can consume in bytes.",
			dialect:     apiDialect,
		},
		{
			name: "recordsize", api: "recordsize", path: "recordsize",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true, property: true,
			description: "The suggested block size for files in this filesystem dataset.",
			dialect:     apiDialect,
		},
		{
			name: "pool", api: "pool", path: "pool",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "The name of the ZFS pool containing this dataset.",
		},
		{
			name: "encrypted", api: "encrypted", path: "encrypted",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether the dataset is encrypted.",
		},
		{
			name: "encryption_root", api: "encryption_root", path: "encryption_root",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "The root dataset where encryption is enabled. `null` if the dataset is not encrypted.",
		},
		{
			name: "key_loaded", api: "key_loaded", path: "key_loaded",
			kind: kindBool, role: roleComputed,
			nullable: true, readable: true,
			description: "Whether the encryption key is currently loaded for encrypted datasets. `null` for unencrypted datasets.",
		},
		{
			name: "locked", api: "locked", path: "locked",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether an encrypted dataset is currently locked (key not loaded).",
		},
		{
			name: "xattr", api: "xattr", path: "xattr",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Extended attributes storage method (on/off).",
			dialect:     apiDialect,
		},
		{
			name: "compressratio", api: "compressratio", path: "compressratio",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "The achieved compression ratio as a decimal (e.g., '2.50x').",
			dialect:     apiDialect,
		},
		{
			name: "origin", api: "origin", path: "origin",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "The snapshot from which this clone was created. Empty for non-clone datasets.",
			dialect:     apiDialect,
		},
		{
			name: "key_format", api: "key_format", path: "key_format",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Format of the encryption key (hex/raw/passphrase). Only relevant for encrypted datasets.",
			dialect:     apiDialect,
		},
		{
			name: "encryption_algorithm", api: "encryption_algorithm", path: "encryption_algorithm",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Encryption algorithm used (e.g., AES-256-GCM). Only relevant for encrypted datasets.",
			dialect:     apiDialect,
		},
		{
			name: "used", api: "used", path: "used",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Total amount of disk space consumed by this dataset and all its children.",
			dialect:     apiDialect,
		},
		{
			name: "usedbychildren", api: "usedbychildren", path: "usedbychildren",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Amount of disk space consumed by child datasets.",
			dialect:     apiDialect,
		},
		{
			name: "usedbydataset", api: "usedbydataset", path: "usedbydataset",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Amount of disk space consumed by this dataset itself, excluding children and snapshots.",
			dialect:     apiDialect,
		},
		{
			name: "usedbyrefreservation", api: "usedbyrefreservation", path: "usedbyrefreservation",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Amount of disk space consumed by the refreservation of this dataset.",
			dialect:     apiDialect,
		},
		{
			name: "usedbysnapshots", api: "usedbysnapshots", path: "usedbysnapshots",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Amount of disk space consumed by snapshots of this dataset.",
			dialect:     apiDialect,
		},
		{
			name: "available", api: "available", path: "available",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Amount of disk space available to this dataset and its children.",
			dialect:     apiDialect,
		},
		{
			name: "pbkdf2iters", api: "pbkdf2iters", path: "pbkdf2iters",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Number of PBKDF2 iterations used for passphrase-based encryption keys.",
			dialect:     apiDialect,
		},
		{
			name: "creation", api: "creation", path: "creation",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true, property: true, rawProperty: true,
			description: "Timestamp when this dataset was created.",
			dialect:     apiDialect,
		},
		{
			name: "mountpoint", api: "mountpoint", path: "mountpoint",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Filesystem path where this dataset is mounted. Null for unmounted datasets or volumes.",
		},
		{
			name: "encryption_options_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `encryption_options` again. Terraform never stores `encryption_options`.",
		},
	},
}

// NewDataset returns the dataset resource.
func NewDataset() resource.Resource {
	return &datasetResource{crudResource{model: modelDataset}}
}

type datasetResource struct{ crudResource }

// NewDatasetDataSource returns the dataset data source.
func NewDatasetDataSource() datasource.DataSource {
	return &dataSource{model: modelDataset, attrs: dataAttrs(modelDataset)}
}

// NewDatasetList returns the dataset list resource, for terraform query.
func NewDatasetList() list.ListResource {
	return &listResource{crudResource{model: modelDataset}}
}

func (r *datasetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a ZFS filesystem dataset. Properties set to `INHERIT` (the default for most) follow the parent dataset. Destroying the resource fails while the dataset has children.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the dataset to create. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"comments": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Comments or description for the dataset. Defaults to `\"INHERIT\"`.",
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"sync": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Synchronous write behavior for the dataset. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("STANDARD", "ALWAYS", "DISABLED", "INHERIT")},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"snapdev": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Controls visibility of volume snapshots under /dev/zvol/.",
				Validators:          []validator.String{stringvalidator.OneOf("HIDDEN", "VISIBLE", "INHERIT")},
			},
			"compression": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Compression algorithm to use for the dataset. Higher numbered variants provide better compression     but use more CPU. 'INHERIT' uses the parent dataset's setting. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("ON", "OFF", "LZ4", "GZIP", "GZIP-1", "GZIP-9", "ZSTD", "ZSTD-FAST", "ZLE", "LZJB", "ZSTD-1", "ZSTD-2", "ZSTD-3", "ZSTD-4", "ZSTD-5", "ZSTD-6", "ZSTD-7", "ZSTD-8", "ZSTD-9", "ZSTD-10", "ZSTD-11", "ZSTD-12", "ZSTD-13", "ZSTD-14", "ZSTD-15", "ZSTD-16", "ZSTD-17", "ZSTD-18", "ZSTD-19", "ZSTD-FAST-1", "ZSTD-FAST-2", "ZSTD-FAST-3", "ZSTD-FAST-4", "ZSTD-FAST-5", "ZSTD-FAST-6", "ZSTD-FAST-7", "ZSTD-FAST-8", "ZSTD-FAST-9", "ZSTD-FAST-10", "ZSTD-FAST-20", "ZSTD-FAST-30", "ZSTD-FAST-40", "ZSTD-FAST-50", "ZSTD-FAST-60", "ZSTD-FAST-70", "ZSTD-FAST-80", "ZSTD-FAST-90", "ZSTD-FAST-100", "ZSTD-FAST-500", "ZSTD-FAST-1000", "INHERIT")},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"exec": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether files in this dataset can be executed. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("ON", "OFF", "INHERIT")},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"managedby": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Identifies which service or system manages this dataset. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"quota_warning": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Percentage of dataset quota at which to issue a warning. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"quota_critical": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Percentage of dataset quota at which to issue a critical alert. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"refquota_warning": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Percentage of reference quota at which to issue a warning. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"refquota_critical": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Percentage of reference quota at which to issue a critical alert. 0-100 or 'INHERIT'. Defaults to `\"INHERIT\"`.",
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"reservation": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Minimum disk space guaranteed to this dataset and its children in bytes.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"refreservation": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Minimum disk space guaranteed to this dataset itself in bytes.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"special_small_block_size": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Size threshold below which blocks are stored on special vdevs.",
			},
			"copies": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Number of copies of data blocks to maintain for redundancy. Defaults to `\"INHERIT\"`.",
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"snapdir": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Controls visibility of the `.zfs/snapshot` directory. 'DISABLED' hides snapshots, 'VISIBLE' shows them,     'HIDDEN' makes them accessible but not listed. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("DISABLED", "VISIBLE", "HIDDEN", "INHERIT")},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"deduplication": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Deduplication setting. 'ON' enables dedup, 'VERIFY' enables with checksum verification, 'OFF' disables. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("ON", "VERIFY", "OFF", "INHERIT")},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"checksum": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Checksum algorithm to verify data integrity. Higher security algorithms like SHA256 provide better     protection but use more CPU. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("ON", "OFF", "FLETCHER2", "FLETCHER4", "SHA256", "SHA512", "SKEIN", "EDONR", "BLAKE3", "INHERIT")},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"readonly": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether the dataset is read-only. Defaults to `\"INHERIT\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("ON", "OFF", "INHERIT")},
				Default:             stringdefault.StaticString("INHERIT"),
			},
			"share_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optimization type for the dataset based on its intended use. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Validators:          []validator.String{stringvalidator.OneOf("GENERIC", "MULTIPROTOCOL", "NFS", "SMB", "APPS")},
			},
			"encryption_options": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration for encryption of dataset for `name` pool. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
				Attributes: map[string]schema.Attribute{
					"generate_key": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Automatically generate the key to be used for dataset encryption. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
					},
					"pbkdf2iters": schema.NumberAttribute{
						Optional:            true,
						MarkdownDescription: "Number of PBKDF2 iterations for key derivation from passphrase. Higher iterations improve security     against brute force attacks but increase unlock time. Default 350,000 balances security and performance. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 100000}},
					},
					"algorithm": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Encryption algorithm to use for dataset encryption. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						Validators:          []validator.String{stringvalidator.OneOf("AES-128-CCM", "AES-192-CCM", "AES-256-CCM", "AES-128-GCM", "AES-192-GCM", "AES-256-GCM")},
					},
					"passphrase": schema.StringAttribute{
						Optional: true, Sensitive: true, WriteOnly: true,
						MarkdownDescription: "Must be specified if encryption for root dataset is desired with a passphrase as a key.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(8)},
					},
					"key": schema.StringAttribute{
						Optional: true, Sensitive: true, WriteOnly: true,
						MarkdownDescription: "A hex-encoded key specified as an alternative to using `passphrase`.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(64), stringvalidator.LengthAtMost(64)},
					},
				},
			},
			"encryption": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Create a ZFS encrypted root dataset for `name` pool.\nThere is 1 case where ZFS encryption is not allowed for a dataset:\n1) If the parent dataset is encrypted with a passphrase and `name` is being created with a key for encrypting the        dataset. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
			},
			"inherit_encryption": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to inherit encryption settings from the parent dataset. Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
			},
			"create_ancestors": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to create any missing parent datasets. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
			},
			"aclmode": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "How Access Control Lists are handled when chmod is used.",
				Validators:          []validator.String{stringvalidator.OneOf("PASSTHROUGH", "RESTRICTED", "DISCARD", "INHERIT")},
			},
			"acltype": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The type of Access Control List system to use.",
				Validators:          []validator.String{stringvalidator.OneOf("OFF", "NFSV4", "POSIX", "INHERIT")},
			},
			"atime": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether file access times are updated when files are accessed.",
				Validators:          []validator.String{stringvalidator.OneOf("ON", "OFF", "INHERIT")},
			},
			"casesensitivity": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "File name case sensitivity setting. Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.OneOf("SENSITIVE", "INSENSITIVE", "INHERIT")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"quota": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum disk space this dataset and its children can consume in bytes.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"refquota": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum disk space this dataset itself can consume in bytes.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"recordsize": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The suggested block size for files in this filesystem dataset.",
			},
			"pool": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the ZFS pool containing this dataset.",
			},
			"encrypted": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the dataset is encrypted.",
			},
			"encryption_root": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The root dataset where encryption is enabled. `null` if the dataset is not encrypted.",
			},
			"key_loaded": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the encryption key is currently loaded for encrypted datasets. `null` for unencrypted datasets.",
			},
			"locked": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether an encrypted dataset is currently locked (key not loaded).",
			},
			"xattr": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Extended attributes storage method (on/off).",
			},
			"compressratio": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The achieved compression ratio as a decimal (e.g., '2.50x').",
			},
			"origin": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The snapshot from which this clone was created. Empty for non-clone datasets.",
			},
			"key_format": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Format of the encryption key (hex/raw/passphrase). Only relevant for encrypted datasets.",
			},
			"encryption_algorithm": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Encryption algorithm used (e.g., AES-256-GCM). Only relevant for encrypted datasets.",
			},
			"used": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Total amount of disk space consumed by this dataset and all its children.",
			},
			"usedbychildren": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Amount of disk space consumed by child datasets.",
			},
			"usedbydataset": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Amount of disk space consumed by this dataset itself, excluding children and snapshots.",
			},
			"usedbyrefreservation": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Amount of disk space consumed by the refreservation of this dataset.",
			},
			"usedbysnapshots": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Amount of disk space consumed by snapshots of this dataset.",
			},
			"available": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Amount of disk space available to this dataset and its children.",
			},
			"pbkdf2iters": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Number of PBKDF2 iterations used for passphrase-based encryption keys.",
			},
			"creation": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Timestamp when this dataset was created.",
			},
			"mountpoint": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Filesystem path where this dataset is mounted. Null for unmounted datasets or volumes.",
			},
			"encryption_options_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `encryption_options` again. Terraform never stores `encryption_options`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
