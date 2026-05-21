// Package parser provides functionality for parsing and validating
// standard 5-field cron expressions of the form:
//
//	<minute> <hour> <day-of-month> <month> <day-of-week>
//
// Each field supports the following syntax:
//
//	*        — wildcard (any value)
//	*/n      — step value (every n units)
//	 a-b     — inclusive range
//	a,b,c   — comma-separated list of values
//
// Example usage:
//
//	expr, err := parser.Parse("0 12 * * 1-5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if errs := parser.Validate(expr); len(errs) > 0 {
//		for _, e := range errs {
//			fmt.Println("validation error:", e)
//		}
//	}
package parser
