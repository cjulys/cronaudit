package explain

import (
	"fmt"
	"io"
	"strings"
)

// Fprint writes a formatted explanation to w.
// Output includes the original expression, a per-field breakdown, and a summary.
func Fprint(w io.Writer, ex *Explanation) error {
	if _, err := fmt.Fprintf(w, "Expression : %s\n", ex.Expression); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Summary    : %s\n", ex.Summary); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, strings.Repeat("-", 40)); err != nil {
		return err
	}
	for _, f := range ex.Fields {
		if _, err := fmt.Fprintf(w, "  %-16s %s  (%s)\n", f.Name+":", f.Raw, f.Human); err != nil {
			return err
		}
	}
	return nil
}

// Sprint returns the formatted explanation as a string.
func Sprint(ex *Explanation) string {
	var sb strings.Builder
	_ = Fprint(&sb, ex)
	return sb.String()
}
