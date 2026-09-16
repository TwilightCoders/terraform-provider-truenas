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
)

var modelMailConfig = &model{
	typeName:     "mail_config",
	namespace:    "mail",
	primaryKey:   "id",
	idKind:       kindInt,
	singleton:    true,
	updateMethod: "mail.update",
	getMethod:    "mail.config",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "fromemail", api: "fromemail", path: "fromemail",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "The sending address that the mail server will use for sending emails.",
		},
		{
			name: "fromname", api: "fromname", path: "fromname",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Display name that will appear as the sender name in outgoing emails.",
		},
		{
			name: "outgoingserver", api: "outgoingserver", path: "outgoingserver",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Hostname or IP address of the SMTP server used for sending emails.",
		},
		{
			name: "port", api: "port", path: "port",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "TCP port number for the SMTP server connection.",
		},
		{
			name: "security", api: "security", path: "security",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Type of encryption.",
		},
		{
			name: "smtp", api: "smtp", path: "smtp",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			description: "Whether SMTP authentication is enabled and `user`, `pass` are required.",
		},
		{
			name: "user", api: "user", path: "user",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "SMTP username.",
		},
		{
			name: "pass", api: "pass", path: "pass",
			kind: kindString, role: roleOptional,
			nullable: true, sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "SMTP password.",
		},
		{
			name: "oauth", api: "oauth", path: "oauth",
			kind: kindAny, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "OAuth configuration for email providers that support it or `null` for basic authentication.",
		},
		{
			name: "pass_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `pass` again. Terraform never stores `pass`.",
		},
	},
}

// NewMailConfig returns the mail_config resource.
func NewMailConfig() resource.Resource {
	return &mailConfigResource{crudResource{model: modelMailConfig}}
}

type mailConfigResource struct{ crudResource }

// NewMailConfigDataSource returns the mail_config data source.
func NewMailConfigDataSource() datasource.DataSource {
	return &dataSource{model: modelMailConfig, attrs: dataAttrs(modelMailConfig)}
}

func (r *mailConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Update Mail Service Configuration.

This resource manages existing settings: creating it applies the configured attributes, and destroying it only removes it from Terraform state.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"fromemail": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "The sending address that the mail server will use for sending emails.",
			},
			"fromname": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Display name that will appear as the sender name in outgoing emails.",
			},
			"outgoingserver": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Hostname or IP address of the SMTP server used for sending emails.",
			},
			"port": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "TCP port number for the SMTP server connection.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"security": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Type of encryption.",
				Validators:          []validator.String{stringvalidator.OneOf("PLAIN", "SSL", "TLS")},
			},
			"smtp": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether SMTP authentication is enabled and `user`, `pass` are required.",
			},
			"user": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "SMTP username.",
			},
			"pass": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "SMTP password.",
			},
			"oauth": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "OAuth configuration for email providers that support it or `null` for basic authentication.",
			},
			"pass_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `pass` again. Terraform never stores `pass`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
