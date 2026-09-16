// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelISCSIGlobalConfig = &model{
	typeName:     "iscsi_global_config",
	namespace:    "iscsi.global",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	getMethod:    "iscsi.global.config",
	updateMethod: "iscsi.global.update",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "basename", api: "basename", path: "basename",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Base name prefix for iSCSI target IQNs.",
		},
		{
			name: "isns_servers", api: "isns_servers", path: "isns_servers",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Array of iSNS (Internet Storage Name Service) server addresses.",
			elem: &node{
				name: "", api: "", path: "isns_servers",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "listen_port", api: "listen_port", path: "listen_port",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "TCP port number for iSCSI connections.",
		},
		{
			name: "pool_avail_threshold", api: "pool_avail_threshold", path: "pool_avail_threshold",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Pool available space threshold percentage or `null` to disable.",
		},
		{
			name: "alua", api: "alua", path: "alua",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether Asymmetric Logical Unit Access (ALUA) is enabled. Enabling is limited to TrueNAS Enterprise-licensed     high availability systems. ALUA only works when configured on both the client and server.",
		},
		{
			name: "iser", api: "iser", path: "iser",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether iSCSI Extensions for RDMA (iSER) are enabled. Enabling is limited to TrueNAS Enterprise-licensed     systems and requires the system and network environment have Remote Direct Memory Access (RDMA)-capable hardware.",
		},
	},
}

// NewISCSIGlobalConfig returns the iscsi_global_config resource.
func NewISCSIGlobalConfig() resource.Resource {
	return &iSCSIGlobalConfigResource{crudResource{model: modelISCSIGlobalConfig}}
}

type iSCSIGlobalConfigResource struct{ crudResource }

// NewISCSIGlobalConfigDataSource returns the iscsi_global_config data source.
func NewISCSIGlobalConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelISCSIGlobalConfig, attrs: dataAttrs(modelISCSIGlobalConfig)}
}

func (r *iSCSIGlobalConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "`alua` is a no-op for FreeNAS.\n\nThis resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"basename": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Base name prefix for iSCSI target IQNs.",
			},
			"isns_servers": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of iSNS (Internet Storage Name Service) server addresses.",
				ElementType:         types.StringType,
			},
			"listen_port": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "TCP port number for iSCSI connections.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1025}, numberAtMost{max: 65535}},
			},
			"pool_avail_threshold": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Pool available space threshold percentage or `null` to disable.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 99}},
			},
			"alua": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether Asymmetric Logical Unit Access (ALUA) is enabled. Enabling is limited to TrueNAS Enterprise-licensed     high availability systems. ALUA only works when configured on both the client and server.",
			},
			"iser": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether iSCSI Extensions for RDMA (iSER) are enabled. Enabling is limited to TrueNAS Enterprise-licensed     systems and requires the system and network environment have Remote Direct Memory Access (RDMA)-capable hardware.",
			},
		},
	}
}
