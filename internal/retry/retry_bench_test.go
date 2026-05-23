package retry_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/retry"
	"github.com/cronaudit/internal/schedule"
)

func largeSampleReport(n int) schedule.Report {
	exprs := []string{"* * * * *", "0 * * * *", "0 9 * * 1", "*/15 * * * *", "0 0 1 * *"}
	var entries []schedule.Entry
	for i := 0; i < n; i++ {
		expr := exprs[i%len(exprs)]
		e, err := schedule.NewEntry(expr, schedule.Options{
			Label:  fmt.Sprintf("bench-job-%d", i),
			Origin: "bench",
			N:      5,
		})
		if err == nil {
			entries = append(entries, e)
		}
	}
	return schedule.Report{Entries: entries}
}

func BenchmarkAnalyze_Fixed(b *testing.B) {
	r := largeSampleReport(200)
	cfg := retry.Config{MaxRetries: 3, Window: 10 * time.Minute, Strategy: retry.Fixed}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		retry.Analyze(r, cfg)
	}
}

func BenchmarkAnalyze_Exponential(b *testing.B) {
	r := largeSampleReport(200)
	cfg := retry.Config{MaxRetries: 5, Window: 5 * time.Minute, Strategy: retry.Exponential}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		retry.Analyze(r, cfg)
	}
}

func BenchmarkSprint(b *testing.B) {
	r := largeSampleReport(100)
	cfg := retry.DefaultConfig()
	rep := retry.Analyze(r, cfg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = retry.Sprint(rep)
	}
}

// keep strings import used in test file
var _ = strings.Contains
