package notify

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Fprint writes a human-readable summary of alerts to w.
func Fprint(w io.Writer, alerts []Alert) {
	if len(alerts) == 0 {
		fmt.Fprintln(w, "No alerts.")
		return
	}
	fmt.Fprintf(w, "Alerts (%d):\n", len(alerts))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, a := range alerts {
		fmt.Fprintf(w, "[%s] %s (%s)\n", a.Level, a.Label, a.Expression)
		fmt.Fprintf(w, "       %s\n", a.Message)
	}
}

// Sprint returns the Fprint output as a string.
func Sprint(alerts []Alert) string {
	var buf bytes.Buffer
	Fprint(&buf, alerts)
	return buf.String()
}
