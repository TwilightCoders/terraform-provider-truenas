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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"math/big"
)

var modelVM = &model{
	typeName:     "vm",
	namespace:    "vm",
	primaryKey:   "id",
	idKind:       kindInt,
	getMethod:    "vm.get_instance",
	deleteMethod: "vm.delete",
	createMethod: "vm.create",
	updateMethod: "vm.update",
	attrs: []*node{
		{
			name: "id", api: "id", path: "id",
			kind: kindInt, role: roleComputed,
			stable: true, readable: true, identity: true,
			description: "Identifier assigned by TrueNAS.",
		},
		{
			name: "command_line_args", api: "command_line_args", path: "command_line_args",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Additional command line arguments passed to the VM hypervisor. Defaults to `\"\"`.",
		},
		{
			name: "cpu_mode", api: "cpu_mode", path: "cpu_mode",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "CUSTOM",
			description: "CPU virtualization mode.\n\n* `CUSTOM`: Use specified model.\n* `HOST-MODEL`: Mirror host CPU.\n* `HOST-PASSTHROUGH`: Provide direct access to host CPU features. Defaults to `\"CUSTOM\"`.",
		},
		{
			name: "cpu_model", api: "cpu_model", path: "cpu_model",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Specific CPU model to emulate. `null` to use hypervisor default.",
		},
		{
			name: "name", api: "name", path: "name",
			kind: kindString, role: roleRequired,
			readable: true, updatable: true,
			description: "Display name of the virtual machine.",
		},
		{
			name: "description", api: "description", path: "description",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "",
			description: "Optional description or notes about the virtual machine. Defaults to `\"\"`.",
		},
		{
			name: "vcpus", api: "vcpus", path: "vcpus",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "1",
			description: "Number of virtual CPUs allocated to the VM. Defaults to `1`.",
		},
		{
			name: "cores", api: "cores", path: "cores",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "1",
			description: "Number of CPU cores per socket. Defaults to `1`.",
		},
		{
			name: "threads", api: "threads", path: "threads",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "1",
			description: "Number of threads per CPU core. Defaults to `1`.",
		},
		{
			name: "cpuset", api: "cpuset", path: "cpuset",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Set of host CPU cores to pin VM CPUs to. `null` for no pinning.",
		},
		{
			name: "nodeset", api: "nodeset", path: "nodeset",
			kind: kindString, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints.",
		},
		{
			name: "enable_cpu_topology_extension", api: "enable_cpu_topology_extension", path: "enable_cpu_topology_extension",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to expose detailed CPU topology information to the guest OS. Defaults to `false`.",
		},
		{
			name: "pin_vcpus", api: "pin_vcpus", path: "pin_vcpus",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexibility. Defaults to `false`.",
		},
		{
			name: "suspend_on_snapshot", api: "suspend_on_snapshot", path: "suspend_on_snapshot",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to suspend the VM when taking snapshots. Defaults to `false`.",
		},
		{
			name: "trusted_platform_module", api: "trusted_platform_module", path: "trusted_platform_module",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to enable virtual Trusted Platform Module (TPM) for the VM. Defaults to `false`.",
		},
		{
			name: "memory", api: "memory", path: "memory",
			kind: kindInt, role: roleRequired,
			readable: true, updatable: true,
			description: "Amount of memory allocated to the VM in megabytes.",
		},
		{
			name: "min_memory", api: "min_memory", path: "min_memory",
			kind: kindInt, role: roleOptional,
			nullable: true, readable: true, updatable: true,
			description: "Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink     during low usage but guarantees this minimum. `null` to disable ballooning.",
		},
		{
			name: "hyperv_enlightenments", api: "hyperv_enlightenments", path: "hyperv_enlightenments",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to enable Hyper-V enlightenments for improved Windows guest performance. Defaults to `false`.",
		},
		{
			name: "bootloader", api: "bootloader", path: "bootloader",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "UEFI",
			description: "Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility. Defaults to `\"UEFI\"`.",
		},
		{
			name: "bootloader_ovmf", api: "bootloader_ovmf", path: "bootloader_ovmf",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, replace: true, stable: true, readable: true,
			description: "OVMF firmware file to use for UEFI boot. Changing this forces a new resource.",
		},
		{
			name: "autostart", api: "autostart", path: "autostart",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether to automatically start the VM when the host system boots. Defaults to `true`.",
		},
		{
			name: "hide_from_msr", api: "hide_from_msr", path: "hide_from_msr",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: false,
			description: "Whether to hide hypervisor signatures from guest OS MSR access. Defaults to `false`.",
		},
		{
			name: "ensure_display_device", api: "ensure_display_device", path: "ensure_display_device",
			kind: kindBool, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: true,
			description: "Whether to ensure at least one display device is configured for the VM. Defaults to `true`.",
		},
		{
			name: "time", api: "time", path: "time",
			kind: kindString, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "LOCAL",
			description: "Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time. Defaults to `\"LOCAL\"`.",
		},
		{
			name: "shutdown_timeout", api: "shutdown_timeout", path: "shutdown_timeout",
			kind: kindInt, role: roleOptionalComputed,
			readable: true, updatable: true,
			hasDefault: true, def: "90",
			description: "Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances     allowing sufficient time for clean shutdown while avoiding indefinite hangs. Defaults to `90`.",
		},
		{
			name: "arch_type", api: "arch_type", path: "arch_type",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Guest architecture type. `null` to use hypervisor default.",
		},
		{
			name: "machine_type", api: "machine_type", path: "machine_type",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Virtual machine type/chipset. `null` to use hypervisor default.",
		},
		{
			name: "uuid", api: "uuid", path: "uuid",
			kind: kindString, role: roleOptionalComputed,
			nullable: true, readable: true, updatable: true,
			description: "Unique UUID for the VM. `null` to auto-generate.",
		},
		{
			name: "enable_secure_boot", api: "enable_secure_boot", path: "enable_secure_boot",
			kind: kindBool, role: roleOptionalComputed,
			replace: true, stable: true, readable: true,
			description: "Whether to enable UEFI Secure Boot for enhanced security. Changing this forces a new resource.",
		},
		{
			name: "display_available", api: "display_available", path: "display_available",
			kind: kindBool, role: roleComputed,
			readable:    true,
			description: "Whether at least one display device is available for this VM.",
		},
		{
			name: "status", api: "status", path: "status",
			kind: kindObject, role: roleComputed,
			readable:    true,
			description: "Current runtime status information for the VM.",
			children: []*node{
				{
					name: "state", api: "state", path: "status.state",
					kind: kindString, role: roleComputed,
					readable:    true,
					description: "Current state of the virtual machine.",
				},
				{
					name: "pid", api: "pid", path: "status.pid",
					kind: kindInt, role: roleComputed,
					nullable: true, readable: true,
					description: "Process ID of the running VM. `null` if not running.",
				},
				{
					name: "domain_state", api: "domain_state", path: "status.domain_state",
					kind: kindString, role: roleComputed,
					readable:    true,
					description: "Hypervisor-specific domain state.",
				},
			},
		},
	},
}

