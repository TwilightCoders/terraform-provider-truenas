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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var modelISCSITargetExtent = &model{
	typeName:     "iscsi_target_extent",
	namespace:    "iscsi.targetextent",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "iscsi.targetextent.create",
	updateMethod: "iscsi.targetextent.update",
	getMethod:    "iscsi.targetextent.get_instance",
	deleteMethod: "iscsi.targetextent.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "target", api: "target", path: "target",
			kind: kindInt, role: roleRequired,
			readable: true, updatable: true,
			description: "ID of the iSCSI target to associate with the extent.",
		},
		{
			name: "lunid", api: "lunid", path: "lunid",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "LUN ID to assign or `null` to auto-assign the next available LUN.",
		},
		{
			name: "extent", api: "extent", path: "extent",
			kind: kindInt, role: roleRequired,
			readable: true, updatable: true,
			description: "ID of the iSCSI extent to associate with the target.",
		},
	},
}

// NewISCSITargetExtent returns the iscsi_target_extent resource.
func NewISCSITargetExtent() resource.Resource {
	return &iSCSITargetExtentResource{crudResource{model: modelISCSITargetExtent}}
}

type iSCSITargetExtentResource struct{ crudResource }

// NewISCSITargetExtentDataSource returns the iscsi_target_extent data source.
func NewISCSITargetExtentDataSource() datasource.DataSource {
	return &dataSource{model: modelISCSITargetExtent, attrs: dataAttrs(modelISCSITargetExtent)}
}

// NewISCSITargetExtentList returns the iscsi_target_extent list resource, for terraform query.
func NewISCSITargetExtentList() list.ListResource {
	return &listResource{crudResource{model: modelISCSITargetExtent}}
}

func (r *iSCSITargetExtentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an Associated Target.\n\n`lunid` will be automatically assigned if it is not provided based on the `target`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"target": schema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "ID of the iSCSI target to associate with the extent.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"lunid": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "LUN ID to assign or `null` to auto-assign the next available LUN.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"extent": schema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "ID of the iSCSI extent to associate with the target.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
