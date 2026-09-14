// Package resources declares the provider's resources as engine specs.
package resources

import "github.com/TwilightCoders/terraform-provider-truenas/internal/engine"

// All lists every resource the provider serves.
var All = []engine.Resource{
	CronJob,
	Dataset,
	Group,
	InitScript,
	SMBShare,
	SnapshotTask,
	User,
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
	Type:      "group",
	Namespace: "group",
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
	Type:      "user",
	Namespace: "user",
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
	Description: "Manages a ZFS filesystem dataset. Properties set to `INHERIT` (the default for most) follow the parent dataset. Destroying the resource fails while the dataset has children.",
	Fields: map[string]engine.Field{
		// Child datasets are resources of their own.
		"children": engine.Omit,
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
	Description: "Manages a ZFS volume (zvol). Destroying the resource fails while the volume has snapshots or dependents.",
	Fields: map[string]engine.Field{
		"children":        engine.Omit,
		"user_properties": engine.Omit,
		"mountpoint":      engine.Omit,
	},
}
