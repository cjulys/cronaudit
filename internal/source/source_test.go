package source_test

import (
	"strings"
	"testing"

	"github.com/cronaudit/cronaudit/internal/source"
)

const sampleInput = `# daily backup
0 2 * * *	 daily-backup
# weekly report
0 9 * * 1  weekly-report
*/5 * * * *

# blank lines and comments above should be skipped
`

func TestFromReader_ParsesExpressions(t *testing.T) {
	r := strings.NewReader(sampleInput)
	entries, err := source.FromReader(r, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestFromReader_Labels(t *testing.T) {
	r := strings.NewReader(sampleInput)
	entries, err := source.FromReader(r, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Label != "daily-backup" {
		t.Errorf("expected label %q, got %q", "daily-backup", entries[0].Label)
	}
	if entries[1].Label != "weekly-report" {
		t.Errorf("expected label %q, got %q", "weekly-report", entries[1].Label)
	}
	if entries[2].Label != "" {
		t.Errorf("expected empty label, got %q", entries[2].Label)
	}
}

func TestFromReader_Origin(t *testing.T) {
	r := strings.NewReader("* * * * *")
	entries, err := source.FromReader(r, "stdin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Origin != "stdin" {
		t.Errorf("expected origin %q, got %q", "stdin", entries[0].Origin)
	}
}

func TestFromStrings_FiltersBlankAndComments(t *testing.T) {
	input := []string{"* * * * *", "", "# comment", "0 0 * * *"}
	entries := source.FromStrings(input, "args")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Expression != "* * * * *" {
		t.Errorf("unexpected expression: %q", entries[0].Expression)
	}
	if entries[1].Expression != "0 0 * * *" {
		t.Errorf("unexpected expression: %q", entries[1].Expression)
	}
}

func TestFromReader_EmptyInput(t *testing.T) {
	r := strings.NewReader("")
	entries, err := source.FromReader(r, "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}
