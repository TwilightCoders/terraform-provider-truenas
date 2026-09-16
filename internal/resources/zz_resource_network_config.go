// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelNetworkConfig = &model{
	typeName:     "network_config",
	namespace:    "network.configuration",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	updateMethod: "network.configuration.update",
	getMethod:    "network.configuration.config",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "hostname", api: "hostname", path: "hostname",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "System hostname.",
		},
		{
			name: "domain", api: "domain", path: "domain",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "System domain name.",
		},
		{
			name: "ipv4gateway", api: "ipv4gateway", path: "ipv4gateway",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Used instead of the default gateway provided by DHCP.",
		},
		{
			name: "ipv6gateway", api: "ipv6gateway", path: "ipv6gateway",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "IPv6 default gateway address.",
		},
		{
			name: "nameserver1", api: "nameserver1", path: "nameserver1",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Primary DNS server.",
		},
		{
			name: "nameserver2", api: "nameserver2", path: "nameserver2",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Secondary DNS server.",
		},
		{
			name: "nameserver3", api: "nameserver3", path: "nameserver3",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Tertiary DNS server.",
		},
		{
			name: "httpproxy", api: "httpproxy", path: "httpproxy",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Must be provided if a proxy is to be used for network operations.",
		},
		{
			name: "hosts", api: "hosts", path: "hosts",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Static host entries to add to the hosts file.",
			elem: &node{
				name: "", api: "", path: "hosts",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "domains", api: "domains", path: "domains",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional domain names for DNS search.",
			elem: &node{
				name: "", api: "", path: "domains",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "service_announcement", api: "service_announcement", path: "service_announcement",
			kind: kindObject, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Determines the broadcast protocols that will be used to advertise the server.",
			children: []*node{
				{
					name: "netbios", api: "netbios", path: "service_announcement.netbios",
					kind: kindBool, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Enable the NetBIOS name server (NBNS) which starts concurrently with the SMB service. SMB clients will only     perform NBNS lookups if SMB1 is enabled. NBNS may be required for legacy SMB clients.",
				},
				{
					name: "mdns", api: "mdns", path: "service_announcement.mdns",
					kind: kindBool, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Enable multicast DNS service announcements for enabled services.",
				},
				{
					name: "wsd", api: "wsd", path: "service_announcement.wsd",
					kind: kindBool, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Enable Web Service Discovery support.",
				},
			},
		},
		{
			name: "activity", api: "activity", path: "activity",
			kind: kindObject, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Network activity filtering configuration.",
			children: []*node{
				{
					name: "type", api: "type", path: "activity.type",
					kind: kindString, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Whether to allow or deny the specified network activities.",
				},
				{
					name: "activities", api: "activities", path: "activity.activities",
					kind: kindList, role: roleOptionalComputed,
					readable: true, updatable: true,
					description: "Array of network activity types to allow or deny.",
					elem: &node{
						name: "", api: "", path: "activity.activities",
						kind: kindString, role: roleRequired,
						readable: true,
					},
				},
			},
		},
		{
			name: "hostname_b", api: "hostname_b", path: "hostname_b",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Hostname for the second controller in HA configurations or `null`.",
		},
		{
			name: "hostname_virtual", api: "hostname_virtual", path: "hostname_virtual",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Virtual hostname for HA configurations or `null`.",
		},
		{
			name: "hostname_local", api: "hostname_local", path: "hostname_local",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Local hostname for this system.",
		},
	},
}

// NewNetworkConfig returns the network_config resource.
func NewNetworkConfig() resource.Resource {
	return &networkConfigResource{crudResource{model: modelNetworkConfig}}
}

type networkConfigResource struct{ crudResource }

// NewNetworkConfigDataSource returns the network_config data source.
func NewNetworkConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelNetworkConfig, attrs: dataAttrs(modelNetworkConfig)}
}

func (r *networkConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Update Network Configuration Service configuration.

This resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"hostname": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "System hostname.",
			},
			"domain": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "System domain name.",
			},
			"ipv4gateway": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Used instead of the default gateway provided by DHCP.",
			},
			"ipv6gateway": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "IPv6 default gateway address.",
			},
			"nameserver1": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Primary DNS server.",
			},
			"nameserver2": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Secondary DNS server.",
			},
			"nameserver3": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Tertiary DNS server.",
			},
			"httpproxy": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Must be provided if a proxy is to be used for network operations.",
			},
			"hosts": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Static host entries to add to the hosts file.",
				ElementType:         types.StringType,
			},
			"domains": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional domain names for DNS search.",
				ElementType:         types.StringType,
			},
			"service_announcement": schema.SingleNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Determines the broadcast protocols that will be used to advertise the server.",
				Attributes: map[string]schema.Attribute{
					"netbios": schema.BoolAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Enable the NetBIOS name server (NBNS) which starts concurrently with the SMB service. SMB clients will only     perform NBNS lookups if SMB1 is enabled. NBNS may be required for legacy SMB clients.",
					},
					"mdns": schema.BoolAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Enable multicast DNS service announcements for enabled services.",
					},
					"wsd": schema.BoolAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Enable Web Service Discovery support.",
					},
				},
			},
			"activity": schema.SingleNestedAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Network activity filtering configuration.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Whether to allow or deny the specified network activities.",
						Validators:          []validator.String{stringvalidator.OneOf("ALLOW", "DENY")},
					},
					"activities": schema.ListAttribute{
						Optional: true, Computed: true,
						MarkdownDescription: "Array of network activity types to allow or deny.",
						ElementType:         types.StringType,
					},
				},
			},
			"hostname_b": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Hostname for the second controller in HA configurations or `null`.",
			},
			"hostname_virtual": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Virtual hostname for HA configurations or `null`.",
			},
			"hostname_local": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Local hostname for this system.",
			},
		},
	}
}
