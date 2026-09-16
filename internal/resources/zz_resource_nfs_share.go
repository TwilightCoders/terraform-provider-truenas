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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelNFSShare = &model{
	typeName:     "nfs_share",
	namespace:    "sharing.nfs",
	primaryKey:   "id",
	idKind:       kindInt,
	deleteMethod: "sharing.nfs.delete",
	createMethod: "sharing.nfs.create",
	updateMethod: "sharing.nfs.update",
	getMethod:    "sharing.nfs.get_instance",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "path", api: "path", path: "path",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Local path to be exported. ",
		},
		{
			name: "aliases", api: "aliases", path: "aliases",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []interface{}{},
			description: "IGNORED for now.  Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "aliases",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "comment", api: "comment", path: "comment",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "User comment associated with share.  Defaults to `\"\"`.",
		},
		{
			name: "networks", api: "networks", path: "networks",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []interface{}{},
			description: "List of authorized networks that are allowed to access the share having format     \"network/mask\" CIDR notation. Each entry must be unique. If empty, all networks are allowed.\nExcessively long lists should be avoided. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "networks",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "hosts", api: "hosts", path: "hosts",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []interface{}{},
			description: "List of IP's/hostnames which are allowed to access the share. No quotes or spaces are allowed.\nEach entry must be unique. If empty, all IP's/hostnames are allowed.\nExcessively long lists should be avoided. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "hosts",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "ro", api: "ro", path: "ro",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Export the share as read only.  Defaults to `false`.",
		},
		{
			name: "maproot_user", api: "maproot_user", path: "maproot_user",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Map root user client to a specified user. ",
		},
		{
			name: "maproot_group", api: "maproot_group", path: "maproot_group",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Map root group client to a specified group. ",
		},
		{
			name: "mapall_user", api: "mapall_user", path: "mapall_user",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Map all client users to a specified user. ",
		},
		{
			name: "mapall_group", api: "mapall_group", path: "mapall_group",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Map all client groups to a specified group. ",
		},
		{
			name: "security", api: "security", path: "security",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []interface{}{},
			description: "Specify the security schema.  Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "security",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Enable or disable the share.  Defaults to `true`.",
		},
		{
			name: "expose_snapshots", api: "expose_snapshots", path: "expose_snapshots",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Enterprise feature to enable access to the ZFS snapshot directory for the export.\nExport path must be the root directory of a ZFS dataset. Defaults to `false`.",
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
	},
}

// NewNFSShare returns the nfs_share resource.
func NewNFSShare() resource.Resource {
	return &nFSShareResource{crudResource{model: modelNFSShare}}
}

type nFSShareResource struct{ crudResource }

// NewNFSShareDataSource returns the nfs_share data source.
func NewNFSShareDataSource() datasource.DataSource {
	return &dataSource{model: modelNFSShare, attrs: dataAttrs(modelNFSShare)}
}

// NewNFSShareList returns the nfs_share list resource, for terraform query.
func NewNFSShareList() list.ListResource {
	return &listResource{crudResource{model: modelNFSShare}}
}

func (r *nFSShareResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a NFS Share.\n\n`path` local path to be exported.\n\n`aliases` IGNORED, for now.\n\n`networks` is a list of authorized networks that are allowed to access the share having format\n\"network/mask\" CIDR notation. If empty, all networks are allowed.\n\n`hosts` is a list of IP's/hostnames which are allowed to access the share. If empty, all IP's/hostnames are\nallowed.\n\n`expose_snapshots` enable TrueNAS Enterprise feature to allow access\nto the ZFS snapshot directory over NFS. This feature requires a valid\nenterprise license.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Local path to be exported. ",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"aliases": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "IGNORED for now.  Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"comment": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "User comment associated with share.  Defaults to `\"\"`.",
				Default:             stringdefault.StaticString(""),
			},
			"networks": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "List of authorized networks that are allowed to access the share having format     \"network/mask\" CIDR notation. Each entry must be unique. If empty, all networks are allowed.\nExcessively long lists should be avoided. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"hosts": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "List of IP's/hostnames which are allowed to access the share. No quotes or spaces are allowed.\nEach entry must be unique. If empty, all IP's/hostnames are allowed.\nExcessively long lists should be avoided. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"ro": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Export the share as read only.  Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"maproot_user": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Map root user client to a specified user. ",
			},
			"maproot_group": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Map root group client to a specified group. ",
			},
			"mapall_user": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Map all client users to a specified user. ",
			},
			"mapall_group": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Map all client groups to a specified group. ",
			},
			"security": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specify the security schema.  Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable or disable the share.  Defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"expose_snapshots": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enterprise feature to enable access to the ZFS snapshot directory for the export.\nExport path must be the root directory of a ZFS dataset. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"locked": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: `Read-only value indicating whether the share is located on a locked dataset.

Returns:
    - True: The share is in a locked dataset.
    - False: The share is not in a locked dataset.
    - None: Lock status is not available because path locking information was not requested.`,
			},
		},
	}
}
