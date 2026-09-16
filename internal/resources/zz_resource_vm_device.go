// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var modelVMDevice = &model{
	typeName:     "vm_device",
	namespace:    "vm.device",
	primaryKey:   "id",
	idKind:       kindInt,
	updateMethod: "vm.device.update",
	getMethod:    "vm.device.get_instance",
	deleteMethod: "vm.device.delete",
	createMethod: "vm.device.create",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "attributes", api: "attributes", path: "attributes",
			kind: kindUnion, role: roleRequired,
			readable: true, updatable: true,
			description:   "Device-specific configuration attributes.",
			discriminator: "dtype",
			children: []*node{
				{
					name: "cdrom", api: "CDROM", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `dtype` is `CDROM`.",
					children: []*node{
						{
							name: "path", api: "path", path: "attributes.path",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Path must not contain \"{\", \"}\" characters, and it should start with \"/mnt/\".",
						},
					},
				},
				{
					name: "display", api: "DISPLAY", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `dtype` is `DISPLAY`.",
					children: []*node{
						{
							name: "resolution", api: "resolution", path: "attributes.resolution",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "1024x768",
							description: "Screen resolution for the virtual display. Defaults to `\"1024x768\"`.",
						},
						{
							name: "port", api: "port", path: "attributes.port",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "VNC/SPICE port number for remote display access. `null` for auto-assignment.",
						},
						{
							name: "web_port", api: "web_port", path: "attributes.web_port",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "Web-based display access port number. `null` for auto-assignment.",
						},
						{
							name: "bind", api: "bind", path: "attributes.bind",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "127.0.0.1",
							description: "IP address to bind the display server to. Defaults to `\"127.0.0.1\"`.",
						},
						{
							name: "wait", api: "wait", path: "attributes.wait",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Whether to wait for a client connection before starting the VM. Defaults to `false`.",
						},
						{
							name: "password", api: "password", path: "attributes.password",
							kind: kindString, role: roleOptional,
							nullable: true, sensitive: true, writeOnly: true, readable: true,
							description: "Password for display server authentication.",
						},
						{
							name: "web", api: "web", path: "attributes.web",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: true,
							description: "Whether to enable web-based display access. Defaults to `true`.",
						},
						{
							name: "type", api: "type", path: "attributes.type",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "SPICE",
							description: "Display protocol type. Defaults to `\"SPICE\"`.",
						},
					},
				},
				{
					name: "nic", api: "NIC", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `dtype` is `NIC`.",
					children: []*node{
						{
							name: "trust_guest_rx_filters", api: "trust_guest_rx_filters", path: "attributes.trust_guest_rx_filters",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Whether to trust guest OS receive filter settings for better performance. Defaults to `false`.",
						},
						{
							name: "type", api: "type", path: "attributes.type",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "E1000",
							description: "Network interface controller type. `E1000` for Intel compatibility, `VIRTIO` for performance. Defaults to `\"E1000\"`.",
						},
						{
							name: "nic_attach", api: "nic_attach", path: "attributes.nic_attach",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "Host network interface or bridge to attach to. `null` for no attachment.",
						},
						{
							name: "mac", api: "mac", path: "attributes.mac",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "MAC address for the virtual network interface. `null` for auto-generation.",
						},
					},
				},
				{
					name: "pci", api: "PCI", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `dtype` is `PCI`.",
					children: []*node{
						{
							name: "pptdev", api: "pptdev", path: "attributes.pptdev",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Host PCI device identifier to pass through to the VM.",
						},
					},
				},
				{
					name: "raw", api: "RAW", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `dtype` is `RAW`.",
					children: []*node{
						{
							name: "path", api: "path", path: "attributes.path",
							kind: kindString, role: roleRequired,
							readable:    true,
							description: "Path must not contain \"{\", \"}\" characters.",
						},
						{
							name: "type", api: "type", path: "attributes.type",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "AHCI",
							description: "Disk controller interface type. AHCI for compatibility, VIRTIO for performance. Defaults to `\"AHCI\"`.",
						},
						{
							name: "exists", api: "exists", path: "attributes.exists",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Whether the disk file already exists or should be created. Defaults to `false`.",
						},
						{
							name: "boot", api: "boot", path: "attributes.boot",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Whether this disk should be marked as bootable. Defaults to `false`.",
						},
						{
							name: "size", api: "size", path: "attributes.size",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "Size of the disk in bytes. Required if creating a new disk file.",
						},
						{
							name: "logical_sectorsize", api: "logical_sectorsize", path: "attributes.logical_sectorsize",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "Logical sector size for the disk. `null` for default.",
						},
						{
							name: "physical_sectorsize", api: "physical_sectorsize", path: "attributes.physical_sectorsize",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "Physical sector size for the disk. `null` for default.",
						},
						{
							name: "iotype", api: "iotype", path: "attributes.iotype",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "THREADS",
							description: "I/O backend type for disk operations. Defaults to `\"THREADS\"`.",
						},
						{
							name: "serial", api: "serial", path: "attributes.serial",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "Serial number to assign to the virtual disk. `null` for auto-generated.",
						},
					},
				},
				{
					name: "disk", api: "DISK", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `dtype` is `DISK`.",
					children: []*node{
						{
							name: "path", api: "path", path: "attributes.path",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "Path to existing disk file or ZFS volume. `null` if creating a new ZFS volume.",
						},
						{
							name: "type", api: "type", path: "attributes.type",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "AHCI",
							description: "Disk controller interface type. AHCI for compatibility, VIRTIO for performance. Defaults to `\"AHCI\"`.",
						},
						{
							name: "create_zvol", api: "create_zvol", path: "attributes.create_zvol",
							kind: kindBool, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: false,
							description: "Whether to create a new ZFS volume for this disk. Defaults to `false`.",
						},
						{
							name: "zvol_name", api: "zvol_name", path: "attributes.zvol_name",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "Name for the new ZFS volume. Required if `create_zvol` is true.",
						},
						{
							name: "zvol_volsize", api: "zvol_volsize", path: "attributes.zvol_volsize",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "Size of the new ZFS volume in bytes. Required if `create_zvol` is true.",
						},
						{
							name: "logical_sectorsize", api: "logical_sectorsize", path: "attributes.logical_sectorsize",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "Logical sector size for the disk. `null` for default.",
						},
						{
							name: "physical_sectorsize", api: "physical_sectorsize", path: "attributes.physical_sectorsize",
							kind: kindInt, role: roleOptional,
							nullable: true, readable: true,
							description: "Physical sector size for the disk. `null` for default.",
						},
						{
							name: "iotype", api: "iotype", path: "attributes.iotype",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "THREADS",
							description: "I/O backend type for disk operations. Defaults to `\"THREADS\"`.",
						},
						{
							name: "serial", api: "serial", path: "attributes.serial",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "Serial number to assign to the virtual disk. `null` for auto-generated.",
						},
					},
				},
				{
					name: "usb", api: "USB", path: "attributes",
					kind: kindObject, role: roleOptional,
					nullable: true, readable: true,
					description: "Settings when `dtype` is `USB`.",
					children: []*node{
						{
							name: "usb", api: "usb", path: "attributes.usb",
							kind: kindObject, role: roleOptional,
							nullable: true, readable: true,
							description: "USB device attributes for identification. `null` for USB host controller only.",
							children: []*node{
								{
									name: "vendor_id", api: "vendor_id", path: "attributes.usb.vendor_id",
									kind: kindString, role: roleRequired,
									readable:    true,
									description: "USB vendor identifier in hexadecimal format (e.g., '0x1d6b' for Linux Foundation).",
								},
								{
									name: "product_id", api: "product_id", path: "attributes.usb.product_id",
									kind: kindString, role: roleRequired,
									readable:    true,
									description: "USB product identifier in hexadecimal format (e.g., '0x0002' for 2.0 root hub).",
								},
							},
						},
						{
							name: "controller_type", api: "controller_type", path: "attributes.controller_type",
							kind: kindString, role: roleOptionalComputed,
							readable:   true,
							hasDefault: true, def: "nec-xhci",
							description: "USB controller type for the virtual machine. Defaults to `\"nec-xhci\"`.",
						},
						{
							name: "device", api: "device", path: "attributes.device",
							kind: kindString, role: roleOptional,
							nullable: true, readable: true,
							description: "Host USB device path to pass through. `null` for controller only.",
						},
					},
				},
			},
		},
		{
			name: "vm", api: "vm", path: "vm",
			kind: kindInt, role: roleRequired,
			readable: true, updatable: true,
			description: "ID of the virtual machine this device belongs to.",
		},
		{
			name: "order", api: "order", path: "order",
			kind: kindInt, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Boot order priority for this device. `null` for automatic assignment.",
		},
		{
			name: "attributes_wo_version", api: "", path: "",
			kind: kindInt, role: roleOptional,
			description: "Change this value to send `attributes` again. Terraform never stores `attributes`.",
		},
	},
}

