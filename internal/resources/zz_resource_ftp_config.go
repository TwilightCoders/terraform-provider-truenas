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

var modelFTPConfig = &model{
	typeName:     "ftp_config",
	namespace:    "ftp",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	updateMethod: "ftp.update",
	getMethod:    "ftp.config",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "port", api: "port", path: "port",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "TCP port number on which the FTP service listens for incoming connections.",
		},
		{
			name: "clients", api: "clients", path: "clients",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum number of simultaneous client connections allowed.",
		},
		{
			name: "ipconnections", api: "ipconnections", path: "ipconnections",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum number of connections allowed from a single IP address. 0 means unlimited.",
		},
		{
			name: "loginattempt", api: "loginattempt", path: "loginattempt",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum number of failed login attempts before blocking an IP address. 0 disables this limit.",
		},
		{
			name: "timeout", api: "timeout", path: "timeout",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Idle timeout in seconds before disconnecting inactive clients. 0 disables timeout.",
		},
		{
			name: "timeout_notransfer", api: "timeout_notransfer", path: "timeout_notransfer",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Timeout in seconds for clients that connect but do not transfer data. 0 disables timeout.",
		},
		{
			name: "onlyanonymous", api: "onlyanonymous", path: "onlyanonymous",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to allow only anonymous FTP access, disabling authenticated user login.",
		},
		{
			name: "anonpath", api: "anonpath", path: "anonpath",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Filesystem path for anonymous FTP users. `null` to use the default anonymous FTP directory.",
		},
		{
			name: "onlylocal", api: "onlylocal", path: "onlylocal",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to allow only local system users to login, disabling anonymous access.",
		},
		{
			name: "banner", api: "banner", path: "banner",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Welcome message displayed to FTP clients upon connection.",
		},
		{
			name: "filemask", api: "filemask", path: "filemask",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Default Unix permissions (umask) for files created by FTP users.",
		},
		{
			name: "dirmask", api: "dirmask", path: "dirmask",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Default Unix permissions (umask) for directories created by FTP users.",
		},
		{
			name: "fxp", api: "fxp", path: "fxp",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to enable File eXchange Protocol (FXP) for server-to-server transfers.",
		},
		{
			name: "resume", api: "resume", path: "resume",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to allow clients to resume interrupted file transfers.",
		},
		{
			name: "defaultroot", api: "defaultroot", path: "defaultroot",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to restrict users to their home directories (chroot jail).",
		},
		{
			name: "ident", api: "ident", path: "ident",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to perform RFC 1413 ident lookups on connecting clients.",
		},
		{
			name: "reversedns", api: "reversedns", path: "reversedns",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to perform reverse DNS lookups on client IP addresses for logging.",
		},
		{
			name: "masqaddress", api: "masqaddress", path: "masqaddress",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Public IP address to advertise to clients for passive mode connections when behind NAT.",
		},
		{
			name: "passiveportsmin", api: "passiveportsmin", path: "passiveportsmin",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Minimum port number for passive mode data connections. Must be 0 or between 1024-65535.",
		},
		{
			name: "passiveportsmax", api: "passiveportsmax", path: "passiveportsmax",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum port number for passive mode data connections. Must be 0 or between 1024-65535.",
		},
		{
			name: "localuserbw", api: "localuserbw", path: "localuserbw",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum upload bandwidth in KiB/s for local users. 0 means unlimited.",
		},
		{
			name: "localuserdlbw", api: "localuserdlbw", path: "localuserdlbw",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum download bandwidth in KiB/s for local users. 0 means unlimited.",
		},
		{
			name: "anonuserbw", api: "anonuserbw", path: "anonuserbw",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum upload bandwidth in KiB/s for anonymous users. 0 means unlimited.",
		},
		{
			name: "anonuserdlbw", api: "anonuserdlbw", path: "anonuserdlbw",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Maximum download bandwidth in KiB/s for anonymous users. 0 means unlimited.",
		},
		{
			name: "tls", api: "tls", path: "tls",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to enable TLS/SSL encryption for FTP connections.",
		},
		{
			name: "tls_policy", api: "tls_policy", path: "tls_policy",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "TLS policy for connections. Values include: `\"on\"` (required), `\"off\"` (disabled), `\"data\"` (data only),     `\"auth\"` (authentication only), `\"ctrl\"` (control only), or combinations with `+` and `!` modifiers.",
		},
		{
			name: "tls_opt_allow_client_renegotiations", api: "tls_opt_allow_client_renegotiations", path: "tls_opt_allow_client_renegotiations",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to allow TLS clients to initiate renegotiation of the TLS connection.",
		},
		{
			name: "tls_opt_allow_dot_login", api: "tls_opt_allow_dot_login", path: "tls_opt_allow_dot_login",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to allow .ftpaccess files to override TLS requirements for specific users.",
		},
		{
			name: "tls_opt_allow_per_user", api: "tls_opt_allow_per_user", path: "tls_opt_allow_per_user",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to allow per-user TLS configuration overrides.",
		},
		{
			name: "tls_opt_common_name_required", api: "tls_opt_common_name_required", path: "tls_opt_common_name_required",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to require client certificates to have a Common Name field.",
		},
		{
			name: "tls_opt_enable_diags", api: "tls_opt_enable_diags", path: "tls_opt_enable_diags",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to enable detailed TLS diagnostic logging.",
		},
		{
			name: "tls_opt_export_cert_data", api: "tls_opt_export_cert_data", path: "tls_opt_export_cert_data",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to export client certificate data to environment variables.",
		},
		{
			name: "tls_opt_no_empty_fragments", api: "tls_opt_no_empty_fragments", path: "tls_opt_no_empty_fragments",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to disable empty TLS record fragments to improve compatibility with some clients.      Disabling increases vulnerability to some attack vectors.",
		},
		{
			name: "tls_opt_no_session_reuse_required", api: "tls_opt_no_session_reuse_required", path: "tls_opt_no_session_reuse_required",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to disable the requirement for TLS session reuse.",
		},
		{
			name: "tls_opt_stdenvvars", api: "tls_opt_stdenvvars", path: "tls_opt_stdenvvars",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to export standard TLS environment variables for use by external programs.",
		},
		{
			name: "tls_opt_dns_name_required", api: "tls_opt_dns_name_required", path: "tls_opt_dns_name_required",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to require client certificates to contain a DNS name in the Subject Alternative Name extension.     The `reversedns` setting must also be enabled.",
		},
		{
			name: "tls_opt_ip_address_required", api: "tls_opt_ip_address_required", path: "tls_opt_ip_address_required",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether to require client certificates to contain an IP address in the Subject Alternative Name extension.",
		},
		{
			name: "ssltls_certificate", api: "ssltls_certificate", path: "ssltls_certificate",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "ID of the certificate to use for TLS/SSL connections. `null` to use the default system certificate.",
		},
		{
			name: "options", api: "options", path: "options",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional ProFTPD configuration directives to include in the server configuration.     Manual directives may render the FTP service non-functional and should be used with caution.",
		},
		{
			name: "unset", api: "", path: "unset",
			kind: kindList, role: roleOptional,
			description: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
			elem: &node{
				name: "unset", api: "", path: "unset",
				kind: kindString, role: roleRequired,
			},
		},
	},
}

