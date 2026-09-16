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
)

var modelISCSIPortal = &model{
	typeName:     "iscsi_portal",
	namespace:    "iscsi.portal",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "iscsi.portal.create",
	updateMethod: "iscsi.portal.update",
	getMethod:    "iscsi.portal.get_instance",
	deleteMethod: "iscsi.portal.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "listen", api: "listen", path: "listen",
			kind: kindList, role: roleRequired,
			readable: true, updatable: true,
			description: "Array of IP addresses for the portal to listen on.",
			elem: &node{
				name: "", api: "", path: "listen",
				kind: kindObject, role: roleRequired,
				readable: true,
				children: []*node{
					{
						name: "ip", api: "ip", path: "listen.ip",
						kind: kindString, role: roleRequired,
						readable:    true,
						description: "IP address for the iSCSI portal to listen on.",
					},
					{
						name: "port", api: "port", path: "listen.port",
						kind: kindInt, role: roleComputed,
						readable:    true,
						description: "TCP port number for the iSCSI portal.",
					},
				},
			},
		},
		{
			name: "comment", api: "comment", path: "comment",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Optional comment describing the portal. Defaults to `\"\"`.",
		},
		{
			name: "tag", api: "tag", path: "tag",
			kind: kindInt, role: roleComputed,
			readable:    true,
			description: "Numeric tag used to associate this portal with iSCSI targets.",
		},
	},
}

// NewISCSIPortal returns the iscsi_portal resource.
func NewISCSIPortal() resource.Resource {
	return &iSCSIPortalResource{crudResource{model: modelISCSIPortal}}
}

type iSCSIPortalResource struct{ crudResource }

// NewISCSIPortalDataSource returns the iscsi_portal data source.
func NewISCSIPortalDataSource() datasource.DataSource {
	return &dataSource{model: modelISCSIPortal, attrs: dataAttrs(modelISCSIPortal)}
}

// NewISCSIPortalList returns the iscsi_portal list resource, for terraform query.
func NewISCSIPortalList() list.ListResource {
	return &listResource{crudResource{model: modelISCSIPortal}}
}

func (r *iSCSIPortalResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a new iSCSI Portal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"listen": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Array of IP addresses for the portal to listen on.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "IP address for the iSCSI portal to listen on.",
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"port": schema.NumberAttribute{
							Computed:            true,
							MarkdownDescription: "TCP port number for the iSCSI portal.",
							Validators:          []validator.Number{numberIsInteger{}},
						},
					},
				},
			},
			"comment": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Optional comment describing the portal. Defaults to `\"\"`.",
				Default:             stringdefault.StaticString(""),
			},
			"tag": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "Numeric tag used to associate this portal with iSCSI targets.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
