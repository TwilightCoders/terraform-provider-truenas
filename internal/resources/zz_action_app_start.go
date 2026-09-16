// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionAppStart = &methodAction{
	typeName: "app_start",
	method:   "app.start",
	args: []*node{
		{
			name: "name", api: "app_name", path: "app_name",
			kind: kindString, role: roleRequired,
			description: "Name of the application to start.",
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Starts an app and waits for it to finish.",
		Attributes: map[string]actionschema.Attribute{
			"name": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the application to start.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	},
}

// NewAppStartAction returns the app_start action.
func NewAppStartAction() action.Action {
	a := *actionAppStart
	return &a
}
