// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var actionServiceControl = &methodAction{
	typeName: "service_control",
	method:   "service.control",
	args: []*node{
		{
			name: "verb", api: "verb", path: "verb",
			kind: kindString, role: roleRequired,
			description: "The service operation to perform.",
		},
		{
			name: "service", api: "service", path: "service",
			kind: kindString, role: roleRequired,
			description: "Name of the service to control.",
		},
		{
			name: "options", api: "options", path: "options",
			kind: kindObject, role: roleOptional,
			description: "Options for controlling the service operation behavior. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
			children: []*node{
				{
					name: "ha_propagate", api: "ha_propagate", path: "options.ha_propagate",
					kind: kindBool, role: roleOptional,
					description: "Whether to propagate the service operation to the HA peer in a high-availability setup. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				},
				{
					name: "silent", api: "silent", path: "options.silent",
					kind: kindBool, role: roleOptional,
					description: "Return `false` instead of an error if the operation fails. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				},
				{
					name: "timeout", api: "timeout", path: "options.timeout",
					kind: kindInt, role: roleOptional,
					nullable:    true,
					description: "Maximum time in seconds to wait for the service operation to complete. `null` for no timeout. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				},
			},
		},
	},
	schema: actionschema.Schema{
		MarkdownDescription: "Starts, stops, restarts or reloads a system service and waits for it to finish.",
		Attributes: map[string]actionschema.Attribute{
			"verb": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The service operation to perform.",
				Validators:          []validator.String{stringvalidator.OneOf("START", "STOP", "RESTART", "RELOAD")},
			},
			"service": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the service to control.",
			},
			"options": actionschema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Options for controlling the service operation behavior. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
				Attributes: map[string]actionschema.Attribute{
					"ha_propagate": actionschema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether to propagate the service operation to the HA peer in a high-availability setup. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
					},
					"silent": actionschema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Return `false` instead of an error if the operation fails. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
					},
					"timeout": actionschema.NumberAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum time in seconds to wait for the service operation to complete. `null` for no timeout. TrueNAS does not report this value, so changes made outside Terraform are not detected.",
						Validators:          []validator.Number{numberIsInteger{}},
					},
				},
			},
		},
	},
}

// NewServiceControlAction returns the service_control action.
func NewServiceControlAction() action.Action {
	a := *actionServiceControl
	return &a
}
