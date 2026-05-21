package filter_test

import (
	"fmt"
	"testing"

	"github.com/example/cronaudit/internal/filter"
	"github.com/example/cronaudit/internal/schedule"
)

func largeSampleReport(n int) *schedule.Report {
	entries := make([]schedule.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = schedule.Entry{
			Label:      fmt.Sprintf("job-%d", i),
			Origin:     []string{"crontab", "systemd", "k8s"}[i%3],
			Expression: fmt.Sprintf("%d * * * *", i%60),
		}
	}
	return &schedule.Report{Entries: entries}
}

func BenchmarkApply_LabelFilter(b *testing.B) {
	report := largeSampleReport(1000)
	opts := filter.Options{Label: "job-5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Apply(report, opts)
	}
}

func BenchmarkApply_OriginFilter(b *testing.B) {
	report := largeSampleReport(1000)
	opts := filter.Options{Origin: "crontab"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Apply(report, opts)
	}
}

func BenchmarkApply_CombinedFilter(b *testing.B) {
	report := largeSampleReport(1000)
	opts := filter.Options{Origin: "systemd", Label: "job-1"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.Apply(report, opts)
	}
}
