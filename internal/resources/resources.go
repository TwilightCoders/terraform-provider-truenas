// Package resources declares the provider's resources as engine specs.
package resources

import "github.com/TwilightCoders/terraform-provider-truenas/internal/engine"

// All lists every resource the provider serves.
var All = []engine.Resource{
	CronJob,
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
