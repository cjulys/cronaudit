// Package retry models hypothetical retry windows for cron entries.
//
// Given a [schedule.Report] and a [Config], [Analyze] projects when retries
// would fire if a job fails on its first scheduled run. Two spacing strategies
// are supported:
//
//   - [Fixed]       – retries spaced by a constant window duration.
//   - [Exponential] – each retry doubles the previous window.
//
// Example:
//
//	cfg := retry.DefaultConfig()
//	rep := retry.Analyze(schedReport, cfg)
//	fmt.Print(retry.Sprint(rep))
package retry
