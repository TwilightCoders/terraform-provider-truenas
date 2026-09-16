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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var modelTunable = &model{
	typeName:     "tunable",
	namespace:    "tunable",
	primaryKey:   "id",
	idKind:       kindInt,
	createMethod: "tunable.create",
	updateMethod: "tunable.update",
	getMethod:    "tunable.get_instance",
	deleteMethod: "tunable.delete",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "type", api: "type", path: "type",
			kind: kindString, role: roleOptionalComputed,
			replace: true, stable: true, readable: true,
			description: "* `SYSCTL`: `var` is a sysctl name (e.g. `kernel.watchdog`) and `value` is its corresponding value (e.g. `0`).\n* `UDEV`: `var` is a udev rules file name (e.g. `10-disable-usb`, `.rules` suffix will be appended automatically)     and `value` is its contents (e.g. `BUS==\"usb\", OPTIONS+=\"ignore_device\"`).\n* `ZFS`: `var` is a ZFS kernel module parameter name (e.g. `zfs_dirty_data_max_max`) and `value` is its value     (e.g. `783091712`). Changing this forces a new resource.",
		},
		{
			name: "var", api: "var", path: "var",
			kind: kindString, role: roleRequired,
			replace: true, readable: true,
			description: "Name or identifier of the system parameter to tune. Changing this forces a new resource.",
		},
		{
			name: "value", api: "value", path: "value",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Value to assign to the tunable parameter.",
		},
		{
			name: "comment", api: "comment", path: "comment",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Optional descriptive comment explaining the purpose of this tunable. Defaults to `\"\"`.",
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether this tunable is active and should be applied. Defaults to `true`.",
		},
		{
			name: "update_initramfs", api: "update_initramfs", path: "update_initramfs",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "If `false`, then initramfs will not be updated after creating a ZFS tunable and you will need to run     `system boot update_initramfs` manually. Defaults to `true`.",
		},
		{
			name: "orig_value", api: "orig_value", path: "orig_value",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Original system value of the parameter before this tunable was applied.",
		},
	},
}

// NewTunable returns the tunable resource.
func NewTunable() resource.Resource {
	return &tunableResource{crudResource{model: modelTunable}}
}

type tunableResource struct{ crudResource }

// NewTunableDataSource returns the tunable data source.
func NewTunableDataSource() datasource.DataSource {
	return &dataSource{model: modelTunable, attrs: dataAttrs(modelTunable)}
}

// NewTunableList returns the tunable list resource, for terraform query.
func NewTunableList() list.ListResource {
	return &listResource{crudResource{model: modelTunable}}
}

func (r *tunableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a tunable.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "* `SYSCTL`: `var` is a sysctl name (e.g. `kernel.watchdog`) and `value` is its corresponding value (e.g. `0`).\n* `UDEV`: `var` is a udev rules file name (e.g. `10-disable-usb`, `.rules` suffix will be appended automatically)     and `value` is its contents (e.g. `BUS==\"usb\", OPTIONS+=\"ignore_device\"`).\n* `ZFS`: `var` is a ZFS kernel module parameter name (e.g. `zfs_dirty_data_max_max`) and `value` is its value     (e.g. `783091712`). Changing this forces a new resource.",
				Validators:          []validator.String{stringvalidator.OneOf("SYSCTL", "UDEV", "ZFS")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"var": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name or identifier of the system parameter to tune. Changing this forces a new resource.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Value to assign to the tunable parameter.",
			},
			"comment": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Optional descriptive comment explaining the purpose of this tunable. Defaults to `\"\"`.",
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether this tunable is active and should be applied. Defaults to `true`.",
			},
			"update_initramfs": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "If `false`, then initramfs will not be updated after creating a ZFS tunable and you will need to run     `system boot update_initramfs` manually. Defaults to `true`.",
			},
			"orig_value": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Original system value of the parameter before this tunable was applied.",
			},
		},
	}
}
