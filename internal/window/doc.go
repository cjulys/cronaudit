// Package window provides time-window querying for cron schedule reports.
//
// Given a [schedule.Report] and a [from, to) time interval, Query returns
// every valid entry whose pre-computed next runs fall within that window,
// together with the matching timestamps.
//
// Typical usage:
//
//	from := time.Now()
//	to   := from.Add(24 * time.Hour)
//	res  := window.Query(report, from, to)
//	fmt.Print(window.Sprint(res))
package window
