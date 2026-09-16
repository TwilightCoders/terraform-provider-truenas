// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionAppRedeploy = &methodAction{
	typeName: "app_redeploy",
	method:   "app.redeploy",
	args: []*node{
		{
			name: "name", api: "app_name", path: "app_name",
			kind: kindString, role: roleRequired,
			description: "Name of the application to redeploy (stop, pull latest images, and restart).",
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Redeploys an app, re-reading its compose configuration, and waits for it to finish.",
		Attributes: map[string]actionschema.Attribute{
			"name": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the application to redeploy (stop, pull latest images, and restart).",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	},
}

// NewAppRedeployAction returns the app_redeploy action.
func NewAppRedeployAction() action.Action {
	a := *actionAppRedeploy
	return &a
}
