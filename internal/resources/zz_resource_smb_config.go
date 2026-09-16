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

var modelSMBConfig = &model{
	typeName:     "smb_config",
	namespace:    "smb",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	getMethod:    "smb.config",
	updateMethod: "smb.update",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "netbiosname", api: "netbiosname", path: "netbiosname",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "The NetBIOS name of this server. ",
		},
		{
			name: "netbiosalias", api: "netbiosalias", path: "netbiosalias",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Alternative netbios names of the TrueNAS server. These names are announced through NetBIOS name server and     registered in Active Directory when TrueNAS joins the domain.",
			elem: &node{
				name: "", api: "", path: "netbiosalias",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "workgroup", api: "workgroup", path: "workgroup",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Workgroup name. When TrueNAS joins active directory, it automatically changes this value to match the NetBIOS     domain of the Active Directory domain. ",
		},
		{
			name: "description", api: "description", path: "description",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Description of the SMB server. SMB clients may see this description during some operations. ",
		},
		{
			name: "enable_smb1", api: "enable_smb1", path: "enable_smb1",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable SMB1 support on the server. WARNING: using the SMB1 protocol is not recommended. ",
		},
		{
			name: "unixcharset", api: "unixcharset", path: "unixcharset",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Select character set for file names on local filesystem. Use this option only if you know the names are not     UTF-8. ",
		},
		{
			name: "localmaster", api: "localmaster", path: "localmaster",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "When set to `true` the NetBIOS name server in TrueNAS participates in elections for the local master browser.\nWhen set to `false` the NetBIOS name server does not attempt to become a local master browser on a subnet and     loses all browsing elections.\n\nNOTE: This parameter has no effect if the NetBIOS name server is disabled. ",
		},
		{
			name: "syslog", api: "syslog", path: "syslog",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Send log messages to syslog. Enable this option if you want SMB server error logs to be included in     information sent to a remote syslog server. NOTE: This requires that remote syslog is globally configured on     TrueNAS. ",
		},
		{
			name: "aapl_extensions", api: "aapl_extensions", path: "aapl_extensions",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable support for SMB2/3 AAPL protocol extensions. This setting makes the TrueNAS server advertise support     for Apple protocol extensions as a MacOS server. Enabling this is required for Time Machine support. ",
		},
		{
			name: "admin_group", api: "admin_group", path: "admin_group",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "The selected group has full administrator privileges on TrueNAS via the SMB protocol. ",
		},
		{
			name: "guest", api: "guest", path: "guest",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "SMB guest account username. This username provides access to legacy SMB shares with guest access enabled.     It must be a valid, existing local user account. ",
		},
		{
			name: "filemask", api: "filemask", path: "filemask",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "`smb.conf` create mask. DEFAULT applies current server default which is 664. ",
		},
		{
			name: "dirmask", api: "dirmask", path: "dirmask",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "`smb.conf` directory mask. DEFAULT applies current server default which is 775. ",
		},
		{
			name: "ntlmv1_auth", api: "ntlmv1_auth", path: "ntlmv1_auth",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable legacy and very insecure NTLMv1 authentication. This should never be done except     in extreme edge cases and may be against regulations in non-home environments. ",
		},
		{
			name: "multichannel", api: "multichannel", path: "multichannel",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Enable SMB3 multi-channel support. ",
		},
		{
			name: "encryption", api: "encryption", path: "encryption",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "SMB2/3 transport encryption setting for the TrueNAS SMB server.\n\n* `NEGOTIATE`: Enable negotiation of data encryption. Encrypt data only if the client explicitly requests it.\n* `DESIRED`: Enable negotiation of data encryption. Encrypt data on sessions and share connections for clients       that support it.\n* `REQUIRED`: Require data encryption for sessions and share connections.\n  NOTE: Clients that do not support encryption cannot access SMB shares.\n* `DEFAULT`: Use the TrueNAS SMB server default encryption settings. Currently, this is the same as `NEGOTIATE`.",
		},
		{
			name: "bindip", api: "bindip", path: "bindip",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "List of IP addresses used by the TrueNAS SMB server. ",
			elem: &node{
				name: "", api: "", path: "bindip",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "server_sid", api: "server_sid", path: "server_sid",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "The unique identifier for the TrueNAS SMB server. It also serves as the domain SID for all local SMB user and     group accounts. ",
		},
		{
			name: "smb_options", api: "smb_options", path: "smb_options",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional unvalidated and unsupported configuration options for the SMB server.\nWARNING: Using `smb_options` may produce unexpected server behavior. ",
		},
		{
			name: "debug", api: "debug", path: "debug",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Set SMB log levels to debug. Use this setting only when troubleshooting a specific SMB issue. Do not use it     in production environments. ",
		},
	},
}

