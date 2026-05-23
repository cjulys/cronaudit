// Package budget analyses schedule entries against execution-frequency and
// CPU-cost budgets.
//
// Use Analyze to evaluate a schedule.Report against a Config, then use
// Fprint or Sprint to render the results in a human-readable form.
//
// Example:
//
//	cfg := budget.DefaultConfig()
//	cfg.MaxRunsPerDay = 144 // tighter limit
//	results := budget.Analyze(report, cfg)
//	fmt.Print(budget.Sprint(results))
package budget
