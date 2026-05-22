package overlap

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// Fprint writes a human-readable overlap report to w.
func Fprint(w io.Writer, result Result) error {
	if !result.HasConflicts() {
		_, err := fmt.Fprintln(w, "No scheduling conflicts detected.")
		return err
	}

	_, err := fmt.Fprintf(w, "Found %d scheduling conflict(s):\n\n", len(result.Conflicts))
	if err != nil {
		return err
	}

	for i, c := range result.Conflicts {
		_, err = fmt.Fprintf(w,
			"  [%d] %q and %q overlap at %s\n",
			i+1,
			c.A.Label,
			c.B.Label,
			c.OverlapAt.Format(time.RFC3339),
		)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w,
			"      %s: %s\n      %s: %s\n\n",
			c.A.Label, c.A.Expression,
			c.B.Label, c.B.Expression,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Sprint returns the overlap report as a string.
func Sprint(result Result) string {
	var buf bytes.Buffer
	_ = Fprint(&buf, result)
	return buf.String()
}
