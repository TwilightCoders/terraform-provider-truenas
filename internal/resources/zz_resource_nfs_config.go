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

var modelNFSConfig = &model{
	typeName:     "nfs_config",
	namespace:    "nfs",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	getMethod:    "nfs.config",
	updateMethod: "nfs.update",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "servers", api: "servers", path: "servers",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Specify the number of nfsd. Default: Number of nfsd is equal number of CPU. ",
		},
		{
			name: "allow_nonroot", api: "allow_nonroot", path: "allow_nonroot",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Allow non-root mount requests.  This equates to 'insecure' share option. ",
		},
		{
			name: "protocols", api: "protocols", path: "protocols",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Specify supported NFS protocols:  NFSv3, NFSv4 or both can be listed. ",
			elem: &node{
				name: "", api: "", path: "protocols",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "v4_krb", api: "v4_krb", path: "v4_krb",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Force Kerberos authentication on NFS shares. ",
		},
		{
			name: "v4_domain", api: "v4_domain", path: "v4_domain",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Specify a DNS domain (NFSv4 only). ",
		},
		{
			name: "bindip", api: "bindip", path: "bindip",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Limit the server IP addresses available for NFS. ",
			elem: &node{
				name: "", api: "", path: "bindip",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "mountd_port", api: "mountd_port", path: "mountd_port",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Specify the mountd port binding. ",
		},
		{
			name: "rpcstatd_port", api: "rpcstatd_port", path: "rpcstatd_port",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Specify the rpc.statd port binding. ",
		},
		{
			name: "rpclockd_port", api: "rpclockd_port", path: "rpclockd_port",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Specify the rpc.lockd port binding. ",
		},
		{
			name: "mountd_log", api: "mountd_log", path: "mountd_log",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable or disable mountd logging. ",
		},
		{
			name: "statd_lockd_log", api: "statd_lockd_log", path: "statd_lockd_log",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable or disable statd and lockd logging. ",
		},
		{
			name: "userd_manage_gids", api: "userd_manage_gids", path: "userd_manage_gids",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable to allow server to manage gids. ",
		},
		{
			name: "rdma", api: "rdma", path: "rdma",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable or disable NFS over RDMA.  Requires RDMA capable NIC. ",
		},
		{
			name: "v4_krb_enabled", api: "v4_krb_enabled", path: "v4_krb_enabled",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Status of NFSv4 authentication requirement (status only). ",
		},
		{
			name: "keytab_has_nfs_spn", api: "keytab_has_nfs_spn", path: "keytab_has_nfs_spn",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Report status of NFS Principal Name in keytab (status only). ",
		},
		{
			name: "managed_nfsd", api: "managed_nfsd", path: "managed_nfsd",
			kind: kindBool, role: roleComputed,
			readable: true,
			description: `Report status of 'servers' field.
If true, the number of nfsd is managed by the server (status only). `,
		},
	},
}

// NewNFSConfig returns the nfs_config resource.
func NewNFSConfig() resource.Resource {
	return &nFSConfigResource{crudResource{model: modelNFSConfig}}
}

type nFSConfigResource struct{ crudResource }

// NewNFSConfigDataSource returns the nfs_config data source.
func NewNFSConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelNFSConfig, attrs: dataAttrs(modelNFSConfig)}
}

func (r *nFSConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update NFS Service Configuration.\n\n`servers` - Represents number of servers to create.\n            By default, the number of nfsd is determined by the capabilities of the system.\n            To specify the number of nfsd, set a value between 1 and 256.\n            'Unset' the field to return to default.\n            This field will always report the number of nfsd to start.\n\n            INPUT: 1 .. 256 or 'unset'\n                where unset will enable the automatic determination\n                and 1 ..256 will set the number of nfsd\n            Default: Number of nfsd is automatically determined and will be no less\n                than 1 and no more than 32\n\n            The number of mountd will be 1/4 the number of reported nfsd.\n\n`allow_nonroot` - If 'enabled' it allows non-root mount requests to be served.\n\n                INPUT: enable/disable (True/False)\n                Default: disabled\n\n`bindip` -  Limit the server IP addresses available for NFS\n            By default, NFS will listen on all IP addresses that are active on the server.\n            To specify the server interface or a set of interfaces provide a list of IP's.\n            If the field is unset/empty, NFS listens on all available server addresses.\n\n            INPUT: list of IP addresses available configured on the server\n            Default: Use all available addresses (empty list)\n\n`protocols` - enable/disable NFSv3, NFSv4\n            Both can be enabled or NFSv4 or NFSv4 by themselves.  At least one must be enabled.\n            Note:  The 'showmount' command is available only if NFSv3 is enabled.\n\n            INPUT: Select NFSv3 or NFSv4 or NFSv3,NFSv4\n            Default: NFSv3,NFSv4\n\n`v4_krb` -  Force Kerberos authentication on NFS shares\n            If enabled, NFS shares will fail if the Kerberos ticket is unavilable\n\n            INPUT: enable/disable\n            Default: disabled\n\n`v4_domain` -   Specify a DNS domain (NFSv4 only)\n            If set, the value will be used to override the default DNS domain name for NFSv4.\n            Specifies the 'Domain' idmapd.conf setting.\n\n            INPUT: a string\n            Default: unset, i.e. an empty string.\n\n`mountd_port` - mountd port binding\n            The value set specifies the port mountd(8) binds to.\n\n            INPUT: unset or an integer between 1 .. 65535\n            Default: unset\n\n`rpcstatd_port` - statd port binding\n            The value set specifies the port rpc.statd(8) binds to.\n\n            INPUT: unset or an integer between 1 .. 65535\n            Default: unset\n\n`rpclockd_port` - lockd port binding\n            The value set specifies the port rpclockd_port(8) binds to.\n\n            INPUT: unset or an integer between 1 .. 65535\n            Default: unset\n\n`rdma` -    Enable/Disable NFS over RDMA support\n            Available on supported platforms and requires an installed and RDMA capable NIC.\n            NFS over RDMA uses port 20040.\n\n            INPUT: Enable/Disable\n            Default: Disable\n\nThis resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"servers": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specify the number of nfsd. Default: Number of nfsd is equal number of CPU. ",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 256}},
			},
			"allow_nonroot": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Allow non-root mount requests.  This equates to 'insecure' share option. ",
			},
			"protocols": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specify supported NFS protocols:  NFSv3, NFSv4 or both can be listed. ",
				ElementType:         types.StringType,
			},
			"v4_krb": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Force Kerberos authentication on NFS shares. ",
			},
			"v4_domain": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specify a DNS domain (NFSv4 only). ",
			},
			"bindip": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Limit the server IP addresses available for NFS. ",
				ElementType:         types.StringType,
			},
			"mountd_port": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specify the mountd port binding. ",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"rpcstatd_port": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specify the rpc.statd port binding. ",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"rpclockd_port": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Specify the rpc.lockd port binding. ",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"mountd_log": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable or disable mountd logging. ",
			},
			"statd_lockd_log": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable or disable statd and lockd logging. ",
			},
			"userd_manage_gids": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable to allow server to manage gids. ",
			},
			"rdma": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable or disable NFS over RDMA.  Requires RDMA capable NIC. ",
			},
			"v4_krb_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Status of NFSv4 authentication requirement (status only). ",
			},
			"keytab_has_nfs_spn": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Report status of NFS Principal Name in keytab (status only). ",
			},
			"managed_nfsd": schema.BoolAttribute{
				Computed: true,
				MarkdownDescription: `Report status of 'servers' field.
If true, the number of nfsd is managed by the server (status only). `,
			},
		},
	}
}
