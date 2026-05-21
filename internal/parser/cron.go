package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Field represents a single cron expression field with its constraints.
type Field struct {
	Name string
	Min  int
	Max  int
}

// StandardFields defines the five standard cron fields.
var StandardFields = []Field{
	{Name: "minute", Min: 0, Max: 59},
	{Name: "hour", Min: 0, Max: 23},
	{Name: "day-of-month", Min: 1, Max: 31},
	{Name: "month", Min: 1, Max: 12},
	{Name: "day-of-week", Min: 0, Max: 6},
}

// CronExpression holds a parsed cron expression.
type CronExpression struct {
	Raw    string
	Fields []string
}

// ParseError describes a validation error for a specific field.
type ParseError struct {
	Field   string
	Value   string
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("field %q value %q: %s", e.Field, e.Value, e.Message)
}

// Parse tokenizes a cron expression string into a CronExpression.
func Parse(expr string) (*CronExpression, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("empty cron expression")
	}
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}
	return &CronExpression{Raw: expr, Fields: parts}, nil
}

// Validate checks each field of the CronExpression against allowed ranges.
func Validate(c *CronExpression) []error {
	var errs []error
	for i, field := range StandardFields {
		if err := validateField(c.Fields[i], field); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func validateField(value string, field Field) error {
	if value == "*" {
		return nil
	}
	// Handle step values e.g. */5
	if strings.HasPrefix(value, "*/") {
		step, err := strconv.Atoi(value[2:])
		if err != nil || step < 1 {
			return &ParseError{Field: field.Name, Value: value, Message: "invalid step value"}
		}
		return nil
	}
	// Handle ranges e.g. 1-5
	if strings.Contains(value, "-") {
		parts := strings.SplitN(value, "-", 2)
		lo, err1 := strconv.Atoi(parts[0])
		hi, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || lo > hi || lo < field.Min || hi > field.Max {
			return &ParseError{Field: field.Name, Value: value, Message: fmt.Sprintf("range must be within %d-%d", field.Min, field.Max)}
		}
		return nil
	}
	// Handle lists e.g. 1,3,5
	for _, part := range strings.Split(value, ",") {
		n, err := strconv.Atoi(part)
		if err != nil || n < field.Min || n > field.Max {
			return &ParseError{Field: field.Name, Value: value, Message: fmt.Sprintf("value must be within %d-%d", field.Min, field.Max)}
		}
	}
	return nil
}
