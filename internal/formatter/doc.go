// Package formatter provides output formatters for cron audit schedule reports.
//
// Supported formats:
//
//   - text: human-readable tabular output suitable for terminal display
//   - json: structured JSON output suitable for machine consumption
//
// Usage:
//
//	f, err := formatter.New(formatter.FormatText)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := f.Write(os.Stdout, report); err != nil {
//		log.Fatal(err)
//	}
package formatter
