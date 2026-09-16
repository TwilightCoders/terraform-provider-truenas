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

var modelISCSIAuth = &model{
	typeName:     "iscsi_auth",
	namespace:    "iscsi.auth",
	primaryKey:   "id",
	idKind:       kindInt,
	getMethod:    "iscsi.auth.get_instance",
	deleteMethod: "iscsi.auth.delete",
	createMethod: "iscsi.auth.create",
	updateMethod: "iscsi.auth.update",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "tag", api: "tag", path: "tag",
			kind: kindInt, role: roleRequired,
			readable: true, updatable: true,
			description: "Numeric tag used to associate this credential with iSCSI targets.",
		},
		{
			name: "user", api: "user", path: "user",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Username for iSCSI CHAP authentication.",
		},
		{
			name: "secret", api: "secret", path: "secret",
			kind: kindString, role: roleRequired,
			sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Password/secret for iSCSI CHAP authentication.",
		},
		{
			name: "peeruser", api: "peeruser", path: "peeruser",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Username for mutual CHAP authentication or empty string if not configured. Defaults to `\"\"`.",
		},
		{
			name: "peersecret", api: "peersecret", path: "peersecret",
			kind: kindString, role: roleOptional,
			sensitive: true, writeOnly: true, readable: true, updatable: true,
			description: "Password/secret for mutual CHAP authentication or empty string if not configured.",
		},
		{
			name: "discovery_auth", api: "discovery_auth", path: "discovery_auth",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "NONE",
			description: "Authentication method for target discovery. If \"CHAP_MUTUAL\" is selected for target discovery, it is only     permitted for a single entry systemwide. Defaults to `\"NONE\"`.",
		},
		{
			name: "secret_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `secret` again. Terraform never stores `secret`.",
		},
		{
			name: "peersecret_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `peersecret` again. Terraform never stores `peersecret`.",
		},
	},
}

// NewISCSIAuth returns the iscsi_auth resource.
func NewISCSIAuth() resource.Resource {
	return &iSCSIAuthResource{crudResource{model: modelISCSIAuth}}
}

type iSCSIAuthResource struct{ crudResource }

// NewISCSIAuthDataSource returns the iscsi_auth data source.
func NewISCSIAuthDataSource() datasource.DataSource {
	return &dataSource{model: modelISCSIAuth, attrs: dataAttrs(modelISCSIAuth)}
}

// NewISCSIAuthList returns the iscsi_auth list resource, for terraform query.
func NewISCSIAuthList() list.ListResource {
	return &listResource{crudResource{model: modelISCSIAuth}}
}

func (r *iSCSIAuthResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Authorized Access.\n\n`tag` should be unique among all configured iSCSI Authorized Accesses.\n\n`secret` and `peersecret` should have length between 12-16 letters inclusive.\n\n`peeruser` and `peersecret` are provided only when configuring mutual CHAP. `peersecret` should not be\nsimilar to `secret`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"tag": schema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "Numeric tag used to associate this credential with iSCSI targets.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Username for iSCSI CHAP authentication.",
			},
			"secret": schema.StringAttribute{
				Required: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Password/secret for iSCSI CHAP authentication.",
			},
			"peeruser": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Username for mutual CHAP authentication or empty string if not configured. Defaults to `\"\"`.",
				Default:             stringdefault.StaticString(""),
			},
			"peersecret": schema.StringAttribute{
				Optional: true, Sensitive: true, WriteOnly: true,
				MarkdownDescription: "Password/secret for mutual CHAP authentication or empty string if not configured.",
			},
			"discovery_auth": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Authentication method for target discovery. If \"CHAP_MUTUAL\" is selected for target discovery, it is only     permitted for a single entry systemwide. Defaults to `\"NONE\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("NONE", "CHAP", "CHAP_MUTUAL")},
				Default:             stringdefault.StaticString("NONE"),
			},
			"secret_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `secret` again. Terraform never stores `secret`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"peersecret_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `peersecret` again. Terraform never stores `peersecret`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
