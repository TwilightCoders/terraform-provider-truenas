// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelISCSIInitiator = &model{
	typeName:     "iscsi_initiator",
	namespace:    "iscsi.initiator",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "iscsi.initiator.create",
	updateMethod: "iscsi.initiator.update",
	getMethod:    "iscsi.initiator.get_instance",
	deleteMethod: "iscsi.initiator.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "initiators", api: "initiators", path: "initiators",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: []any{},
			description: "Array of iSCSI Qualified Names (IQNs) or IP addresses of authorized initiators. Defaults to `[]`.",
			elem: &node{
				name: "", api: "", path: "initiators",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "comment", api: "comment", path: "comment",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Optional comment describing the authorized initiator group. Defaults to `\"\"`.",
		},
	},
}

// NewISCSIInitiator returns the iscsi_initiator resource.
func NewISCSIInitiator() resource.Resource {
	return &iSCSIInitiatorResource{crudResource{model: modelISCSIInitiator}}
}

type iSCSIInitiatorResource struct{ crudResource }

// NewISCSIInitiatorDataSource returns the iscsi_initiator data source.
func NewISCSIInitiatorDataSource() datasource.DataSource {
	return &dataSource{model: modelISCSIInitiator, attrs: dataAttrs(modelISCSIInitiator)}
}

// NewISCSIInitiatorList returns the iscsi_initiator list resource, for terraform query.
func NewISCSIInitiatorList() list.ListResource {
	return &listResource{crudResource{model: modelISCSIInitiator}}
}

func (r *iSCSIInitiatorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Initiator.\n\n`initiators` is a list of initiator hostnames which are authorized to access an iSCSI Target. To allow all\npossible initiators, `initiators` can be left empty.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"initiators": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of iSCSI Qualified Names (IQNs) or IP addresses of authorized initiators. Defaults to `[]`.",
				ElementType:         types.StringType,
			},
			"comment": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Optional comment describing the authorized initiator group. Defaults to `\"\"`.",
			},
		},
	}
}
