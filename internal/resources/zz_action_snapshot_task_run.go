// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionSnapshotTaskRun = &methodAction{
	typeName: "snapshot_task_run",
	method:   "pool.snapshottask.run",
	args: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleRequired,
			description: "ID of the periodic snapshot task to run immediately.",
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Takes the snapshot for a periodic snapshot task now.",
		Attributes: map[string]actionschema.Attribute{
			"id": actionschema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "ID of the periodic snapshot task to run immediately.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	},
}

// NewSnapshotTaskRunAction returns the snapshot_task_run action.
func NewSnapshotTaskRunAction() action.Action {
	a := *actionSnapshotTaskRun
	return &a
}
