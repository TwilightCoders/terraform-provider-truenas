// Package resources declares the provider's resources as engine specs.
package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/engine"
)

// Custom lists hand-written resources that the engine cannot derive.
var Custom = []func() resource.Resource{
	NewApp,
	NewFilesystemACL,
}

// All lists every resource the provider serves.
var All = []engine.Resource{
	ACMEDNSAuthenticator,
	Certificate,
	CloudSyncCredentials,
	CloudSyncTask,
	CronJob,
	Dataset,
	GeneralConfig,
	Group,
	InitScript,
	NetworkConfig,
	NFSConfig,
	Replication,
	SMBConfig,
	SMBShare,
	SnapshotTask,
	SSHConfig,
	User,
	VM,
	VMDevice,
	Zvol,
}

// CronJob manages scheduled commands (cronjob).
var CronJob = engine.Resource{
	Type:      "cron_job",
	Namespace: "cronjob",
	Fields: map[string]engine.Field{
		// The middleware fields mean "ignore output", which reads backwards as "stdout = true".
		"stdout": {Name: "ignore_stdout"},
		"stderr": {Name: "ignore_stderr"},
	},
}

// Group manages local groups (group).
var Group = engine.Resource{
	Type:        "group",
	Namespace:   "group",
	ListFilters: [][]any{{"builtin", "=", false}},
	Fields: map[string]engine.Field{
		// "group" is a read-only alias of name.
		"group": engine.Omit,
	},
}

// InitScript manages commands and scripts run at boot or shutdown (initshutdownscript).
var InitScript = engine.Resource{
	Type:        "init_script",
	Namespace:   "initshutdownscript",
	Description: "Runs a command or script at `PREINIT`, `POSTINIT` or `SHUTDOWN`. Set `command` when `type` is `COMMAND` and `script` when it is `SCRIPT`.",
}

// SMBShare manages SMB shares (sharing.smb).
var SMBShare = engine.Resource{
	Type:      "smb_share",
	Namespace: "sharing.smb",
	Fields: map[string]engine.Field{
		// The server fills purpose-specific options when they are not given.
		"options": engine.Computed,
	},
}

// SnapshotTask manages periodic snapshot tasks (pool.snapshottask).
var SnapshotTask = engine.Resource{
	Type:      "snapshot_task",
	Namespace: "pool.snapshottask",
	Fields: map[string]engine.Field{
		// Run status with timestamps; it changes on every run and is not configuration.
		"state": engine.Omit,
	},
}

// User manages local user accounts (user).
var User = engine.Resource{
	Type:        "user",
	Namespace:   "user",
	ListFilters: [][]any{{"builtin", "=", false}},
	Fields: map[string]engine.Field{
		"password": engine.WriteOnly,
		// Password hashes; never useful in state.
		"unixhash": engine.Omit,
		"smbhash":  engine.Omit,
		// Generates a password the provider could only store in state.
		"random_password": engine.Omit,
		// Accepts a group id, returns the group object.
		"group": {Ref: "id"},
	},
}

// Dataset manages ZFS filesystems (pool.dataset with type FILESYSTEM).
var Dataset = engine.Resource{
	Type:        "dataset",
	Namespace:   "pool.dataset",
	Variant:     "FILESYSTEM",
	ListOptions: map[string]any{"extra": map[string]any{"flat": true}},
	Description: "Manages a ZFS filesystem dataset. Properties set to `INHERIT` (the default for most) follow the parent dataset. Destroying the resource fails while the dataset has children.",
	Fields: map[string]engine.Field{
		// Child datasets are resources of their own.
		"children": engine.Omit,
		// The encryption passphrase or key must never reach state.
		"encryption_options": {WriteOnlySecrets: true},
		// Read as a map of property wrappers but written as a list; not modeled yet.
		"user_properties": engine.Omit,
		// Volume-only properties that the shared read shape reports for filesystems too.
		"volsize":      engine.Omit,
		"volblocksize": engine.Omit,
		"sparse":       engine.Omit,
	},
}

// Zvol manages ZFS volumes (pool.dataset with type VOLUME).
var Zvol = engine.Resource{
	Type:        "zvol",
	Namespace:   "pool.dataset",
	Variant:     "VOLUME",
	ListOptions: map[string]any{"extra": map[string]any{"flat": true}},
	Description: "Manages a ZFS volume (zvol). Destroying the resource fails while the volume has snapshots or dependents.",
	Fields: map[string]engine.Field{
		"children":           engine.Omit,
		"user_properties":    engine.Omit,
		"mountpoint":         engine.Omit,
		"encryption_options": {WriteOnlySecrets: true},
	},
}

// NetworkConfig manages global network settings: hostname, domain, gateways, DNS (network.configuration).
var NetworkConfig = engine.Resource{
	Type:      "network_config",
	Namespace: "network.configuration",
	Fields: map[string]engine.Field{
		// Live interface and DHCP state, not configuration.
		"state": engine.Omit,
	},
}

// NFSConfig manages the NFS service settings (nfs).
var NFSConfig = engine.Resource{
	Type:      "nfs_config",
	Namespace: "nfs",
}

// SMBConfig manages the SMB service settings (smb).
var SMBConfig = engine.Resource{
	Type:      "smb_config",
	Namespace: "smb",
}

