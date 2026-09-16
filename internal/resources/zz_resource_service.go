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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var modelService = &model{
	typeName:     "service",
	namespace:    "service",
	primaryKey:   "id",
	idKind:       kindInt,
	adopt:        true,
	adoptBy:      "service",
	adoptByAttr:  "service",
	updateMethod: "service.update",
	getMethod:    "service.get_instance",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "service", api: "service", path: "service",
			kind: kindString, role: roleRequired,
			replace: true, readable: true,
			description: "Name of the system service. Changing this forces a new resource.",
		},
		{
			name: "enable", api: "enable", path: "enable",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether the service should start on boot.",
		},
		{
			name: "state", api: "state", path: "state",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Current state of the service (e.g., 'RUNNING', 'STOPPED').",
		},
	},
}

// NewService returns the service resource.
func NewService() resource.Resource {
	return &serviceResource{crudResource{model: modelService}}
}

type serviceResource struct{ crudResource }

// NewServiceDataSource returns the service data source.
func NewServiceDataSource() datasource.DataSource {
	return &dataSource{model: modelService, attrs: dataAttrs(modelService)}
}

// NewServiceList returns the service list resource, for terraform query.
func NewServiceList() list.ListResource {
	return &listResource{crudResource{model: modelService}}
}

func (r *serviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages whether a system service (for example `cifs`, `nfs`, `ssh`) starts at boot. The service must already exist. Start or stop it immediately with the `truenas_service_control` action.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"service": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the system service. Changing this forces a new resource.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"enable": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether the service should start on boot.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current state of the service (e.g., 'RUNNING', 'STOPPED').",
			},
		},
	}
}
