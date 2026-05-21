package explain

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cronaudit/internal/parser"
)

// Field names in standard cron order.
var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// monthNames maps 1-based month numbers to names.
var monthNames = []string{"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December"}

// dowNames maps 0-based weekday numbers to names.
var dowNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// Explanation holds a human-readable breakdown of a cron expression.
type Explanation struct {
	Expression string
	Fields     []FieldExplanation
	Summary    string
}

// FieldExplanation describes a single cron field.
type FieldExplanation struct {
	Name  string
	Raw   string
	Human string
}

// Explain parses expr and returns a structured Explanation.
// Returns an error if the expression is invalid.
func Explain(expr string) (*Explanation, error) {
	fields, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("explain: %w", err)
	}
	if err := parser.Validate(fields); err != nil {
		return nil, fmt.Errorf("explain: %w", err)
	}

	fe := make([]FieldExplanation, len(fields))
	for i, f := range fields {
		fe[i] = FieldExplanation{
			Name:  fieldNames[i],
			Raw:   f,
			Human: explainField(f, i),
		}
	}

	return &Explanation{
		Expression: expr,
		Fields:     fe,
		Summary:    buildSummary(fe),
	}, nil
}

// explainField converts a single cron field token to plain English.
func explainField(field string, idx int) string {
	if field == "*" {
		return "every " + fieldNames[idx]
	}
	if strings.HasPrefix(field, "*/") {
		step := field[2:]
		return fmt.Sprintf("every %s %s(s)", step, fieldNames[idx])
	}
	if strings.Contains(field, "-") {
		parts := strings.SplitN(field, "-", 2)
		return fmt.Sprintf("%s through %s", labelValue(parts[0], idx), labelValue(parts[1], idx))
	}
	if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		labeled := make([]string, len(parts))
		for i, p := range parts {
			labeled[i] = labelValue(p, idx)
		}
		return strings.Join(labeled, ", ")
	}
	return "at " + labelValue(field, idx)
}

// labelValue optionally maps a numeric string to a name for month/dow fields.
func labelValue(v string, idx int) string {
	n, err := strconv.Atoi(v)
	if err != nil {
		return v
	}
	switch idx {
	case 3:
		if n >= 1 && n <= 12 {
			return monthNames[n]
		}
	case 4:
		if n >= 0 && n <= 6 {
			return dowNames[n]
		}
	}
	return v
}

// buildSummary joins field explanations into a single sentence.
func buildSummary(fields []FieldExplanation) string {
	parts := make([]string, len(fields))
	for i, f := range fields {
		parts[i] = f.Human
	}
	return strings.Join(parts, ", ")
}
