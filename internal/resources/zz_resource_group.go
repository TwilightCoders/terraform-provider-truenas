// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelGroup = &model{
	typeName:     "group",
	namespace:    "group",
	primaryKey:   "id",
	idKind:       kindInt,
	getMethod:    "group.get_instance",
	deleteMethod: "group.delete",
	createMethod: "group.create",
	updateMethod: "group.update",
	listFilters:  [][]interface{}{[]interface{}{"builtin", "=", false}},
	listOptions:  map[string]interface{}(nil),
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "gid", api: "gid", path: "gid",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "If `null`, it is automatically filled with the next one available. Changing this forces a new resource.",
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "A string used to identify a group.",
		},
		{
			name: "sudo_commands", api: "sudo_commands", path: "sudo_commands",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "A list of commands that group members may execute with elevated privileges. User is prompted for password     when executing any command from the list.  Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "sudo_commands",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "sudo_commands_nopasswd", api: "sudo_commands_nopasswd", path: "sudo_commands_nopasswd",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "A list of commands that group members may execute with elevated privileges. User is not prompted for password     when executing any command from the list.  Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "sudo_commands_nopasswd",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "smb", api: "smb", path: "smb",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "If set to `True`, the group can be used for SMB share ACL entries. The group is mapped to an NT group account     on the TrueNAS SMB server and has a `sid` value.  Defaults to `true`.",
		},
		{
			name: "userns_idmap", api: "userns_idmap", path: "userns_idmap",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true, inherit: true, intOrString: true,
			description: `Specifies the subgid mapping for this group. If DIRECT then the GID will be     directly mapped to all containers. Alternatively, the target GID may be     explicitly specified. If null, then the GID will not be mapped.

**NOTE: This field will be ignored for groups that have been assigned TrueNAS roles.**`,
			dialect: apiDialect,
		},
		{
			name: "users", api: "users", path: "users",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "A list a API user identifiers for local users who are members of this group. These IDs match the `id` field     from `user.query`.\n\nNOTE: This field is empty for groups that come from directory services (`local` is `False`).  Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "users",
				kind: kindInt, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "builtin", api: "builtin", path: "builtin",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "If `True`, the group is an internal system account for the TrueNAS server. Typically, one should     create dedicated groups for access to the TrueNAS server webui and shares. ",
		},
		{
			name: "local", api: "local", path: "local",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "If `True`, the group is local to the TrueNAS server. If `False`, the group is provided by a directory service.",
		},
		{
			name: "sid", api: "sid", path: "sid",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses this value to     check share access and for other purposes. ",
		},
		{
			name: "roles", api: "roles", path: "roles",
			kind: kindList, role: roleComputed,
			readable:    true,
			description: "List of roles assigned to this groups. Roles control administrative access to TrueNAS through the web UI and     API. You can change group roles by using `privilege.create`, `privilege.update`, and `privilege.delete`. ",
			elem: &node{
				name: "", api: "", path: "roles",
				kind: kindString, role: roleComputed,
				readable: true,
			},
		},
		{
			name: "immutable", api: "immutable", path: "immutable",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "This is a read-only field showing if the group entry can be changed. If `True`, the group is immutable and     cannot be changed. If `False`, the group can be changed. ",
		},
	},
}

// NewGroup returns the group resource.
func NewGroup() resource.Resource {
	return &groupResource{crudResource{model: modelGroup}}
}

type groupResource struct{ crudResource }

// NewGroupDataSource returns the group data source.
func NewGroupDataSource() datasource.DataSource {
	return &dataSource{model: modelGroup, attrs: dataAttrs(modelGroup)}
}

// NewGroupList returns the group list resource, for terraform query.
func NewGroupList() list.ListResource {
	return &listResource{crudResource{model: modelGroup}}
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"gid": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If `null`, it is automatically filled with the next one available. Changing this forces a new resource.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 9e+07}},
				PlanModifiers:       []planmodifier.Number{numberplanmodifier.RequiresReplace(), numberplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A string used to identify a group.",
			},
			"sudo_commands": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A list of commands that group members may execute with elevated privileges. User is prompted for password     when executing any command from the list.  Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A list of commands that group members may execute with elevated privileges. User is not prompted for password     when executing any command from the list.  Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"smb": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If set to `True`, the group can be used for SMB share ACL entries. The group is mapped to an NT group account     on the TrueNAS SMB server and has a `sid` value.  Defaults to `true`.",
			},
			"userns_idmap": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: `Specifies the subgid mapping for this group. If DIRECT then the GID will be     directly mapped to all containers. Alternatively, the target GID may be     explicitly specified. If null, then the GID will not be mapped.

**NOTE: This field will be ignored for groups that have been assigned TrueNAS roles.**`,
			},
			"users": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "A list a API user identifiers for local users who are members of this group. These IDs match the `id` field     from `user.query`.\n\nNOTE: This field is empty for groups that come from directory services (`local` is `False`).  Defaults to `[]`.",
				ElementType:         types.NumberType,
			},
			"builtin": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If `True`, the group is an internal system account for the TrueNAS server. Typically, one should     create dedicated groups for access to the TrueNAS server webui and shares. ",
			},
			"local": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If `True`, the group is local to the TrueNAS server. If `False`, the group is provided by a directory service.",
			},
			"sid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses this value to     check share access and for other purposes. ",
			},
			"roles": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "List of roles assigned to this groups. Roles control administrative access to TrueNAS through the web UI and     API. You can change group roles by using `privilege.create`, `privilege.update`, and `privilege.delete`. ",
				ElementType:         types.StringType,
			},
			"immutable": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "This is a read-only field showing if the group entry can be changed. If `True`, the group is immutable and     cannot be changed. If `False`, the group can be changed. ",
			},
		},
	}
}
