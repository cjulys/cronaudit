// Package tag provides tagging support for cron schedule entries.
//
// Tags are arbitrary strings attached to entries by their label,
// allowing consumers to group, filter, or annotate entries outside
// of the core scheduling logic.
//
// Example usage:
//
//	tr := tag.New(report)
//	tr.Tag("backup-job", "critical", "nightly")
//	tr.Tag("health-check", "monitoring")
//
//	critical := tr.ByTag("critical")
//	all := tr.AllTags()
package tag