// NewVM returns the vm resource.
func NewVM() resource.Resource {
	return &vMResource{crudResource{model: modelVM}}
}

type vMResource struct{ crudResource }

// NewVMDataSource returns the vm data source.
func NewVMDataSource() datasource.DataSource {
	return &dataSource{model: modelVM, attrs: dataAttrs(modelVM)}
}

// NewVMList returns the vm list resource, for terraform query.
func NewVMList() list.ListResource {
	return &listResource{crudResource{model: modelVM}}
}

func (r *vMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual machine. Attach disks, NICs and displays with `truenas_vm_device`. Destroying the VM never deletes zvols attached to it.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Identifier assigned by TrueNAS.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"command_line_args": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Additional command line arguments passed to the VM hypervisor. Defaults to `\"\"`.",
				Default:             stringdefault.StaticString(""),
			},
			"cpu_mode": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "CPU virtualization mode.\n\n* `CUSTOM`: Use specified model.\n* `HOST-MODEL`: Mirror host CPU.\n* `HOST-PASSTHROUGH`: Provide direct access to host CPU features. Defaults to `\"CUSTOM\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("CUSTOM", "HOST-MODEL", "HOST-PASSTHROUGH")},
				Default:             stringdefault.StaticString("CUSTOM"),
			},
			"cpu_model": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specific CPU model to emulate. `null` to use hypervisor default.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Display name of the virtual machine.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Optional description or notes about the virtual machine. Defaults to `\"\"`.",
				Default:             stringdefault.StaticString(""),
			},
			"vcpus": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Number of virtual CPUs allocated to the VM. Defaults to `1`.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}},
				Default:             numberdefault.StaticBigFloat(big.NewFloat(1)),
			},
			"cores": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Number of CPU cores per socket. Defaults to `1`.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}},
				Default:             numberdefault.StaticBigFloat(big.NewFloat(1)),
			},
			"threads": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Number of threads per CPU core. Defaults to `1`.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 1}},
				Default:             numberdefault.StaticBigFloat(big.NewFloat(1)),
			},
			"cpuset": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Set of host CPU cores to pin VM CPUs to. `null` for no pinning.",
			},
			"nodeset": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Set of NUMA nodes to constrain VM memory allocation. `null` for no constraints.",
			},
			"enable_cpu_topology_extension": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to expose detailed CPU topology information to the guest OS. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"pin_vcpus": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to pin virtual CPUs to specific host CPU cores. Improves performance but reduces host flexibility. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"suspend_on_snapshot": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to suspend the VM when taking snapshots. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"trusted_platform_module": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable virtual Trusted Platform Module (TPM) for the VM. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"memory": schema.NumberAttribute{
				Required:            true,
				MarkdownDescription: "Amount of memory allocated to the VM in megabytes.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 20}},
			},
			"min_memory": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum memory allocation for dynamic memory ballooning in megabytes. Allows VM memory to shrink     during low usage but guarantees this minimum. `null` to disable ballooning.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 20}},
			},
			"hyperv_enlightenments": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable Hyper-V enlightenments for improved Windows guest performance. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"bootloader": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Boot firmware type. `UEFI` for modern UEFI, `UEFI_CSM` for legacy BIOS compatibility. Defaults to `\"UEFI\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("UEFI_CSM", "UEFI")},
				Default:             stringdefault.StaticString("UEFI"),
			},
			"bootloader_ovmf": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "OVMF firmware file to use for UEFI boot. Changing this forces a new resource.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
			},
			"autostart": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to automatically start the VM when the host system boots. Defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"hide_from_msr": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to hide hypervisor signatures from guest OS MSR access. Defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"ensure_display_device": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to ensure at least one display device is configured for the VM. Defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"time": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Guest OS time zone reference. `LOCAL` uses host timezone, `UTC` uses coordinated universal time. Defaults to `\"LOCAL\"`.",
				Validators:          []validator.String{stringvalidator.OneOf("LOCAL", "UTC")},
				Default:             stringdefault.StaticString("LOCAL"),
			},
			"shutdown_timeout": schema.NumberAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Maximum time in seconds to wait for graceful shutdown before forcing power off. Default 90s balances     allowing sufficient time for clean shutdown while avoiding indefinite hangs. Defaults to `90`.",
				Validators:          []validator.Number{numberIsInteger{}, numberAtLeast{min: 5}, numberAtMost{max: 300}},
				Default:             numberdefault.StaticBigFloat(big.NewFloat(90)),
			},
			"arch_type": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Guest architecture type. `null` to use hypervisor default.",
			},
			"machine_type": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Virtual machine type/chipset. `null` to use hypervisor default.",
			},
			"uuid": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Unique UUID for the VM. `null` to auto-generate.",
			},
			"enable_secure_boot": schema.BoolAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Whether to enable UEFI Secure Boot for enhanced security. Changing this forces a new resource.",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace(), boolplanmodifier.UseStateForUnknown()},
			},
			"display_available": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether at least one display device is available for this VM.",
			},
			"status": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Current runtime status information for the VM.",
				Attributes: map[string]schema.Attribute{
					"state": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Current state of the virtual machine.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"pid": schema.NumberAttribute{
						Computed:            true,
						MarkdownDescription: "Process ID of the running VM. `null` if not running.",
						Validators:          []validator.Number{numberIsInteger{}},
					},
					"domain_state": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Hypervisor-specific domain state.",
						Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
		},
	}
}
