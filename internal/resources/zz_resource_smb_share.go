// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
)

var modelSMBShare = &model{
	typeName:     "smb_share",
	namespace:    "sharing.smb",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "sharing.smb.create",
	updateMethod: "sharing.smb.update",
	getMethod:    "sharing.smb.get_instance",
	deleteMethod: "sharing.smb.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "purpose", api: "purpose", path: "purpose",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "DEFAULT_SHARE",
			description: "This parameter sets the purpose of the SMB share. It controls how the SMB share behaves and what features are     available through options. The DEFAULT_SHARE setting is best for most applications, and should be used, unless     there is a specific reason to change it.\n\n* `DEFAULT_SHARE`: Set the SMB share for best compatibility with common SMB clients.\n\n* `LEGACY_SHARE`: Set the SMB share for compatibility with older TrueNAS versions. Automated backend migrations       use this to help the administrator move to better-supported share settings. It should not be used for new SMB       shares.\n\n* `TIMEMACHINE_SHARE`: The SMB share is presented to MacOS clients as a time machine target.\n  NOTE: `aapl_extensions` must be set in the global `smb.config`.\n\n* `MULTIPROTOCOL_SHARE`: The SMB share is configured for multi-protocol access. Set this if the `path` is shared       through NFS, FTP, or used by containers or apps.\n  NOTE: This setting can reduce SMB share performance because it turns off some SMB features for safer       interoperability with external processes.\n\n* `TIME_LOCKED_SHARE`: The SMB share makes files read-only through the SMB protocol after the set grace_period       ends.\n  WARNING: This setting does not work if the `path` is accessed locally or if another SMB share without the       `TIME_LOCKED_SHARE` purpose uses the same path.\n  WARNING: This setting might not meet regulatory requirements for write-once storage.\n\n* `PRIVATE_DATASETS_SHARE`: The server uses the specified `dataset_naming_schema` in `options` to make a new ZFS       dataset when the client connects. The server uses this dataset as the share path during the SMB session.\n\n* `EXTERNAL_SHARE`: The SMB share is a DFS proxy to a share hosted on an external SMB server.\n\n* `VEEAM_REPOSITORY_SHARE`: The SMB share is a repository for Veeam Backup & Replication and supports Fast Clone.\n  NOTE: This feature is available only for TrueNAS Enterprise customers.\n\n* `FCP_SHARE`: The SMB share is a used for Final Cut Pro storage. This feature automatically configures the share       to provide storage according to Apple support guidelines described in https://support.apple.com/en-ca/101919.       NOTE: `aapl_extensions` must be set in the global `smb.config`.       WARNING: This feature forcibly enables `aapl_name_mangling` on the SMB share which may cause unexpected behavior       for data that was written without this feature enabled. Defaults to `\"DEFAULT_SHARE\"`.",
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "SMB share name. SMB share names are case-insensitive and must be unique, and are subject     to the following restrictions:\n\n* A share name must be no more than 80 characters in length.\n\n* The following characters are illegal in a share name: `\\ / [ ] : | < > + = ; , * ? \"`\n\n* Unicode control characters are illegal in a share name.\n\n* The following share names are not allowed: global, printers, homes.",
		},
		{
			name: "path", api: "path", path: "path",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Local server path to share by using the SMB protocol. The path must start with `/mnt/` and must be in a     ZFS pool.\n\nUse the string `EXTERNAL` if the share works as a DFS proxy.\n\nWARNING: The TrueNAS server does not check if external paths are reachable. ",
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "If unset, the SMB share is not available over the SMB protocol.  Defaults to `true`.",
		},
		{
			name: "comment", api: "comment", path: "comment",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Text field that is seen next to a share when an SMB client requests a list of SMB shares on the TrueNAS     server.  Defaults to `\"\"`.",
		},
		{
			name: "readonly", api: "readonly", path: "readonly",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "If set, SMB clients cannot create or change files and directories in the SMB share.\n\nNOTE: If set, the share path is still writeable by local processes or other file sharing protocols.  Defaults to `false`.",
		},
		{
			name: "browsable", api: "browsable", path: "browsable",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "If set, the share is included when an SMB client requests a list of SMB shares on the TrueNAS server.  Defaults to `true`.",
		},
		{
			name: "access_based_share_enumeration", api: "access_based_share_enumeration", path: "access_based_share_enumeration",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "If set, the share is only included when an SMB client requests a list of shares on the SMB server if     the share (not filesystem) access control list (see `sharing.smb.getacl`) grants access to the user.  Defaults to `false`.",
		},
		{
			name: "audit", api: "audit", path: "audit",
			kind: kindObject, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: map[string]any{"enable": false, "ignore_list": []any{}, "watch_list": []any{}},
			description: "Audit configuration for monitoring SMB share access and operations. Defaults to `{\"enable\":false,\"ignore_list\":[],\"watch_list\":[]}`.",
			children: []*node{
				{
					name: "enable", api: "enable", path: "audit.enable",
					kind: kindBool, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: false,
					description: "Turn on auditing for the SMB share. SMB share auditing may not be enabled if `enable_smb1` is `true`     in the SMB service configuration. Defaults to `false`.",
				},
				{
					name: "watch_list", api: "watch_list", path: "audit.watch_list",
					kind: kindList, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: []any{},
					description: "Only audit the listed group accounts. If the list is empty, all groups will be audited.  Defaults to `[]`.",
					elem: &node{
						name: "", api: "", path: "audit.watch_list",
						kind: kindString, role: roleRequired,
						readable: true,
					},
				},
				{
					name: "ignore_list", api: "ignore_list", path: "audit.ignore_list",
					kind: kindList, role: roleOptionalComputed,
					readable: true, updatable: true,
					hasDefault: true, def: []any{},
					description: "List of groups that will not be audited.  Defaults to `[]`.",
					elem: &node{
						name: "", api: "", path: "audit.ignore_list",
						kind: kindString, role: roleRequired,
						readable: true,
					},
				},
			},
		},
		{
			name: "options", api: "options", path: "options",
			kind: kindUnion, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description:   "Additional configuration related to the configured SMB share purpose. If null, then the default     options related to the share purpose will be applied. ",
			discriminator: "purpose",
			children: []*node{
				{
					name: "legacy_share", api: "LEGACY_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `LEGACY_SHARE`.",
					children: []*node{
						{
							name: "recyclebin", api: "recyclebin", path: "options.recyclebin",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, deleted files are moved to per-user subdirectories in the `.recycle` directory. The     SMB server creates the `.recycle` directory at the root of the SMB share if the file is in the same     ZFS dataset as the share `path`. If the file is in a child ZFS dataset, the server uses the     `mountpoint` of that dataset to create the `.recycle` directory.\n\nNOTE: This feature does not work with recycle bin features in client operating systems.\n\nWARNING: Do not use this feature instead of backups or ZFS snapshots.  Defaults to `false`.",
						},
						{
							name: "path_suffix", api: "path_suffix", path: "options.path_suffix",
							kind: kindString, role: roleOptionalComputed,
							nullable: true, readable: true,
							description: "Path suffix template for dynamic path generation. Uses SMB variable substitution patterns like `%D` (domain)     and `%U` (username).",
						},
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "guestok", api: "guestok", path: "options.guestok",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, guest access to the share is allowed. This should not be used in production environments.\n\nNOTE: If a user account does not exist, the SMB server maps access to the guest account.\n\nWARNING: Additional client-side configuration downgrading security settings may be required in order     to use this feature.  Defaults to `false`.",
						},
						{
							name: "streams", api: "streams", path: "options.streams",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: true,
							description: "If set, support for SMB alternate data streams is enabled.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `true`.",
						},
						{
							name: "durablehandle", api: "durablehandle", path: "options.durablehandle",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: true,
							description: "If set, support for SMB durable handles is enabled.\n\nWARNING: This feature is incompatible with multiprotocol and local filesystem access.  Defaults to `true`.",
						},
						{
							name: "shadowcopy", api: "shadowcopy", path: "options.shadowcopy",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: true,
							description: "If set, previous versions of files contained in ZFS snapshots are accessible through standard SMB protocol     operations on previous versions of files.  Defaults to `true`.",
						},
						{
							name: "fsrvp", api: "fsrvp", path: "options.fsrvp",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, enable support for the File Server Remote VSS Protocol. This allows clients to manage     snapshots for the specified SMB share.  Defaults to `false`.",
						},
						{
							name: "home", api: "home", path: "options.home",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Use the `path` to store user home directories. Each user has a personal home directory and share.     Users cannot access other user directories when connecting to shares.\n\nNOTE: This parameter changes the share `name` to `homes`. It also creates a dynamic share that mirrors     the username of the user. Both shares use the same `path`. You can hide the homes share by turning off     `browsable`. The dynamic user home share cannot be hidden.\n\nWARNING: This parameter changes the global server configuration. The SMB server will not authenticate     users without a valid home directory or shell. Defaults to `false`.",
						},
						{
							name: "acl", api: "acl", path: "options.acl",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: true,
							description: "If set, enable mapping of local filesystem ACLs to NT ACLs for SMB clients.  Defaults to `true`.",
						},
						{
							name: "afp", api: "afp", path: "options.afp",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, SMB server will read and store file metadata in an on-disk format compatible with the     legacy AFP file server.\n\nWARNING: This should not be set unless the SMB server is sharing data that was originally written     via the AFP protocol.  Defaults to `false`.",
						},
						{
							name: "timemachine", api: "timemachine", path: "options.timemachine",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, MacOS clients can use the share as a time machine target.  Defaults to `false`.",
						},
						{
							name: "timemachine_quota", api: "timemachine_quota", path: "options.timemachine_quota",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: 0,
							description: "If set, it defines the maximum size of a single time machine sparsebundle volume by limiting the     reported disk size to the SMB client. A value of zero means no quota is applied to the share.\n\nNOTE: Modern MacOS versions you set Time Machine quotas client-side. This gives more predictable     server and client behavior. Defaults to `0`.",
						},
						{
							name: "aapl_name_mangling", api: "aapl_name_mangling", path: "options.aapl_name_mangling",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
						},
						{
							name: "vuid", api: "vuid", path: "options.vuid",
							kind: kindString, role: roleOptionalComputed,
							nullable: true, readable: true,
							description: "This value is the Time Machine volume UUID for the SMB share. The TrueNAS server uses this value in the mDNS     advertisement for the Time Machine share. MacOS clients may use it to identify the volume. When you create or     update a share, setting this value to null makes the TrueNAS server generate a new UUID for the share. ",
						},
						{
							name: "auxsmbconf", api: "auxsmbconf", path: "options.auxsmbconf",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "",
							description: "Additional parameters to set on the SMB share. Parameters must be separated by the new-line character.\n\nWARNING: These parameters are not validated and may cause undefined server behavior including     data corruption or data loss.\n\nWARNING: Auxiliary parameters are an unsupported configuration. Defaults to `\"\"`.",
						},
					},
				},
				{
					name: "default_share", api: "DEFAULT_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `DEFAULT_SHARE`.",
					children: []*node{
						{
							name: "aapl_name_mangling", api: "aapl_name_mangling", path: "options.aapl_name_mangling",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
						},
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
				{
					name: "timemachine_share", api: "TIMEMACHINE_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `TIMEMACHINE_SHARE`.",
					children: []*node{
						{
							name: "timemachine_quota", api: "timemachine_quota", path: "options.timemachine_quota",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: 0,
							description: "If set, it defines the maximum size in bytes of a single time machine sparsebundle volume by limiting the     reported disk size to the SMB client. A value of zero means no quota is set.\n\nNOTE: Modern MacOS versions you set Time Machine quotas client-side. This gives more predictable     server and client behavior. Defaults to `0`.",
						},
						{
							name: "auto_snapshot", api: "auto_snapshot", path: "options.auto_snapshot",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, the server makes a ZFS snapshot of the share dataset when the client makes a new     Time Machine backup.  Defaults to `false`.",
						},
						{
							name: "auto_dataset_creation", api: "auto_dataset_creation", path: "options.auto_dataset_creation",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, the server uses the `dataset_naming_schema` to make a new ZFS dataset when the client connects.     The server uses this dataset as the share path during the SMB session.\n\nNOTE: this setting requires the share path to be a dataset mountpoint. Defaults to `false`.",
						},
						{
							name: "dataset_naming_schema", api: "dataset_naming_schema", path: "options.dataset_naming_schema",
							kind: kindString, role: roleOptionalComputed,
							nullable: true, readable: true,
							description: "The naming schema to use when `auto_dataset_creation` is specified. If you do not set a schema,     the server uses `%U` (username) if it is not joined to Active Directory. If the server is joined to     Active Directory it uses `%D/%U` (domain/username). See the `VARIABLE SUBSTITUTIONS` section in the smb.conf     manpage for valid strings.\n\nWARNING: ZFS dataset naming rules are more restrictive than normal path rules. For example, if `%u` is specified     then the character `\\` may be inserted in the username (which is not supported in ZFS).",
						},
						{
							name: "vuid", api: "vuid", path: "options.vuid",
							kind: kindString, role: roleOptionalComputed,
							nullable: true, readable: true,
							description: "This value is the Time Machine volume UUID for the SMB share. The TrueNAS server uses this value in the mDNS     advertisement for the Time Machine share. MacOS clients may use it to identify the volume. When you create or     update a share, setting this value to null makes the TrueNAS server generate a new UUID for the share. ",
						},
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
				{
					name: "multiprotocol_share", api: "MULTIPROTOCOL_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `MULTIPROTOCOL_SHARE`.",
					children: []*node{
						{
							name: "aapl_name_mangling", api: "aapl_name_mangling", path: "options.aapl_name_mangling",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
						},
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
				{
					name: "time_locked_share", api: "TIME_LOCKED_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `TIME_LOCKED_SHARE`.",
					children: []*node{
						{
							name: "grace_period", api: "grace_period", path: "options.grace_period",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: 900,
							description: "Time in seconds when write access to the file or directory is allowed.  Defaults to `900`.",
						},
						{
							name: "aapl_name_mangling", api: "aapl_name_mangling", path: "options.aapl_name_mangling",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
						},
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
				{
					name: "private_datasets_share", api: "PRIVATE_DATASETS_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `PRIVATE_DATASETS_SHARE`.",
					children: []*node{
						{
							name: "dataset_naming_schema", api: "dataset_naming_schema", path: "options.dataset_naming_schema",
							kind: kindString, role: roleOptionalComputed,
							nullable: true, readable: true,
							description: "The naming schema to use. If you do not set a schema, the server uses `%U` (username) if it is not joined to     Active Directory. If the server is joined to Active Directory it uses `%D/%U` (domain/username).\n\nWARNING: ZFS dataset naming rules are more restrictive than normal path rules.",
						},
						{
							name: "auto_quota", api: "auto_quota", path: "options.auto_quota",
							kind: kindInt, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: 0,
							description: "Set the specified ZFS quota (in gibibytes) on new datasets. If the value is zero, TrueNAS disables     automatic quotas for the share. Defaults to `0`.",
						},
						{
							name: "aapl_name_mangling", api: "aapl_name_mangling", path: "options.aapl_name_mangling",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
						},
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
				{
					name: "external_share", api: "EXTERNAL_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `EXTERNAL_SHARE`.",
					children: []*node{
						{
							name: "remote_path", api: "remote_path", path: "options.remote_path",
							kind: kindList, role: roleRequired,
							readable:    true,
							description: "This is the path to the external server and share. Each server entry must include a full domain name or IP     address and share name. Separate the server and share with the `\\` character.\n\nWARNING: The SMB server and TrueNAS middleware do not check if external paths are reachable. ",
							elem: &node{
								name: "", api: "", path: "options.remote_path",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
				{
					name: "veeam_repository_share", api: "VEEAM_REPOSITORY_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `VEEAM_REPOSITORY_SHARE`.",
					children: []*node{
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
				{
					name: "fcp_share", api: "FCP_SHARE", path: "options",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `purpose` is `FCP_SHARE`.",
					children: []*node{
						{
							name: "aapl_name_mangling", api: "aapl_name_mangling", path: "options.aapl_name_mangling",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: true,
							description: "Illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.  Defaults to `true`.",
						},
						{
							name: "hostsallow", api: "hostsallow", path: "options.hostsallow",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsallow",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
						{
							name: "hostsdeny", api: "hostsdeny", path: "options.hostsdeny",
							kind: kindList, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: []any{},
							description: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
							elem: &node{
								name: "", api: "", path: "options.hostsdeny",
								kind: kindString, role: roleRequired,
								readable: true,
							},
						},
					},
				},
			},
		},
		{
			name: "locked", api: "locked", path: "locked",
			kind: kindBool, role: roleComputed,
			nullable: true, readable: true,
			description: `Read-only value indicating whether the share is located on a locked dataset.

Returns:
    - True: The share is in a locked dataset.
    - False: The share is not in a locked dataset.
    - None: Lock status is not available because path locking information was not requested.`,
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

// NewSMBShare returns the smb_share resource.
func NewSMBShare() resource.Resource {
	return &sMBShareResource{crudResource{model: modelSMBShare}}
}

type sMBShareResource struct{ crudResource }

// NewSMBShareDataSource returns the smb_share data source.
func NewSMBShareDataSource() datasource.DataSource {
	return &dataSource{model: modelSMBShare, attrs: dataAttrs(modelSMBShare)}
}

// NewSMBShareList returns the smb_share list resource, for terraform query.
func NewSMBShareList() list.ListResource {
	return &listResource{crudResource{model: modelSMBShare}}
}

func (r *sMBShareResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"purpose": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "This parameter sets the purpose of the SMB share. It controls how the SMB share behaves and what features are     available through options. The DEFAULT_SHARE setting is best for most applications, and should be used, unless     there is a specific reason to change it.\n\n* `DEFAULT_SHARE`: Set the SMB share for best compatibility with common SMB clients.\n\n* `LEGACY_SHARE`: Set the SMB share for compatibility with older TrueNAS versions. Automated backend migrations       use this to help the administrator move to better-supported share settings. It should not be used for new SMB       shares.\n\n* `TIMEMACHINE_SHARE`: The SMB share is presented to MacOS clients as a time machine target.\n  NOTE: `aapl_extensions` must be set in the global `smb.config`.\n\n* `MULTIPROTOCOL_SHARE`: The SMB share is configured for multi-protocol access. Set this if the `path` is shared       through NFS, FTP, or used by containers or apps.\n  NOTE: This setting can reduce SMB share performance because it turns off some SMB features for safer       interoperability with external processes.\n\n* `TIME_LOCKED_SHARE`: The SMB share makes files read-only through the SMB protocol after the set grace_period       ends.\n  WARNING: This setting does not work if the `path` is accessed locally or if another SMB share without the       `TIME_LOCKED_SHARE` purpose uses the same path.\n  WARNING: This setting might not meet regulatory requirements for write-once storage.\n\n* `PRIVATE_DATASETS_SHARE`: The server uses the specified `dataset_naming_schema` in `options` to make a new ZFS       dataset when the client connects. The server uses this dataset as the share path during the SMB session.\n\n* `EXTERNAL_SHARE`: The SMB share is a DFS proxy to a share hosted on an external SMB server.\n\n* `VEEAM_REPOSITORY_SHARE`: The SMB share is a repository for Veeam Backup & Replication and supports Fast Clone.\n  NOTE: This feature is available only for TrueNAS Enterprise customers.\n\n* `FCP_SHARE`: The SMB share is a used for Final Cut Pro storage. This feature automatically configures the share       to provide storage according to Apple support guidelines described in https://support.apple.com/en-ca/101919.       NOTE: `aapl_extensions` must be set in the global `smb.config`.       WARNING: This feature forcibly enables `aapl_name_mangling` on the SMB share which may cause unexpected behavior       for data that was written without this feature enabled. Defaults to `\"DEFAULT_SHARE\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("DEFAULT_SHARE", "LEGACY_SHARE", "TIMEMACHINE_SHARE", "MULTIPROTOCOL_SHARE", "TIME_LOCKED_SHARE", "PRIVATE_DATASETS_SHARE", "EXTERNAL_SHARE", "VEEAM_REPOSITORY_SHARE", "FCP_SHARE")},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "SMB share name. SMB share names are case-insensitive and must be unique, and are subject     to the following restrictions:\n\n* A share name must be no more than 80 characters in length.\n\n* The following characters are illegal in a share name: `\\ / [ ] : | < > + = ; , * ? \"`\n\n* Unicode control characters are illegal in a share name.\n\n* The following share names are not allowed: global, printers, homes.",
			},
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Local server path to share by using the SMB protocol. The path must start with `/mnt/` and must be in a     ZFS pool.\n\nUse the string `EXTERNAL` if the share works as a DFS proxy.\n\nWARNING: The TrueNAS server does not check if external paths are reachable. ",
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If unset, the SMB share is not available over the SMB protocol.  Defaults to `true`.",
			},
			"comment": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Text field that is seen next to a share when an SMB client requests a list of SMB shares on the TrueNAS     server.  Defaults to `\"\"`.",
			},
			"readonly": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If set, SMB clients cannot create or change files and directories in the SMB share.\n\nNOTE: If set, the share path is still writeable by local processes or other file sharing protocols.  Defaults to `false`.",
			},
			"browsable": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If set, the share is included when an SMB client requests a list of SMB shares on the TrueNAS server.  Defaults to `true`.",
			},
			"access_based_share_enumeration": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If set, the share is only included when an SMB client requests a list of shares on the SMB server if     the share (not filesystem) access control list (see `sharing.smb.getacl`) grants access to the user.  Defaults to `false`.",
			},
			"audit": schema.SingleNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Audit configuration for monitoring SMB share access and operations. Defaults to `{\"enable\":false,\"ignore_list\":[],\"watch_list\":[]}`.",
				Attributes: map[string]schema.Attribute{
					"enable": schema.BoolAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Turn on auditing for the SMB share. SMB share auditing may not be enabled if `enable_smb1` is `true`     in the SMB service configuration. Defaults to `false`.",
						Default:             booldefault.StaticBool(false),
					},
					"watch_list": schema.ListAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Only audit the listed group accounts. If the list is empty, all groups will be audited.  Defaults to `[]`.",
						ElementType:         types.StringType,
					},
					"ignore_list": schema.ListAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "List of groups that will not be audited.  Defaults to `[]`.",
						ElementType:         types.StringType,
					},
				},
			},
			"options": schema.SingleNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional configuration related to the configured SMB share purpose. If null, then the default     options related to the share purpose will be applied. ",
				Attributes: map[string]schema.Attribute{
					"legacy_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `LEGACY_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"recyclebin": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, deleted files are moved to per-user subdirectories in the `.recycle` directory. The     SMB server creates the `.recycle` directory at the root of the SMB share if the file is in the same     ZFS dataset as the share `path`. If the file is in a child ZFS dataset, the server uses the     `mountpoint` of that dataset to create the `.recycle` directory.\n\nNOTE: This feature does not work with recycle bin features in client operating systems.\n\nWARNING: Do not use this feature instead of backups or ZFS snapshots.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"path_suffix": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Path suffix template for dynamic path generation. Uses SMB variable substitution patterns like `%D` (domain)     and `%U` (username).",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"guestok": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, guest access to the share is allowed. This should not be used in production environments.\n\nNOTE: If a user account does not exist, the SMB server maps access to the guest account.\n\nWARNING: Additional client-side configuration downgrading security settings may be required in order     to use this feature.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"streams": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, support for SMB alternate data streams is enabled.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `true`.",
								Default:             booldefault.StaticBool(true),
							},
							"durablehandle": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, support for SMB durable handles is enabled.\n\nWARNING: This feature is incompatible with multiprotocol and local filesystem access.  Defaults to `true`.",
								Default:             booldefault.StaticBool(true),
							},
							"shadowcopy": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, previous versions of files contained in ZFS snapshots are accessible through standard SMB protocol     operations on previous versions of files.  Defaults to `true`.",
								Default:             booldefault.StaticBool(true),
							},
							"fsrvp": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, enable support for the File Server Remote VSS Protocol. This allows clients to manage     snapshots for the specified SMB share.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"home": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Use the `path` to store user home directories. Each user has a personal home directory and share.     Users cannot access other user directories when connecting to shares.\n\nNOTE: This parameter changes the share `name` to `homes`. It also creates a dynamic share that mirrors     the username of the user. Both shares use the same `path`. You can hide the homes share by turning off     `browsable`. The dynamic user home share cannot be hidden.\n\nWARNING: This parameter changes the global server configuration. The SMB server will not authenticate     users without a valid home directory or shell. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"acl": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, enable mapping of local filesystem ACLs to NT ACLs for SMB clients.  Defaults to `true`.",
								Default:             booldefault.StaticBool(true),
							},
							"afp": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, SMB server will read and store file metadata in an on-disk format compatible with the     legacy AFP file server.\n\nWARNING: This should not be set unless the SMB server is sharing data that was originally written     via the AFP protocol.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"timemachine": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, MacOS clients can use the share as a time machine target.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"timemachine_quota": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, it defines the maximum size of a single time machine sparsebundle volume by limiting the     reported disk size to the SMB client. A value of zero means no quota is applied to the share.\n\nNOTE: Modern MacOS versions you set Time Machine quotas client-side. This gives more predictable     server and client behavior. Defaults to `0`.",
								Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 1.099511627776e+14}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(0)),
							},
							"aapl_name_mangling": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"vuid": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "This value is the Time Machine volume UUID for the SMB share. The TrueNAS server uses this value in the mDNS     advertisement for the Time Machine share. MacOS clients may use it to identify the volume. When you create or     update a share, setting this value to null makes the TrueNAS server generate a new UUID for the share. ",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"auxsmbconf": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Additional parameters to set on the SMB share. Parameters must be separated by the new-line character.\n\nWARNING: These parameters are not validated and may cause undefined server behavior including     data corruption or data loss.\n\nWARNING: Auxiliary parameters are an unsupported configuration. Defaults to `\"\"`.",
								Default:             stringdefault.StaticString(""),
							},
						},
					},
					"default_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `DEFAULT_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"aapl_name_mangling": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
						},
					},
					"timemachine_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `TIMEMACHINE_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"timemachine_quota": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, it defines the maximum size in bytes of a single time machine sparsebundle volume by limiting the     reported disk size to the SMB client. A value of zero means no quota is set.\n\nNOTE: Modern MacOS versions you set Time Machine quotas client-side. This gives more predictable     server and client behavior. Defaults to `0`.",
								Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 1.099511627776e+14}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(0)),
							},
							"auto_snapshot": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, the server makes a ZFS snapshot of the share dataset when the client makes a new     Time Machine backup.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"auto_dataset_creation": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, the server uses the `dataset_naming_schema` to make a new ZFS dataset when the client connects.     The server uses this dataset as the share path during the SMB session.\n\nNOTE: this setting requires the share path to be a dataset mountpoint. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"dataset_naming_schema": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "The naming schema to use when `auto_dataset_creation` is specified. If you do not set a schema,     the server uses `%U` (username) if it is not joined to Active Directory. If the server is joined to     Active Directory it uses `%D/%U` (domain/username). See the `VARIABLE SUBSTITUTIONS` section in the smb.conf     manpage for valid strings.\n\nWARNING: ZFS dataset naming rules are more restrictive than normal path rules. For example, if `%u` is specified     then the character `\\` may be inserted in the username (which is not supported in ZFS).",
							},
							"vuid": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "This value is the Time Machine volume UUID for the SMB share. The TrueNAS server uses this value in the mDNS     advertisement for the Time Machine share. MacOS clients may use it to identify the volume. When you create or     update a share, setting this value to null makes the TrueNAS server generate a new UUID for the share. ",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
						},
					},
					"multiprotocol_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `MULTIPROTOCOL_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"aapl_name_mangling": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
						},
					},
					"time_locked_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `TIME_LOCKED_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"grace_period": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Time in seconds when write access to the file or directory is allowed.  Defaults to `900`.",
								Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 60}, numberAtMost{max: 1.5552e+07}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(900)),
							},
							"aapl_name_mangling": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
						},
					},
					"private_datasets_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `PRIVATE_DATASETS_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"dataset_naming_schema": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "The naming schema to use. If you do not set a schema, the server uses `%U` (username) if it is not joined to     Active Directory. If the server is joined to Active Directory it uses `%D/%U` (domain/username).\n\nWARNING: ZFS dataset naming rules are more restrictive than normal path rules.",
							},
							"auto_quota": schema.NumberAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Set the specified ZFS quota (in gibibytes) on new datasets. If the value is zero, TrueNAS disables     automatic quotas for the share. Defaults to `0`.",
								Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}},
								Default:             numberdefault.StaticBigFloat(big.NewFloat(0)),
							},
							"aapl_name_mangling": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "If set, illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.\n\nWARNING: This value should not be changed once data is written to the SMB share.  Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
						},
					},
					"external_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `EXTERNAL_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"remote_path": schema.ListAttribute{
								Required:            true,
								MarkdownDescription: "This is the path to the external server and share. Each server entry must include a full domain name or IP     address and share name. Separate the server and share with the `\\` character.\n\nWARNING: The SMB server and TrueNAS middleware do not check if external paths are reachable. ",
								ElementType:         types.StringType,
							},
						},
					},
					"veeam_repository_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `VEEAM_REPOSITORY_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
						},
					},
					"fcp_share": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `purpose` is `FCP_SHARE`.",
						Attributes: map[string]schema.Attribute{
							"aapl_name_mangling": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Illegal NTFS characters commonly used by MacOS clients are stored with their native values on the SMB     server's local filesystem.\n\nNOTE: Files with illegal NTFS characters in their names may not be accessible to non-MacOS SMB clients.  Defaults to `true`.",
								Default:             booldefault.StaticBool(true),
							},
							"hostsallow": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are allowed to access the SMB share. The EXCEPT keyword     may be used to limit a wildcard list.\n\nNOTE: Hostname lookups are disabled on the SMB server for performance reasons.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
							"hostsdeny": schema.ListAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "A list of IP addresses or subnets that are not allowed to access the SMB share. The keyword     `ALL` or the netmask `0.0.0.0/0` may be used to deny all by default.  Defaults to `[]`.",
								ElementType:         types.StringType,
							},
						},
					},
				},
			},
			"locked": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: `Read-only value indicating whether the share is located on a locked dataset.

Returns:
    - True: The share is in a locked dataset.
    - False: The share is not in a locked dataset.
    - None: Lock status is not available because path locking information was not requested.`,
			},
			"unset": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
				ElementType:         types.StringType,
			},
		},
	}
}
