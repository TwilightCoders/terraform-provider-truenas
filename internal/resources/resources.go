// Package resources declares the provider's resources as engine specs.
package resources

import "github.com/TwilightCoders/terraform-provider-truenas/internal/engine"

// All lists every resource the provider serves.
var All = []engine.Resource{
	CronJob,
	Group,
	InitScript,
	SMBShare,
	SnapshotTask,
	User,
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
