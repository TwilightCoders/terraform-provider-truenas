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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var modelISCSIExtent = &model{
	typeName:     "iscsi_extent",
	namespace:    "iscsi.extent",
	primaryKey:   "id",
	idKind:       kindInt,
	deleteMethod: "iscsi.extent.delete",
	createMethod: "iscsi.extent.create",
	updateMethod: "iscsi.extent.update",
	getMethod:    "iscsi.extent.get_instance",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Name of the iSCSI extent.",
		},
		{
			name: "type", api: "type", path: "type",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "DISK",
			description: "Type of the extent storage backend. Defaults to `\"DISK\"`.",
		},
		{
			name: "disk", api: "disk", path: "disk",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Disk device to use for the extent or `null` if using a file.",
		},
		{
			name: "serial", api: "serial", path: "serial",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Serial number for the extent or `null` to auto-generate.",
		},
		{
			name: "path", api: "path", path: "path",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "File path for file-based extents or `null` if using a disk.",
		},
		{
			name: "filesize", api: "filesize", path: "filesize",
			kind: kindAny, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "0",
			description: "Size of the file-based extent in bytes. Defaults to `\"0\"`.",
		},
		{
			name: "blocksize", api: "blocksize", path: "blocksize",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: 512,
			description: "Block size for the extent in bytes. Defaults to `512`.",
		},
		{
			name: "pblocksize", api: "pblocksize", path: "pblocksize",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to use physical block size reporting. Defaults to `false`.",
		},
		{
			name: "avail_threshold", api: "avail_threshold", path: "avail_threshold",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Available space threshold percentage or `null` to disable.",
		},
		{
			name: "comment", api: "comment", path: "comment",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Optional comment describing the extent. Defaults to `\"\"`.",
		},
		{
			name: "insecure_tpc", api: "insecure_tpc", path: "insecure_tpc",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether to enable insecure Third Party Copy (TPC) operations. Defaults to `true`.",
		},
		{
			name: "xen", api: "xen", path: "xen",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to enable Xen compatibility mode. Defaults to `false`.",
		},
		{
			name: "rpm", api: "rpm", path: "rpm",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "SSD",
			description: "Reported RPM type for the extent. Defaults to `\"SSD\"`.",
		},
		{
			name: "ro", api: "ro", path: "ro",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether the extent is read-only. Defaults to `false`.",
		},
		{
			name: "enabled", api: "enabled", path: "enabled",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether the extent is enabled and available for use. Defaults to `true`.",
		},
		{
			name: "product_id", api: "product_id", path: "product_id",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Product ID string for the extent or `null` for default.",
		},
		{
			name: "naa", api: "naa", path: "naa",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Network Address Authority (NAA) identifier for the extent.",
		},
		{
			name: "vendor", api: "vendor", path: "vendor",
			kind: kindString, role: roleComputed,
			readable:    true,
			description: "Vendor string reported by the extent.",
		},
		{
			name: "locked", api: "locked", path: "locked",
			kind: kindBool, role: roleComputed,
			nullable: true, readable: true,
			description: "Read-only value indicating whether the iscsi extent is located on a locked dataset.\n\n- `true`: The extent is in a locked dataset.\n- `false`: The extent is not in a locked dataset.\n- `null`: Lock status is not available because path locking information was not requested.",
		},
		{
			name: "unset", api: "", path: "unset",
			kind: kindList, role: roleOptional,
			description: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
			elem: &node{
				name: "unset", api: "", path: "unset",
				kind: kindString, role: roleRequired,
			},
		},
	},
}

// NewISCSIExtent returns the iscsi_extent resource.
func NewISCSIExtent() resource.Resource {
	return &iSCSIExtentResource{crudResource{model: modelISCSIExtent}}
}

type iSCSIExtentResource struct{ crudResource }

// NewISCSIExtentDataSource returns the iscsi_extent data source.
func NewISCSIExtentDataSource() datasource.DataSource {
	return &dataSource{model: modelISCSIExtent, attrs: dataAttrs(modelISCSIExtent)}
}

// NewISCSIExtentList returns the iscsi_extent list resource, for terraform query.
func NewISCSIExtentList() list.ListResource {
	return &listResource{crudResource{model: modelISCSIExtent}}
}

func (r *iSCSIExtentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an iSCSI Extent.\n\nWhen `type` is set to FILE, attribute `filesize` is used and it represents number of bytes. `filesize` if\nnot zero should be a multiple of `blocksize`. `path` is a required attribute with `type` set as FILE.\n\nWith `type` being set to DISK, a valid ZFS volume is required.\n\n`insecure_tpc` when enabled allows an initiator to bypass normal access control and access any scannable\ntarget. This allows xcopy operations otherwise blocked by access control.\n\n`xen` is a boolean value which is set to true if Xen is being used as the iSCSI initiator.\n\n`ro` when set to true prevents the initiator from writing to this LUN.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the iSCSI extent.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.LengthAtMost(64)},
			},
			"type": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Type of the extent storage backend. Defaults to `\"DISK\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("DISK", "FILE")},
			},
			"disk": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Disk device to use for the extent or `null` if using a file.",
			},
			"serial": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Serial number for the extent or `null` to auto-generate.",
			},
			"path": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "File path for file-based extents or `null` if using a disk.",
			},
			"filesize": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Size of the file-based extent in bytes. Defaults to `\"0\"`.",
			},
			"blocksize": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Block size for the extent in bytes. Defaults to `512`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"pblocksize": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to use physical block size reporting. Defaults to `false`.",
			},
			"avail_threshold": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Available space threshold percentage or `null` to disable.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}, numberAtMost{max: 99}},
			},
			"comment": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Optional comment describing the extent. Defaults to `\"\"`.",
			},
			"insecure_tpc": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable insecure Third Party Copy (TPC) operations. Defaults to `true`.",
			},
			"xen": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable Xen compatibility mode. Defaults to `false`.",
			},
			"rpm": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Reported RPM type for the extent. Defaults to `\"SSD\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("UNKNOWN", "SSD", "5400", "7200", "10000", "15000")},
			},
			"ro": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether the extent is read-only. Defaults to `false`.",
			},
			"enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether the extent is enabled and available for use. Defaults to `true`.",
			},
			"product_id": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Product ID string for the extent or `null` for default.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1), stringvalidator.LengthAtMost(16)},
			},
			"naa": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Network Address Authority (NAA) identifier for the extent.",
				Validators:          []validator.String{stringvalidator.LengthAtMost(34)},
			},
			"vendor": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Vendor string reported by the extent.",
			},
			"locked": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Read-only value indicating whether the iscsi extent is located on a locked dataset.\n\n- `true`: The extent is in a locked dataset.\n- `false`: The extent is not in a locked dataset.\n- `null`: Lock status is not available because path locking information was not requested.",
			},
			"unset": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Attributes to clear, by name. An attribute this configuration does not mention is left as the server has it; naming it here removes the value instead. Only attributes that accept an empty value can be listed.",
				ElementType:         types.StringType,
			},
		},
	}
}