// NewFTPConfig returns the ftp_config resource.
func NewFTPConfig() resource.Resource {
	return &fTPConfigResource{crudResource{model: modelFTPConfig}}
}

type fTPConfigResource struct{ crudResource }

// NewFTPConfigDataSource returns the ftp_config data source.
func NewFTPConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelFTPConfig, attrs: dataAttrs(modelFTPConfig)}
}

func (r *fTPConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update ftp service configuration.\n\n`clients` is an integer value which sets the maximum number of simultaneous clients allowed. It defaults to 32.\n\n`ipconnections` is an integer value which shows the maximum number of connections per IP address. It defaults\nto 0 which equals to unlimited.\n\n`timeout` is the maximum number of seconds that proftpd will allow clients to stay connected without receiving\nany data on either the control or data connection.\n\n`timeout_notransfer` is the maximum number of seconds a client is allowed to spend connected, after\nauthentication, without issuing a command which results in creating an active or passive data connection\n(i.e. sending/receiving a file, or receiving a directory listing).\n\n`onlyanonymous` allows anonymous FTP logins with access to the directory specified by `anonpath`.\n\n`banner` is a message displayed to local login users after they successfully authenticate. It is not displayed\nto anonymous login users.\n\n`filemask` sets the default permissions for newly created files which by default are 077.\n\n`dirmask` sets the default permissions for newly created directories which by default are 077.\n\n`resume` if set allows FTP clients to resume interrupted transfers.\n\n`fxp` if set to true indicates that File eXchange Protocol is enabled. Generally it is discouraged as it\nmakes the server vulnerable to FTP bounce attacks.\n\n`defaultroot` when set ensures that for local users, home directory access is only granted if the user\nis a member of group wheel.\n\n`ident` is a boolean value which when set to true indicates that IDENT authentication is required. If identd\nis not running on the client, this can result in timeouts.\n\n`masqaddress` is the public IP address or hostname which is set if FTP clients cannot connect through a\nNAT device.\n\n`localuserbw` is a positive integer value which indicates maximum upload bandwidth in KB/s for local user.\nDefault of zero indicates unlimited upload bandwidth ( from the FTP server configuration ).\n\n`localuserdlbw` is a positive integer value which indicates maximum download bandwidth in KB/s for local user.\nDefault of zero indicates unlimited download bandwidth ( from the FTP server configuration ).\n\n`anonuserbw` is a positive integer value which indicates maximum upload bandwidth in KB/s for anonymous user.\nDefault of zero indicates unlimited upload bandwidth ( from the FTP server configuration ).\n\n`anonuserdlbw` is a positive integer value which indicates maximum download bandwidth in KB/s for anonymous\nuser. Default of zero indicates unlimited download bandwidth ( from the FTP server configuration ).\n\n`tls` is a boolean value which when set indicates that encrypted connections are enabled. This requires a\ncertificate to be configured first with the certificate service and the id of certificate is passed on in\n`ssltls_certificate`.\n\n`tls_policy` defines whether the control channel, data channel, both channels, or neither channel of an FTP\nsession must occur over SSL/TLS.\n\n`tls_opt_enable_diags` is a boolean value when set, logs verbosely. This is helpful when troubleshooting a\nconnection.\n\n`options` is a string used to add proftpd(8) parameters not covered by ftp service.\n\nThis resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"port": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "TCP port number on which the FTP service listens for incoming connections.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"clients": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum number of simultaneous client connections allowed.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 10000}},
			},
			"ipconnections": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum number of connections allowed from a single IP address. 0 means unlimited.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 1000}},
			},
			"loginattempt": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum number of failed login attempts before blocking an IP address. 0 disables this limit.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 1000}},
			},
			"timeout": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Idle timeout in seconds before disconnecting inactive clients. 0 disables timeout.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 10000}},
			},
			"timeout_notransfer": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Timeout in seconds for clients that connect but do not transfer data. 0 disables timeout.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}, numberAtMost{max: 10000}},
			},
			"onlyanonymous": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to allow only anonymous FTP access, disabling authenticated user login.",
			},
			"anonpath": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Filesystem path for anonymous FTP users. `null` to use the default anonymous FTP directory.",
			},
			"onlylocal": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to allow only local system users to login, disabling anonymous access.",
			},
			"banner": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Welcome message displayed to FTP clients upon connection.",
			},
			"filemask": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Default Unix permissions (umask) for files created by FTP users.",
			},
			"dirmask": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Default Unix permissions (umask) for directories created by FTP users.",
			},
			"fxp": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable File eXchange Protocol (FXP) for server-to-server transfers.",
			},
			"resume": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to allow clients to resume interrupted file transfers.",
			},
			"defaultroot": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to restrict users to their home directories (chroot jail).",
			},
			"ident": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to perform RFC 1413 ident lookups on connecting clients.",
			},
			"reversedns": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to perform reverse DNS lookups on client IP addresses for logging.",
			},
			"masqaddress": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Public IP address to advertise to clients for passive mode connections when behind NAT.",
			},
			"passiveportsmin": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Minimum port number for passive mode data connections. Must be 0 or between 1024-65535.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"passiveportsmax": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum port number for passive mode data connections. Must be 0 or between 1024-65535.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"localuserbw": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum upload bandwidth in KiB/s for local users. 0 means unlimited.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}},
			},
			"localuserdlbw": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum download bandwidth in KiB/s for local users. 0 means unlimited.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}},
			},
			"anonuserbw": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum upload bandwidth in KiB/s for anonymous users. 0 means unlimited.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}},
			},
			"anonuserdlbw": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum download bandwidth in KiB/s for anonymous users. 0 means unlimited.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}},
			},
			"tls": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable TLS/SSL encryption for FTP connections.",
			},
			"tls_policy": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "TLS policy for connections. Values include: `\"on\"` (required), `\"off\"` (disabled), `\"data\"` (data only),     `\"auth\"` (authentication only), `\"ctrl\"` (control only), or combinations with `+` and `!` modifiers.",
				Validators:          []validator.String{stringvalidator.OneOf("", "on", "off", "data", "!data", "auth", "ctrl", "ctrl+data", "ctrl+!data", "auth+data", "auth+!data")},
			},
			"tls_opt_allow_client_renegotiations": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to allow TLS clients to initiate renegotiation of the TLS connection.",
			},
			"tls_opt_allow_dot_login": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to allow .ftpaccess files to override TLS requirements for specific users.",
			},
			"tls_opt_allow_per_user": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to allow per-user TLS configuration overrides.",
			},
			"tls_opt_common_name_required": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to require client certificates to have a Common Name field.",
			},
			"tls_opt_enable_diags": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable detailed TLS diagnostic logging.",
			},
			"tls_opt_export_cert_data": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to export client certificate data to environment variables.",
			},
			"tls_opt_no_empty_fragments": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to disable empty TLS record fragments to improve compatibility with some clients.      Disabling increases vulnerability to some attack vectors.",
			},
			"tls_opt_no_session_reuse_required": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to disable the requirement for TLS session reuse.",
			},
			"tls_opt_stdenvvars": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to export standard TLS environment variables for use by external programs.",
			},
			"tls_opt_dns_name_required": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to require client certificates to contain a DNS name in the Subject Alternative Name extension.     The `reversedns` setting must also be enabled.",
			},
			"tls_opt_ip_address_required": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to require client certificates to contain an IP address in the Subject Alternative Name extension.",
			},
			"ssltls_certificate": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "ID of the certificate to use for TLS/SSL connections. `null` to use the default system certificate.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"options": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional ProFTPD configuration directives to include in the server configuration.     Manual directives may render the FTP service non-functional and should be used with caution.",
			},
			"unset": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
				ElementType:         types.StringType,
			},
		},
	}
}
