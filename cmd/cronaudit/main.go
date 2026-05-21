// Command cronaudit parses and validates cron expressions from one or
// more sources and prints a unified schedule report.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/example/cronaudit/internal/formatter"
	"github.com/example/cronaudit/internal/schedule"
	"github.com/example/cronaudit/internal/source"
)

func main() {
	format := flag.String("format", "text", "output format: text or json")
	n := flag.Int("n", 5, "number of upcoming runs to compute per entry")
	inline := flag.String("expr", "", "comma-separated cron expressions to evaluate inline")
	flag.Parse()

	var entries []schedule.Entry

	// Load expressions from files provided as positional arguments.
	for _, path := range flag.Args() {
		loaded, err := source.FromFile(path)
		if err != nil {
			log.Fatalf("cronaudit: reading %s: %v", path, err)
		}
		entries = append(entries, loaded...)
	}

	// Load inline expressions supplied via -expr flag.
	if *inline != "" {
		parts := strings.Split(*inline, ",")
		loaded, err := source.FromStrings(parts, "inline")
		if err != nil {
			log.Fatalf("cronaudit: parsing inline expressions: %v", err)
		}
		entries = append(entries, loaded...)
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "cronaudit: no expressions provided (pass files or use -expr)")
		flag.Usage()
		os.Exit(1)
	}

	report := schedule.NewReport(entries, *n)

	f, err := formatter.New(*format)
	if err != nil {
		log.Fatalf("cronaudit: %v", err)
	}

	if err := f.Write(os.Stdout, report); err != nil {
		log.Fatalf("cronaudit: writing output: %v", err)
	}
}
