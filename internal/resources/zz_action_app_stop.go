// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionAppStop = &methodAction{
	typeName: "app_stop",
	method:   "app.stop",
	args: []*node{
		{
			name: "name", api: "app_name", path: "app_name",
			kind: kindString, role: roleRequired,
			description: "Name of the application to stop.",
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Stops an app and waits for it to finish.",
		Attributes: map[string]actionschema.Attribute{
			"name": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the application to stop.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	},
}

// NewAppStopAction returns the app_stop action.
func NewAppStopAction() action.Action {
	a := *actionAppStop
	return &a
}