// NewVMDevice returns the vm_device resource.
func NewVMDevice() resource.Resource {
	return &vMDeviceResource{crudResource{model: modelVMDevice}}
}

type vMDeviceResource struct{ crudResource }

// NewVMDeviceDataSource returns the vm_device data source.
func NewVMDeviceDataSource() datasource.DataSource {
	return &dataSource{model: modelVMDevice, attrs: dataAttrs(modelVMDevice)}
}

// NewVMDeviceList returns the vm_device list resource, for terraform query.
func NewVMDeviceList() list.ListResource {
	return &listResource{crudResource{model: modelVMDevice}}
}

func (r *vMDeviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual machine device. Set exactly one device type under `attributes`. Display passwords are write-only. Most device changes require the VM to be stopped.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"attributes": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Device-specific configuration attributes.",
				Attributes: map[string]schema.Attribute{
					"cdrom": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `dtype` is `CDROM`.",
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Path must not contain \"{\", \"}\" characters, and it should start with \"/mnt/\".",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"display": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `dtype` is `DISPLAY`.",
						Attributes: map[string]schema.Attribute{
							"resolution": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Screen resolution for the virtual display. Defaults to `\"1024x768\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("1920x1200", "1920x1080", "1600x1200", "1600x900", "1400x1050", "1280x1024", "1280x720", "1024x768", "800x600", "640x480")},
								Default:             stringdefault.StaticString("1024x768"),
							},
							"port": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "VNC/SPICE port number for remote display access. `null` for auto-assignment.",
								Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 5900}, numberAtMost{max: 65535}},
							},
							"web_port": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Web-based display access port number. `null` for auto-assignment.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"bind": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "IP address to bind the display server to. Defaults to `\"127.0.0.1\"`.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
								Default:             stringdefault.StaticString("127.0.0.1"),
							},
							"wait": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether to wait for a client connection before starting the VM. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"password": schema.StringAttribute{
								Optional: true, Sensitive: true, WriteOnly: true,
								MarkdownDescription: "Password for display server authentication.",
							},
							"web": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether to enable web-based display access. Defaults to `true`.",
								Default:             booldefault.StaticBool(true),
							},
							"type": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Display protocol type. Defaults to `\"SPICE\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("SPICE", "VNC")},
								Default:             stringdefault.StaticString("SPICE"),
							},
						},
					},
					"nic": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `dtype` is `NIC`.",
						Attributes: map[string]schema.Attribute{
							"trust_guest_rx_filters": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether to trust guest OS receive filter settings for better performance. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"type": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Network interface controller type. `E1000` for Intel compatibility, `VIRTIO` for performance. Defaults to `\"E1000\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("E1000", "VIRTIO")},
								Default:             stringdefault.StaticString("E1000"),
							},
							"nic_attach": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Host network interface or bridge to attach to. `null` for no attachment.",
							},
							"mac": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "MAC address for the virtual network interface. `null` for auto-generation.",
							},
						},
					},
					"pci": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `dtype` is `PCI`.",
						Attributes: map[string]schema.Attribute{
							"pptdev": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Host PCI device identifier to pass through to the VM.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"raw": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `dtype` is `RAW`.",
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Path must not contain \"{\", \"}\" characters.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"type": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Disk controller interface type. AHCI for compatibility, VIRTIO for performance. Defaults to `\"AHCI\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("AHCI", "VIRTIO")},
								Default:             stringdefault.StaticString("AHCI"),
							},
							"exists": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether the disk file already exists or should be created. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"boot": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether this disk should be marked as bootable. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"size": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Size of the disk in bytes. Required if creating a new disk file.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"logical_sectorsize": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Logical sector size for the disk. `null` for default.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"physical_sectorsize": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Physical sector size for the disk. `null` for default.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"iotype": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "I/O backend type for disk operations. Defaults to `\"THREADS\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("NATIVE", "THREADS", "IO_URING")},
								Default:             stringdefault.StaticString("THREADS"),
							},
							"serial": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Serial number to assign to the virtual disk. `null` for auto-generated.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"disk": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `dtype` is `DISK`.",
						Attributes: map[string]schema.Attribute{
							"path": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Path to existing disk file or ZFS volume. `null` if creating a new ZFS volume.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"type": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Disk controller interface type. AHCI for compatibility, VIRTIO for performance. Defaults to `\"AHCI\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("AHCI", "VIRTIO")},
								Default:             stringdefault.StaticString("AHCI"),
							},
							"create_zvol": schema.BoolAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "Whether to create a new ZFS volume for this disk. Defaults to `false`.",
								Default:             booldefault.StaticBool(false),
							},
							"zvol_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Name for the new ZFS volume. Required if `create_zvol` is true.",
							},
							"zvol_volsize": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Size of the new ZFS volume in bytes. Required if `create_zvol` is true.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"logical_sectorsize": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Logical sector size for the disk. `null` for default.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"physical_sectorsize": schema.NumberAttribute{
								Optional:            true,
								MarkdownDescription: "Physical sector size for the disk. `null` for default.",
								Validators:          []validator.Number{numberIsInteger{}},
							},
							"iotype": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "I/O backend type for disk operations. Defaults to `\"THREADS\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("NATIVE", "THREADS", "IO_URING")},
								Default:             stringdefault.StaticString("THREADS"),
							},
							"serial": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Serial number to assign to the virtual disk. `null` for auto-generated.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
					"usb": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Settings when `dtype` is `USB`.",
						Attributes: map[string]schema.Attribute{
							"usb": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "USB device attributes for identification. `null` for USB host controller only.",
								Attributes: map[string]schema.Attribute{
									"vendor_id": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "USB vendor identifier in hexadecimal format (e.g., '0x1d6b' for Linux Foundation).",
										Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
									},
									"product_id": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "USB product identifier in hexadecimal format (e.g., '0x0002' for 2.0 root hub).",
										Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
									},
								},
							},
							"controller_type": schema.StringAttribute{
								Optional: true, Computed: true,
								MarkdownDescription: "USB controller type for the virtual machine. Defaults to `\"nec-xhci\"`.",
								Validators:          []validator.String{stringvalidator.OneOf("piix3-uhci", "piix4-uhci", "ehci", "ich9-ehci1", "vt82c686b-uhci", "pci-ohci", "nec-xhci", "qemu-xhci")},
								Default:             stringdefault.StaticString("nec-xhci"),
							},
							"device": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Host USB device path to pass through. `null` for controller only.",
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
				},
			},
			"vm": schema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "ID of the virtual machine this device belongs to.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"order": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Boot order priority for this device. `null` for automatic assignment.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
			"attributes_wo_version": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Change this value to send `attributes` again. Terraform never stores `attributes`.",
				Validators:          []validator.Number{numberIsInteger{}},
			},
		},
	}
}
