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

var modelISCSITarget = &model{
	typeName:     "iscsi_target",
	namespace:    "iscsi.target",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "iscsi.target.create",
	updateMethod: "iscsi.target.update",
	getMethod:    "iscsi.target.get_instance",
	deleteMethod: "iscsi.target.delete",
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
			description: "Name of the iSCSI target (maximum 120 characters).",
		},
		{
			name: "alias", api: "alias", path: "alias",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Optional alias name for the iSCSI target.",
		},
		{
			name: "mode", api: "mode", path: "mode",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "ISCSI",
			description: "Protocol mode for the target.\n\n* `ISCSI`: iSCSI protocol only\n* `FC`: Fibre Channel protocol only\n* `BOTH`: Both iSCSI and Fibre Channel protocols\n\nFibre Channel may only be selected on TrueNAS Enterprise-licensed systems with a suitable Fibre Channel HBA. Defaults to `\"ISCSI\"`.",
		},
		{
			name: "groups", api: "groups", path: "groups",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []interface{}{},
			description: "Array of portal-initiator group associations for this target. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "groups",
				kind: kindObject, role: roleRequired,
				readable: true,
				children: []*node{
					{
						name: "portal", api: "portal", path: "groups.portal",
						kind: kindInt, role: roleRequired,
						readable:    true,
						description: "ID of the iSCSI portal to use for this target group.",
					},
					{
						name: "initiator", api: "initiator", path: "groups.initiator",
						kind: kindInt, role: roleOptional,
						nullable: true, readable: true,
						description: "ID of the authorized initiator group or `null` to allow any initiator.",
					},
					{
						name: "authmethod", api: "authmethod", path: "groups.authmethod",
						kind: kindString, role: roleOptionalComputed,
						readable:    true,
						description: "Authentication method for this target group.",
					},
					{
						name: "auth", api: "auth", path: "groups.auth",
						kind: kindInt, role: roleOptional,
						nullable: true, readable: true,
						description: "ID of the authentication credential or `null` if no authentication.",
					},
				},
			},
		},
		{
			name: "auth_networks", api: "auth_networks", path: "auth_networks",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []interface{}{},
			description: "Array of network addresses allowed to access this target. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "auth_networks",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "iscsi_parameters", api: "iscsi_parameters", path: "iscsi_parameters",
			kind: kindObject, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Optional iSCSI-specific parameters for this target.",
			children: []*node{
				{
					name: "queuedcommands", api: "QueuedCommands", path: "iscsi_parameters.QueuedCommands",
					kind: kindInt, role: roleOptional,
					nullable: true, readable: true, updatable: true,
					description: "Maximum number of queued commands per iSCSI session.\n\n* `32`: Standard queue depth for most use cases\n* `128`: Higher queue depth for performance-critical applications",
				},
			},
		},
		{
			name: "rel_tgt_id", api: "rel_tgt_id", path: "rel_tgt_id",
			kind: kindInt, role: roleComputed,
			readable:    true,
			description: "Relative target ID number assigned by the system.",
		},
	},
}

// NewISCSITarget returns the iscsi_target resource.
func NewISCSITarget() resource.Resource {
	return &iSCSITargetResource{crudResource{model: modelISCSITarget}}
}

type iSCSITargetResource struct{ crudResource }

// NewISCSITargetDataSource returns the iscsi_target data source.
func NewISCSITargetDataSource() datasource.DataSource {
	return &dataSource{model: modelISCSITarget, attrs: dataAttrs(modelISCSITarget)}
}

// NewISCSITargetList returns the iscsi_target list resource, for terraform query.
func NewISCSITargetList() list.ListResource {
	return &listResource{crudResource{model: modelISCSITarget}}
}

func (r *iSCSITargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Target.\n\n`groups` is a list of group dictionaries which provide information related to using a `portal`, `initiator`,\n`authmethod` and `auth` with this target. `auth` represents a valid iSCSI Authorized Access and defaults to\nnull.\n\n`auth_networks` is a list of IP/CIDR addresses which are allowed to use this initiator. If all networks are\nto be allowed, this field should be left empty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the iSCSI target (maximum 120 characters).",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.LengthAtMost(120)},
			},
			"alias": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional alias name for the iSCSI target.",
			},
			"mode": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Protocol mode for the target.\n\n* `ISCSI`: iSCSI protocol only\n* `FC`: Fibre Channel protocol only\n* `BOTH`: Both iSCSI and Fibre Channel protocols\n\nFibre Channel may only be selected on TrueNAS Enterprise-licensed systems with a suitable Fibre Channel HBA. Defaults to `\"ISCSI\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("ISCSI", "FC", "BOTH")},
				Default:             stringdefault.StaticString("ISCSI"),
			},
			"groups": schema.ListNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of portal-initiator group associations for this target. Defaults to `[]`.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"portal": schema.NumberAttribute{
							Required:            true,
							MarkdownDescription: "ID of the iSCSI portal to use for this target group.",
							Validators:          []validator.Number{numberIsInteger{}},
						},
						"initiator": schema.NumberAttribute{
							Optional:            true,
							MarkdownDescription: "ID of the authorized initiator group or `null` to allow any initiator.",
							Validators:          []validator.Number{numberIsInteger{}},
						},
						"authmethod": schema.StringAttribute{
							Optional: true, Computed: true,
							MarkdownDescription: "Authentication method for this target group.",
							Validators:          []validator.String{stringvalidator.OneOf("NONE", "CHAP", "CHAP_MUTUAL")},
						},
						"auth": schema.NumberAttribute{
							Optional:            true,
							MarkdownDescription: "ID of the authentication credential or `null` if no authentication.",
							Validators:          []validator.Number{numberIsInteger{}},
						},
					},
				},
			},
			"auth_networks": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of network addresses allowed to access this target. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"iscsi_parameters": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Optional iSCSI-specific parameters for this target.",
				Attributes: map[string]schema.Attribute{
					"queuedcommands": schema.NumberAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum number of queued commands per iSCSI session.\n\n* `32`: Standard queue depth for most use cases\n* `128`: Higher queue depth for performance-critical applications",
						Validators:          []validator.Number{numberIsInteger{}},
					},
				},
			},
			"rel_tgt_id": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "Relative target ID number assigned by the system.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
