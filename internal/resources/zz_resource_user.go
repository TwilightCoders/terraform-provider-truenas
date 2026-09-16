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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelUser = &model{
	typeName:     "user",
	namespace:    "user",
	primaryKey:   "id",
	idKind:       kindInt,
	getMethod:    "user.get_instance",
	deleteMethod: "user.delete",
	createMethod: "user.create",
	updateMethod: "user.update",
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
			name: "uid", api: "uid", path: "uid",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "UNIX UID. If not provided, it is automatically filled with the next one available. Changing this forces a new resource.",
		},
		{
			name: "username", api: "username", path: "username",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "String used to uniquely identify the user on the server. In order to be portable across     systems, local user names must be composed of characters from the POSIX portable filename     character set (IEEE Std 1003.1-2024 section 3.265). This means alphanumeric characters,     hyphens, underscores, and periods. Usernames also may not begin with a hyphen or a period.",
		},
		{
			name: "home", api: "home", path: "home",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "/var/empty",
			description: "The local file system path for the user account's home directory.\nTypically, this is required only if the account has shell access (local or SSH) to TrueNAS.\nThis is not required for accounts used only for SMB share access.  Defaults to `\"/var/empty\"`.",
		},
		{
			name: "shell", api: "shell", path: "shell",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "/usr/bin/zsh",
			description: "Available choices can be retrieved with `user.shell_choices`. Defaults to `\"/usr/bin/zsh\"`.",
		},
		{
			name: "full_name", api: "full_name", path: "full_name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Comment field to provide additional information about the user account. Typically, this is     the full name of the user or a short description of a service account. There are no character set restrictions     for this field. This field is for information only. ",
		},
		{
			name: "smb", api: "smb", path: "smb",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash of the     user account's password for local accounts. This feature is unavailable for local accounts when General Purpose OS     STIG compatibility mode is enabled. If set to `true` the user is automatically added to the `builtin_users`     group. Defaults to `true`.",
		},
		{
			name: "userns_idmap", api: "userns_idmap", path: "userns_idmap",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true, inherit: true, intOrString: true,
			description: "Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to all containers. Alternatively, the target UID may be     explicitly specified. If `null`, then the UID will not be mapped.\n\nNOTE: This field will be ignored for users that have been assigned TrueNAS roles.",
			dialect:     apiDialect,
		},
		{
			name: "group", api: "group", path: "group",
			kind: kindInt, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "The group entry `id` for the user's primary group. This is not the same as the Unix group `gid` value.     This is required if `group_create` is `false`. ",
			ref:         "id",
		},
		{
			name: "groups", api: "groups", path: "groups",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Array of additional groups to which the user belongs. NOTE: Groups are identified by their group entry `id`,     not their Unix group ID (`gid`). ",
			elem: &node{
				name: "", api: "", path: "groups",
				kind: kindInt, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "password_disabled", api: "password_disabled", path: "password_disabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "If set to `true` password authentication for the user account is disabled.\n\nNOTE: Users with password authentication disabled may still authenticate to the TrueNAS server by other methods,     such as SSH key-based authentication.\n\nNOTE: Password authentication is required for `smb` users. Defaults to `false`.",
		},
		{
			name: "ssh_password_enabled", api: "ssh_password_enabled", path: "ssh_password_enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Allow the user to authenticate to the TrueNAS SSH server using a password.\n\nWARNING: The established best practice is to use only key-based authentication for SSH servers.  Defaults to `false`.",
		},
		{
			name: "sshpubkey", api: "sshpubkey", path: "sshpubkey",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server. ",
		},
		{
			name: "locked", api: "locked", path: "locked",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS server.  Defaults to `false`.",
		},
		{
			name: "sudo_commands", api: "sudo_commands", path: "sudo_commands",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "An array of commands the user may execute with elevated privileges. User is prompted for password     when executing any command from the array. ",
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
			description: "An array of commands the user may execute with elevated privileges. User is *not* prompted for password     when executing any command from the array. ",
			elem: &node{
				name: "", api: "", path: "sudo_commands_nopasswd",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "email", api: "email", path: "email",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and     notifications. ",
		},
		{
			name: "group_create", api: "group_create", path: "group_create",
			kind: kindBool, role: roleOptional,
			createOnly:  true,
			description: "If set to `true`, the TrueNAS server automatically creates a new local group as the user's primary group.  Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
		},
		{
			name: "home_create", api: "home_create", path: "home_create",
			kind: kindBool, role: roleOptional,
			updatable:   true,
			description: "Create a new home directory for the user in the specified `home` path.  TrueNAS does not report this value, so changes made outside Terraform are not detected.",
		},
		{
			name: "home_mode", api: "home_mode", path: "home_mode",
			kind: kindString, role: roleOptional,
			updatable:   true,
			description: "Filesystem permission to set on the user's home directory.  TrueNAS does not report this value, so changes made outside Terraform are not detected.",
		},
		{
			name: "password", api: "password", path: "password",
			kind: kindString, role: roleOptional,
			nullable: true, sensitive: true, writeOnly: true, updatable: true,
			description: "The password for the user account. This is required if `random_password` is not set. ",
		},
		{
			name: "builtin", api: "builtin", path: "builtin",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "If `true`, the user account is an internal system account for the TrueNAS server. Typically, one should     create dedicated user accounts for access to the TrueNAS server webui and shares. ",
		},
		{
			name: "local", api: "local", path: "local",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "If `true`, the account is local to the TrueNAS server. If `false`, the account is provided by a directory     service. ",
		},
		{
			name: "immutable", api: "immutable", path: "immutable",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "If `true`, the account is system-provided and most fields related to it may not be changed. ",
		},
		{
			name: "twofactor_auth_configured", api: "twofactor_auth_configured", path: "twofactor_auth_configured",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "If `true`, the account has been configured for two-factor authentication. Users are prompted for a     second factor when authenticating to the TrueNAS web UI and API. They may also be prompted when signing     in to the TrueNAS SSH server using a password (depending on global two-factor authentication settings). ",
		},
		{
			name: "sid", api: "sid", path: "sid",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses     this value to check share access and for other purposes. ",
		},
		{
			name: "last_password_change", api: "last_password_change", path: "last_password_change",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "The date of the last password change for local user accounts.",
		},
		{
			name: "password_age", api: "password_age", path: "password_age",
			kind: kindInt, role: roleComputed,
			nullable: true, readable: true,
			description: "The age in days of the password for local user accounts.",
		},
		{
			name: "password_history", api: "password_history", path: "password_history",
			kind: kindList, role: roleComputed,
			nullable: true, readable: true,
			description: "This contains hashes of the ten most recent passwords used by local user accounts, and is     for enforcing password history requirements as defined in system.security.",
			elem: &node{
				name: "", api: "", path: "password_history",
				kind: kindAny, role: roleComputed,
				readable: true,
			},
		},
		{
			name: "password_change_required", api: "password_change_required", path: "password_change_required",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Password change for local user account is required on next login.",
		},
		{
			name: "roles", api: "roles", path: "roles",
			kind: kindList, role: roleComputed,
			readable:    true,
			description: "Array of roles assigned to this user's groups. Roles control administrative access to TrueNAS through     the web UI and API.",
			elem: &node{
				name: "", api: "", path: "roles",
				kind: kindString, role: roleComputed,
				readable: true,
			},
		},
		{
			name: "api_keys", api: "api_keys", path: "api_keys",
			kind: kindList, role: roleComputed,
			readable:    true,
			description: "Array of API key IDs associated with this user account for programmatic access.",
			elem: &node{
				name: "", api: "", path: "api_keys",
				kind: kindInt, role: roleComputed,
				readable: true,
			},
		},
		{
			name: "password_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `password` again. Terraform never stores `password`.",
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

// NewUser returns the user resource.
func NewUser() resource.Resource {
	return &userResource{crudResource{model: modelUser}}
}

type userResource struct{ crudResource }

// NewUserDataSource returns the user data source.
func NewUserDataSource() datasource.DataSource {
	return &dataSource{model: modelUser, attrs: dataAttrs(modelUser)}
}

// NewUserList returns the user list resource, for terraform query.
func NewUserList() list.ListResource {
	return &listResource{crudResource{model: modelUser}}
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"uid": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "UNIX UID. If not provided, it is automatically filled with the next one available. Changing this forces a new resource.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 9e+07}},
				PlanModifiers:       []planmodifier.Number{numberplanmodifier.RequiresReplace(), numberplanmodifier.UseStateForUnknown()},
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "String used to uniquely identify the user on the server. In order to be portable across     systems, local user names must be composed of characters from the POSIX portable filename     character set (IEEE Std 1003.1-2024 section 3.265). This means alphanumeric characters,     hyphens, underscores, and periods. Usernames also may not begin with a hyphen or a period.",
			},
			"home": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The local file system path for the user account's home directory.\nTypically, this is required only if the account has shell access (local or SSH) to TrueNAS.\nThis is not required for accounts used only for SMB share access.  Defaults to `\"/var/empty\"`.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"shell": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Available choices can be retrieved with `user.shell_choices`. Defaults to `\"/usr/bin/zsh\"`.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"full_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Comment field to provide additional information about the user account. Typically, this is     the full name of the user or a short description of a service account. There are no character set restrictions     for this field. This field is for information only. ",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"smb": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The user account may be used to access SMB shares. If set to `true` then TrueNAS stores an NT hash of the     user account's password for local accounts. This feature is unavailable for local accounts when General Purpose OS     STIG compatibility mode is enabled. If set to `true` the user is automatically added to the `builtin_users`     group. Defaults to `true`.",
			},
			"userns_idmap": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specifies the subuid mapping for this user. If DIRECT then the UID will be     directly mapped to all containers. Alternatively, the target UID may be     explicitly specified. If `null`, then the UID will not be mapped.\n\nNOTE: This field will be ignored for users that have been assigned TrueNAS roles.",
			},
			"group": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "The group entry `id` for the user's primary group. This is not the same as the Unix group `gid` value.     This is required if `group_create` is `false`. ",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"groups": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of additional groups to which the user belongs. NOTE: Groups are identified by their group entry `id`,     not their Unix group ID (`gid`). ",
				ElementType:         types.NumberType,
			},
			"password_disabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If set to `true` password authentication for the user account is disabled.\n\nNOTE: Users with password authentication disabled may still authenticate to the TrueNAS server by other methods,     such as SSH key-based authentication.\n\nNOTE: Password authentication is required for `smb` users. Defaults to `false`.",
			},
			"ssh_password_enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Allow the user to authenticate to the TrueNAS SSH server using a password.\n\nWARNING: The established best practice is to use only key-based authentication for SSH servers.  Defaults to `false`.",
			},
			"sshpubkey": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "SSH public keys corresponding to private keys that authenticate this user to the TrueNAS SSH server. ",
			},
			"locked": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If set to `true` the account is locked. The account cannot be used to authenticate to the TrueNAS server.  Defaults to `false`.",
			},
			"sudo_commands": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "An array of commands the user may execute with elevated privileges. User is prompted for password     when executing any command from the array. ",
				ElementType:         types.StringType,
			},
			"sudo_commands_nopasswd": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "An array of commands the user may execute with elevated privileges. User is *not* prompted for password     when executing any command from the array. ",
				ElementType:         types.StringType,
			},
			"email": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Email address of the user. If the user has the `FULL_ADMIN` role, they will receive email alerts and     notifications. ",
			},
			"group_create": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "If set to `true`, the TrueNAS server automatically creates a new local group as the user's primary group.  Used only when creating the resource; changing it forces a new resource. TrueNAS does not report it, so after an import the first plan records it without changing TrueNAS.",
			},
			"home_create": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Create a new home directory for the user in the specified `home` path.  TrueNAS does not report this value, so changes made outside Terraform are not detected.",
			},
			"home_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filesystem permission to set on the user's home directory.  TrueNAS does not report this value, so changes made outside Terraform are not detected.",
			},
			"password": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "The password for the user account. This is required if `random_password` is not set. ",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"builtin": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If `true`, the user account is an internal system account for the TrueNAS server. Typically, one should     create dedicated user accounts for access to the TrueNAS server webui and shares. ",
			},
			"local": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If `true`, the account is local to the TrueNAS server. If `false`, the account is provided by a directory     service. ",
			},
			"immutable": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If `true`, the account is system-provided and most fields related to it may not be changed. ",
			},
			"twofactor_auth_configured": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If `true`, the account has been configured for two-factor authentication. Users are prompted for a     second factor when authenticating to the TrueNAS web UI and API. They may also be prompted when signing     in to the TrueNAS SSH server using a password (depending on global two-factor authentication settings). ",
			},
			"sid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Security Identifier (SID) of the user if the account an `smb` account. The SMB server uses     this value to check share access and for other purposes. ",
			},
			"last_password_change": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date of the last password change for local user accounts.",
			},
			"password_age": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "The age in days of the password for local user accounts.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"password_history": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "This contains hashes of the ten most recent passwords used by local user accounts, and is     for enforcing password history requirements as defined in system.security.",
				ElementType:         types.StringType,
			},
			"password_change_required": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Password change for local user account is required on next login.",
			},
			"roles": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "Array of roles assigned to this user's groups. Roles control administrative access to TrueNAS through     the web UI and API.",
				ElementType:         types.StringType,
			},
			"api_keys": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "Array of API key IDs associated with this user account for programmatic access.",
				ElementType:         types.NumberType,
			},
			"password_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `password` again. Terraform never stores `password`.",
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
