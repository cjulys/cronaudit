package explain_test

import (
	"testing"

	"github.com/cronaudit/internal/explain"
)

var benchExprs = []string{
	"* * * * *",
	"*/15 * * * *",
	"0 9-17 * * 1-5",
	"30 6 1,15 * *",
	"0 0 1 1 *",
	"*/5 8,12,18 * * 0,6",
}

func BenchmarkExplain(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, expr := range benchExprs {
			_, err := explain.Explain(expr)
			if err != nil {
				b.Fatalf("unexpected error for %q: %v", expr, err)
			}
		}
	}
}

func BenchmarkExplain_Complex(b *testing.B) {
	expr := "0,15,30,45 9-17 1-15 1,3,5,7,9,11 1-5"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := explain.Explain(expr)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