// NewSMBConfig returns the smb_config resource.
func NewSMBConfig() resource.Resource {
	return &sMBConfigResource{crudResource{model: modelSMBConfig}}
}

type sMBConfigResource struct{ crudResource }

// NewSMBConfigDataSource returns the smb_config data source.
func NewSMBConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelSMBConfig, attrs: dataAttrs(modelSMBConfig)}
}

func (r *sMBConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update SMB Service Configuration.\n\n`netbiosname` defaults to the original hostname of the system.\n\n`netbiosalias` a list of netbios aliases. If Server is joined to an AD domain, additional Kerberos\nService Principal Names will be generated for these aliases.\n\n`workgroup` specifies the NetBIOS workgroup to which the TrueNAS server belongs. This will be\nautomatically set to the correct value during the process of joining an AD domain.\nNOTE: `workgroup` and `netbiosname` should have different values.\n\n`enable_smb1` allows legacy SMB clients to connect to the server when enabled.\n\n`aapl_extensions` enables support for SMB2 protocol extensions for MacOS clients. This is not a\nrequirement for MacOS support, but is currently a requirement for time machine support.\n\n`localmaster` when set, determines if the system participates in a browser election.\n\n`guest` attribute is specified to select the account to be used for guest access. It defaults to \"nobody\".\n\nThe group specified as the SMB `admin_group` will be automatically added as a foreign group member\nof S-1-5-32-544 (builtin\\admins). This will afford the group all privileges granted to a local admin.\nAny SMB group may be selected (including AD groups).\n\n`ntlmv1_auth` enables a legacy and insecure authentication method, which may be required for legacy or\npoorly-implemented SMB clients.\n\n`encryption` set global server behavior with regard to SMB encrpytion. Options are DEFAULT (which\nfollows the upstream defaults -- currently identical to NEGOTIATE), NEGOTIATE encrypts SMB transport\nonly if requested by the SMB client, DESIRED encrypts SMB transport if supported by the SMB client,\nREQUIRED only allows encrypted transport to the SMB server. Mandatory SMB encryption is not\ncompatible with SMB1 server support in TrueNAS.\n\n`smb_options` smb.conf parameters that are not covered by the above supported configuration options may be\nadded as an smb_option. Not all options are tested or supported, and behavior of smb_options may change\nbetween releases. Stability of smb.conf options is not guaranteed.\n\nThis resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"netbiosname": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The NetBIOS name of this server. ",
			},
			"netbiosalias": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Alternative netbios names of the TrueNAS server. These names are announced through NetBIOS name server and     registered in Active Directory when TrueNAS joins the domain.",
				ElementType:         types.StringType,
			},
			"workgroup": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Workgroup name. When TrueNAS joins active directory, it automatically changes this value to match the NetBIOS     domain of the Active Directory domain. ",
			},
			"description": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Description of the SMB server. SMB clients may see this description during some operations. ",
			},
			"enable_smb1": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable SMB1 support on the server. WARNING: using the SMB1 protocol is not recommended. ",
			},
			"unixcharset": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Select character set for file names on local filesystem. Use this option only if you know the names are not     UTF-8. ",
				Validators:          []validator.String{stringvalidator.OneOf("UTF-8", "GB2312", "HZ-GB-2312", "CP1361", "BIG5", "BIG5HKSCS", "CP037", "CP273", "CP424", "CP437", "CP500", "CP775", "CP850", "CP852", "CP855", "CP857", "CP858", "CP860", "CP861", "CP862", "CP863", "CP864", "CP865", "CP866", "CP869", "CP932", "CP949", "CP950", "CP1026", "CP1125", "CP1140", "CP1250", "CP1251", "CP1252", "CP1253", "CP1254", "CP1255", "CP1256", "CP1257", "CP1258", "EUC_JIS_2004", "EUC_JISX0213", "EUC_JP", "EUC_KR", "GB18030", "GBK", "HZ", "ISO2022_JP", "ISO2022_JP_1", "ISO2022_JP_2", "ISO2022_JP_2004", "ISO2022_JP_3", "ISO2022_JP_EXT", "ISO2022_KR", "ISO8859_1", "ISO8859_2", "ISO8859_3", "ISO8859_4", "ISO8859_5", "ISO8859_6", "ISO8859_7", "ISO8859_8", "ISO8859_9", "ISO8859_10", "ISO8859_11", "ISO8859_13", "ISO8859_14", "ISO8859_15", "ISO8859_16", "JOHAB", "KOI8_R", "KZ1048", "LATIN_1", "MAC_CYRILLIC", "MAC_GREEK", "MAC_ICELAND", "MAC_LATIN2", "MAC_ROMAN", "MAC_TURKISH", "PTCP154", "SHIFT_JIS", "SHIFT_JIS_2004", "SHIFT_JISX0213", "TIS_620", "UTF_16", "UTF_16_BE", "UTF_16_LE")},
			},
			"localmaster": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "When set to `true` the NetBIOS name server in TrueNAS participates in elections for the local master browser.\nWhen set to `false` the NetBIOS name server does not attempt to become a local master browser on a subnet and     loses all browsing elections.\n\nNOTE: This parameter has no effect if the NetBIOS name server is disabled. ",
			},
			"syslog": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Send log messages to syslog. Enable this option if you want SMB server error logs to be included in     information sent to a remote syslog server. NOTE: This requires that remote syslog is globally configured on     TrueNAS. ",
			},
			"aapl_extensions": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable support for SMB2/3 AAPL protocol extensions. This setting makes the TrueNAS server advertise support     for Apple protocol extensions as a MacOS server. Enabling this is required for Time Machine support. ",
			},
			"admin_group": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The selected group has full administrator privileges on TrueNAS via the SMB protocol. ",
			},
			"guest": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "SMB guest account username. This username provides access to legacy SMB shares with guest access enabled.     It must be a valid, existing local user account. ",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"filemask": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "`smb.conf` create mask. DEFAULT applies current server default which is 664. ",
			},
			"dirmask": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "`smb.conf` directory mask. DEFAULT applies current server default which is 775. ",
			},
			"ntlmv1_auth": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable legacy and very insecure NTLMv1 authentication. This should never be done except     in extreme edge cases and may be against regulations in non-home environments. ",
			},
			"multichannel": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Enable SMB3 multi-channel support. ",
			},
			"encryption": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "SMB2/3 transport encryption setting for the TrueNAS SMB server.\n\n* `NEGOTIATE`: Enable negotiation of data encryption. Encrypt data only if the client explicitly requests it.\n* `DESIRED`: Enable negotiation of data encryption. Encrypt data on sessions and share connections for clients       that support it.\n* `REQUIRED`: Require data encryption for sessions and share connections.\n  NOTE: Clients that do not support encryption cannot access SMB shares.\n* `DEFAULT`: Use the TrueNAS SMB server default encryption settings. Currently, this is the same as `NEGOTIATE`.",
				Validators:          []validator.String{stringvalidator.OneOf("DEFAULT", "NEGOTIATE", "DESIRED", "REQUIRED")},
			},
			"bindip": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "List of IP addresses used by the TrueNAS SMB server. ",
				ElementType:         types.StringType,
			},
			"server_sid": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The unique identifier for the TrueNAS SMB server. It also serves as the domain SID for all local SMB user and     group accounts. ",
			},
			"smb_options": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional unvalidated and unsupported configuration options for the SMB server.\nWARNING: Using `smb_options` may produce unexpected server behavior. ",
			},
			"debug": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Set SMB log levels to debug. Use this setting only when troubleshooting a specific SMB issue. Do not use it     in production environments. ",
			},
		},
	}
}
