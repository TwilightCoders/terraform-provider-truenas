// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionCloudsyncSync = &methodAction{
	typeName: "cloudsync_sync",
	method:   "cloudsync.sync",
	args: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleRequired,
			description: "ID of the cloud sync task to run.",
		},
		{
			name: "options", api: "cloud_sync_sync_options", path: "cloud_sync_sync_options",
			kind: kindObject, role: roleOptional,
			description: "Options for the sync operation. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
			children: []*node{
				{
					name: "dry_run", api: "dry_run", path: "cloud_sync_sync_options.dry_run",
					kind: kindBool, role: roleOptional,
					description: "Whether to perform a dry run without making actual changes. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				},
			},
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Runs a cloud sync task now and waits for it to finish.",
		Attributes: map[string]actionschema.Attribute{
			"id": actionschema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "ID of the cloud sync task to run.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"options": actionschema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Options for the sync operation. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				Attributes: map[string]actionschema.Attribute{
					"dry_run": actionschema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether to perform a dry run without making actual changes. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
					},
				},
			},
		},
	},
}

// NewCloudsyncSyncAction returns the cloudsync_sync action.
func NewCloudsyncSyncAction() action.Action {
	a := *actionCloudsyncSync
	return &a
}
