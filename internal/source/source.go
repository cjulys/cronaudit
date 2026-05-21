// Package source provides types and utilities for loading cron expressions
// from various input sources such as files, stdin, and raw strings.
package source

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Entry represents a single cron entry loaded from a source.
type Entry struct {
	// Expression is the raw cron expression string.
	Expression string
	// Label is an optional human-readable name for the entry.
	Label string
	// Origin describes where the entry was loaded from (e.g. filename, "stdin").
	Origin string
}

// FromFile reads cron expressions from a file, one per line.
// Lines beginning with '#' and blank lines are ignored.
// Each line may optionally include a label separated by a tab or multiple spaces.
func FromFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("source: open file %q: %w", path, err)
	}
	defer f.Close()
	return fromReader(f, path)
}

// FromReader reads cron expressions from an io.Reader.
// origin is used to populate Entry.Origin for traceability.
func FromReader(r io.Reader, origin string) ([]Entry, error) {
	return fromReader(r, origin)
}

// FromStrings wraps a slice of raw expression strings into Entries.
// origin is applied to all resulting entries.
func FromStrings(exprs []string, origin string) []Entry {
	entries := make([]Entry, 0, len(exprs))
	for _, e := range exprs {
		e = strings.TrimSpace(e)
		if e == "" || strings.HasPrefix(e, "#") {
			continue
		}
		entries = append(entries, Entry{Expression: e, Origin: origin})
	}
	return entries
}

// fromReader is the shared implementation for file and reader loading.
func fromReader(r io.Reader, origin string) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entry := Entry{Origin: origin}
		// Support optional label after a tab or two+ spaces
		if idx := strings.Index(line, "\t"); idx != -1 {
			entry.Expression = strings.TrimSpace(line[:idx])
			entry.Label = strings.TrimSpace(line[idx+1:])
		} else if idx := strings.Index(line, "  "); idx != -1 {
			entry.Expression = strings.TrimSpace(line[:idx])
			entry.Label = strings.TrimSpace(line[idx+2:])
		} else {
			entry.Expression = line
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("source: read from %q: %w", origin, err)
	}
	return entries, nil
}
