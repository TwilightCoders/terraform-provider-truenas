// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionUiRestart = &methodAction{
	typeName: "ui_restart",
	method:   "system.general.ui_restart",
	args: []*node{
		{
			name: "delay", api: "delay", path: "delay",
			kind: kindInt, role: roleOptional,
			description: "How long to wait before the UI is restarted. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Restarts the web UI after `delay` seconds, applying UI settings. Aborts every HTTP connection, including the provider's.",
		Attributes: map[string]actionschema.Attribute{
			"delay": actionschema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "How long to wait before the UI is restarted. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 0}},
			},
		},
	},
}

// NewUiRestartAction returns the ui_restart action.
func NewUiRestartAction() action.Action {
	a := *actionUiRestart
	return &a
}
