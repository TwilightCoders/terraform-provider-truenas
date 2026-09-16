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

var modelStaticRoute = &model{
	typeName:     "static_route",
	namespace:    "staticroute",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "staticroute.create",
	updateMethod: "staticroute.update",
	getMethod:    "staticroute.get_instance",
	deleteMethod: "staticroute.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "destination", api: "destination", path: "destination",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Destination network or host for this static route.",
		},
		{
			name: "gateway", api: "gateway", path: "gateway",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Gateway IP address for this static route.",
		},
		{
			name: "description", api: "description", path: "description",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Optional description for this static route. Defaults to `\"\"`.",
		},
	},
}

// NewStaticRoute returns the static_route resource.
func NewStaticRoute() resource.Resource {
	return &staticRouteResource{crudResource{model: modelStaticRoute}}
}

type staticRouteResource struct{ crudResource }

// NewStaticRouteDataSource returns the static_route data source.
func NewStaticRouteDataSource() datasource.DataSource {
	return &dataSource{model: modelStaticRoute, attrs: dataAttrs(modelStaticRoute)}
}

// NewStaticRouteList returns the static_route list resource, for terraform query.
func NewStaticRouteList() list.ListResource {
	return &listResource{crudResource{model: modelStaticRoute}}
}

func (r *staticRouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Static Route.\n\nAddress families of `gateway` and `destination` should match when creating a static route.\n\n`description` is an optional attribute for any notes regarding the static route.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"destination": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Destination network or host for this static route.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"gateway": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Gateway IP address for this static route.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Optional description for this static route. Defaults to `\"\"`.",
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}
