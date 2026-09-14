package resources

import "github.com/TwilightCoders/terraform-provider-truenas/internal/engine"

// Actions lists every action the provider serves (Terraform 1.14 and later).
var Actions = []engine.Action{
	{
		Type: "replication_run", Method: "replication.run",
		Description: "Runs a replication task now and waits for it to finish.",
	},
	{
		Type: "cloudsync_sync", Method: "cloudsync.sync",
		Description: "Runs a cloud sync task now and waits for it to finish.",
		Args:        map[string]engine.Field{"cloud_sync_sync_options": {Name: "options"}},
	},
	{
		Type: "snapshot_task_run", Method: "pool.snapshottask.run",
		Description: "Takes the snapshot for a periodic snapshot task now.",
	},
	{
		Type: "scrub_run", Method: "pool.scrub.run",
		Description: "Starts a scrub of a pool when the last one is older than `threshold` days.",
	},
	{
		Type: "app_redeploy", Method: "app.redeploy",
		Description: "Redeploys an app, re-reading its compose configuration, and waits for it to finish.",
		Args:        map[string]engine.Field{"app_name": {Name: "name"}},
	},
	{
		Type: "app_start", Method: "app.start",
		Description: "Starts an app and waits for it to finish.",
		Args:        map[string]engine.Field{"app_name": {Name: "name"}},
	},
	{
		Type: "app_stop", Method: "app.stop",
		Description: "Stops an app and waits for it to finish.",
		Args:        map[string]engine.Field{"app_name": {Name: "name"}},
	},
	{
		Type: "service_control", Method: "service.control",
		Description: "Starts, stops, restarts or reloads a system service and waits for it to finish.",
	},
	{
		Type: "ui_restart", Method: "system.general.ui_restart",
		Description: "Restarts the web UI after `delay` seconds, applying UI settings. Aborts every HTTP connection, including the provider's.",
	},
}
