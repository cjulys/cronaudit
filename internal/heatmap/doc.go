// Package heatmap produces a day-of-week × hour-of-day frequency grid
// that visualises when cron jobs are scheduled to run.
//
// Usage:
//
//	res := heatmap.Build(report)
//	fmt.Println(heatmap.Sprint(res))
//
// The grid cells are rendered as ASCII intensity characters:
//
//	'#'  ≥ 75 % of peak
//	'O'  ≥ 50 %
//	'o'  ≥ 25 %
//	'-'  > 0 %
//	'.'  no runs
package heatmap
