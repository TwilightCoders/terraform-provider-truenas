// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionReplicationRun = &methodAction{
	typeName: "replication_run",
	method:   "replication.run",
	args: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleRequired,
			description: "ID of the replication task to run.",
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Runs a replication task now and waits for it to finish.",
		Attributes: map[string]actionschema.Attribute{
			"id": actionschema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "ID of the replication task to run.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	},
}

// NewReplicationRunAction returns the replication_run action.
func NewReplicationRunAction() action.Action {
	a := *actionReplicationRun
	return &a
}
