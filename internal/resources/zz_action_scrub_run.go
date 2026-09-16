// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionScrubRun = &methodAction{
	typeName: "scrub_run",
	method:   "pool.scrub.run",
	args: []*node{
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			description: "Name of the pool to run scrub on.",
		},
		{
			name: "threshold", api: "threshold", path: "threshold",
			kind: kindInt, role: roleOptional,
			description: "Days before a scrub is due when the scrub should start. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Starts a scrub of a pool when the last one is older than `threshold` days.",
		Attributes: map[string]actionschema.Attribute{
			"name": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the pool to run scrub on.",
			},
			"threshold": actionschema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Days before a scrub is due when the scrub should start. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	},
}

// NewScrubRunAction returns the scrub_run action.
func NewScrubRunAction() action.Action {
	a := *actionScrubRun
	return &a
}
