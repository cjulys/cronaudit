// Package suggest provides cron expression suggestions and corrections
// for common mistakes or non-portable constructs.
package suggest

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/schedule"
)

// Suggestion represents a recommended alternative for a cron expression.
type Suggestion struct {
	Original    string
	Replacement string
	Reason      string
}

// Analyze inspects a report and returns suggestions for each entry
// that contains known non-portable or fragile patterns.
func Analyze(r *schedule.Report) map[string][]Suggestion {
	result := make(map[string][]Suggestion)
	for _, entry := range r.Entries {
		suggestions := analyzeExpression(entry.Expression)
		if len(suggestions) > 0 {
			result[entry.Label] = suggestions
		}
	}
	return result
}

func analyzeExpression(expr string) []Suggestion {
	var suggestions []Suggestion

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return suggestions
	}

	// Suggest replacing @reboot / named macros (not standard 5-field)
	// Detect */1 which is equivalent to *
	for i, f := range fields {
		if f == "*/1" {
			fields[i] = "*"
			suggestions = append(suggestions, Suggestion{
				Original:    expr,
				Replacement: strings.Join(fields, " "),
				Reason:      fmt.Sprintf("field %d: '*/1' is equivalent to '*'; prefer '*' for clarity", i+1),
			})
			fields[i] = "*/1" // restore for next iteration
		}
	}

	// Detect day-of-month and day-of-week both set (ambiguous across implementations)
	dom := fields[2]
	dow := fields[4]
	if dom != "*" && dow != "*" {
		suggestions = append(suggestions, Suggestion{
			Original:    expr,
			Replacement: "",
			Reason:      "both day-of-month and day-of-week are set; behaviour differs between cron implementations (OR vs AND)",
		})
	}

	// Detect minute wildcard with hour wildcard — runs every minute
	if fields[0] == "*" && fields[1] == "*" {
		suggestions = append(suggestions, Suggestion{
			Original:    expr,
			Replacement: "0 * * * *",
			Reason:      "expression runs every minute; if hourly was intended, use '0 * * * *'",
		})
	}

	return suggestions
}
