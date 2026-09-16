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

var modelSSHConfig = &model{
	typeName:     "ssh_config",
	namespace:    "ssh",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	getMethod:    "ssh.config",
	updateMethod: "ssh.update",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "bindiface", api: "bindiface", path: "bindiface",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Array of network interface names to bind the SSH service to.",
			elem: &node{
				name: "", api: "", path: "bindiface",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "tcpport", api: "tcpport", path: "tcpport",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "TCP port number for SSH connections.",
		},
		{
			name: "password_login_groups", api: "password_login_groups", path: "password_login_groups",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Array of group names allowed to authenticate with passwords.",
			elem: &node{
				name: "", api: "", path: "password_login_groups",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "passwordauth", api: "passwordauth", path: "passwordauth",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether password authentication is enabled.",
		},
		{
			name: "kerberosauth", api: "kerberosauth", path: "kerberosauth",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether Kerberos authentication is enabled.",
		},
		{
			name: "tcpfwd", api: "tcpfwd", path: "tcpfwd",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether TCP forwarding is enabled.",
		},
		{
			name: "compression", api: "compression", path: "compression",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether compression is enabled for SSH connections.",
		},
		{
			name: "sftp_log_level", api: "sftp_log_level", path: "sftp_log_level",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Logging level for SFTP subsystem (empty string means default).",
		},
		{
			name: "sftp_log_facility", api: "sftp_log_facility", path: "sftp_log_facility",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Syslog facility for SFTP logging (empty string means default).",
		},
		{
			name: "weak_ciphers", api: "weak_ciphers", path: "weak_ciphers",
			kind: kindList, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Array of weak ciphers to enable for compatibility with legacy clients.",
			elem: &node{
				name: "", api: "", path: "weak_ciphers",
				kind: kindString, role: roleRequired,
				readable: true,
			},
		},
		{
			name: "options", api: "options", path: "options",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Additional SSH daemon configuration options.",
		},
		{
			name: "host_dsa_key_pub", api: "host_dsa_key_pub", path: "host_dsa_key_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "DSA host public key. `null` if not configured.",
		},
		{
			name: "host_dsa_key_cert_pub", api: "host_dsa_key_cert_pub", path: "host_dsa_key_cert_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "DSA host certificate public key. `null` if not configured.",
		},
		{
			name: "host_ecdsa_key_pub", api: "host_ecdsa_key_pub", path: "host_ecdsa_key_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "ECDSA host public key. `null` if not configured.",
		},
		{
			name: "host_ecdsa_key_cert_pub", api: "host_ecdsa_key_cert_pub", path: "host_ecdsa_key_cert_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "ECDSA host certificate public key. `null` if not configured.",
		},
		{
			name: "host_ed25519_key_pub", api: "host_ed25519_key_pub", path: "host_ed25519_key_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Ed25519 host public key. `null` if not configured.",
		},
		{
			name: "host_ed25519_key_cert_pub", api: "host_ed25519_key_cert_pub", path: "host_ed25519_key_cert_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Ed25519 host certificate public key. `null` if not configured.",
		},
		{
			name: "host_key_pub", api: "host_key_pub", path: "host_key_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "Legacy SSH host public key. `null` if not configured.",
		},
		{
			name: "host_rsa_key_pub", api: "host_rsa_key_pub", path: "host_rsa_key_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "RSA host public key. `null` if not configured.",
		},
		{
			name: "host_rsa_key_cert_pub", api: "host_rsa_key_cert_pub", path: "host_rsa_key_cert_pub",
			kind: kindString, role: roleComputed,
			nullable: true, readable: true,
			description: "RSA host certificate public key. `null` if not configured.",
		},
	},
}

// NewSSHConfig returns the ssh_config resource.
func NewSSHConfig() resource.Resource {
	return &sSHConfigResource{crudResource{model: modelSSHConfig}}
}

type sSHConfigResource struct{ crudResource }

// NewSSHConfigDataSource returns the ssh_config data source.
func NewSSHConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelSSHConfig, attrs: dataAttrs(modelSSHConfig)}
}

func (r *sSHConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update settings of SSH daemon service.\n\nIf `bindiface` is empty it will listen for all available addresses.\n\n\n\nThis resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"bindiface": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of network interface names to bind the SSH service to.",
				ElementType:         types.StringType,
			},
			"tcpport": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "TCP port number for SSH connections.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 65535}},
			},
			"password_login_groups": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of group names allowed to authenticate with passwords.",
				ElementType:         types.StringType,
			},
			"passwordauth": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether password authentication is enabled.",
			},
			"kerberosauth": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether Kerberos authentication is enabled.",
			},
			"tcpfwd": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether TCP forwarding is enabled.",
			},
			"compression": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether compression is enabled for SSH connections.",
			},
			"sftp_log_level": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Logging level for SFTP subsystem (empty string means default).",
				Validators:          []validator.String{stringvalidator.OneOf("", "QUIET", "FATAL", "ERROR", "INFO", "VERBOSE", "DEBUG", "DEBUG2", "DEBUG3")},
			},
			"sftp_log_facility": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Syslog facility for SFTP logging (empty string means default).",
				Validators:          []validator.String{stringvalidator.OneOf("", "DAEMON", "USER", "AUTH", "LOCAL0", "LOCAL1", "LOCAL2", "LOCAL3", "LOCAL4", "LOCAL5", "LOCAL6", "LOCAL7")},
			},
			"weak_ciphers": schema.ListAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Array of weak ciphers to enable for compatibility with legacy clients.",
				ElementType:         types.StringType,
			},
			"options": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional SSH daemon configuration options.",
			},
			"host_dsa_key_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DSA host public key. `null` if not configured.",
			},
			"host_dsa_key_cert_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DSA host certificate public key. `null` if not configured.",
			},
			"host_ecdsa_key_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ECDSA host public key. `null` if not configured.",
			},
			"host_ecdsa_key_cert_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ECDSA host certificate public key. `null` if not configured.",
			},
			"host_ed25519_key_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Ed25519 host public key. `null` if not configured.",
			},
			"host_ed25519_key_cert_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Ed25519 host certificate public key. `null` if not configured.",
			},
			"host_key_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Legacy SSH host public key. `null` if not configured.",
			},
			"host_rsa_key_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RSA host public key. `null` if not configured.",
			},
			"host_rsa_key_cert_pub": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RSA host certificate public key. `null` if not configured.",
			},
		},
	}
}