// SSHConfig manages the SSH service settings (ssh).
var SSHConfig = engine.Resource{
	Type:      "ssh_config",
	Namespace: "ssh",
	Fields: map[string]engine.Field{
		// Host private keys: never in Terraform state.
		"privatekey":       engine.Omit,
		"host_dsa_key":     engine.Omit,
		"host_ecdsa_key":   engine.Omit,
		"host_ed25519_key": engine.Omit,
		"host_key":         engine.Omit,
		"host_rsa_key":     engine.Omit,
	},
}

// CloudSyncCredentials manages credentials for cloud sync tasks (cloudsync.credentials).
var CloudSyncCredentials = engine.Resource{
	Type:      "cloudsync_credentials",
	Namespace: "cloudsync.credentials",
	Description: "Manages credentials for cloud sync tasks. Set exactly one storage type under `storage`. " +
		"Secret fields are write-only: they never enter Terraform state. Bump `storage_wo_version` to resend them.",
	Fields: map[string]engine.Field{
		// "provider" is reserved by Terraform.
		"provider": {Name: "storage", WriteOnlySecrets: true},
	},
}

// CloudSyncTask manages cloud sync (rclone) tasks (cloudsync).
var CloudSyncTask = engine.Resource{
	Type:      "cloudsync_task",
	Namespace: "cloudsync",
	Fields: map[string]engine.Field{
		"credentials":         {Ref: "id"},
		"encryption_password": engine.WriteOnly,
		"encryption_salt":     engine.WriteOnly,
		// Last run status, not configuration.
		"job": engine.Omit,
	},
}

// Replication manages ZFS replication tasks (replication).
var Replication = engine.Resource{
	Type:      "replication",
	Namespace: "replication",
	Fields: map[string]engine.Field{
		"ssh_credentials":         {Ref: "id"},
		"periodic_snapshot_tasks": {Ref: "id"},
		"encryption_key":          engine.WriteOnly,
		// Run status, not configuration.
		"state": engine.Omit,
		"job":   engine.Omit,
	},
}

// ACMEDNSAuthenticator manages DNS providers used to solve ACME DNS-01 challenges
// (acme.dns.authenticator).
var ACMEDNSAuthenticator = engine.Resource{
	Type:      "acme_dns_authenticator",
	Namespace: "acme.dns.authenticator",
	Description: "Manages a DNS provider for ACME DNS-01 challenges. Set exactly one provider under `authenticator`. " +
		"Credentials are write-only: they never enter Terraform state. Bump `authenticator_wo_version` to resend them.",
	Fields: map[string]engine.Field{
		"attributes": {Name: "authenticator", WriteOnlySecrets: true},
	},
}

// Certificate manages certificates: imported, CSRs, and ACME certificates issued from a CSR
// (certificate). TrueNAS 25.10 cannot create self-signed certificates or certificate authorities.
var Certificate = engine.Resource{
	Type:      "certificate",
	Namespace: "certificate",
	Description: "Manages a certificate. `create_type` selects how: `CERTIFICATE_CREATE_IMPORTED` (bring `certificate` and " +
		"`privatekey`), `CERTIFICATE_CREATE_CSR` (TrueNAS generates the key and a CSR), `CERTIFICATE_CREATE_IMPORTED_CSR`, or " +
		"`CERTIFICATE_CREATE_ACME` (issue from an existing CSR via `csr_id`, `dns_mapping` and `tos`). " +
		"Private keys never enter Terraform state. TrueNAS renews ACME certificates itself.",
	Fields: map[string]engine.Field{
		// TrueNAS returns the private key on read; keep it out of state.
		"privatekey": engine.WriteOnly,
		"passphrase": engine.WriteOnly,
	},
}

// GeneralConfig manages general system settings: web UI listeners and certificate, timezone,
// keyboard map (system.general).
var GeneralConfig = engine.Resource{
	Type:      "general_config",
	Namespace: "system.general",
	Fields: map[string]engine.Field{
		"ui_certificate": {Ref: "id"},
		"ui_restart_delay": {Description: "Seconds after the update to restart the web UI and apply UI settings. " +
			"Without it, UI changes take effect at the next UI restart. The restart aborts every HTTP connection, " +
			"including the provider's own; reads retry across it."},
	},
}

// VM manages bhyve/libvirt virtual machines (vm). Devices are separate truenas_vm_device resources.
var VM = engine.Resource{
	Type:      "vm",
	Namespace: "vm",
	Description: "Manages a virtual machine. Attach disks, NICs and displays with `truenas_vm_device`. " +
		"Destroying the VM never deletes zvols attached to it.",
	Fields: map[string]engine.Field{
		// Devices are managed as their own resources.
		"devices": engine.Omit,
		// Assigned by TrueNAS when not given.
		"uuid":         engine.Computed,
		"arch_type":    engine.Computed,
		"machine_type": engine.Computed,
	},
}

// VMDevice manages one device of a virtual machine: disk, raw file, CD-ROM, NIC, display, PCI or USB
// passthrough (vm.device).
var VMDevice = engine.Resource{
	Type:      "vm_device",
	Namespace: "vm.device",
	Description: "Manages a virtual machine device. Set exactly one device type under `attributes`. " +
		"Display passwords are write-only. Most device changes require the VM to be stopped.",
	Fields: map[string]engine.Field{
		"attributes": {WriteOnlySecrets: true},
	},
}
