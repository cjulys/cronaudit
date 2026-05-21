package timeutil_test

import (
	"testing"
	"time"

	"github.com/cronaudit/cronaudit/internal/timeutil"
)

func TestHumanizeDuration_Seconds(t *testing.T) {
	now := time.Now()
	future := now.Add(30 * time.Second)
	got := timeutil.HumanizeDuration(now, future)
	if got != "in 30 seconds" {
		t.Errorf("expected 'in 30 seconds', got %q", got)
	}
}

func TestHumanizeDuration_OneSecond(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Second)
	got := timeutil.HumanizeDuration(now, future)
	if got != "in 1 second" {
		t.Errorf("expected 'in 1 second', got %q", got)
	}
}

func TestHumanizeDuration_Minutes(t *testing.T) {
	now := time.Now()
	future := now.Add(5 * time.Minute)
	got := timeutil.HumanizeDuration(now, future)
	if got != "in 5 minutes" {
		t.Errorf("expected 'in 5 minutes', got %q", got)
	}
}

func TestHumanizeDuration_Hours(t *testing.T) {
	now := time.Now()
	future := now.Add(3 * time.Hour)
	got := timeutil.HumanizeDuration(now, future)
	if got != "in 3 hours" {
		t.Errorf("expected 'in 3 hours', got %q", got)
	}
}

func TestHumanizeDuration_Days(t *testing.T) {
	now := time.Now()
	future := now.Add(48 * time.Hour)
	got := timeutil.HumanizeDuration(now, future)
	if got != "in 2 days" {
		t.Errorf("expected 'in 2 days', got %q", got)
	}
}

func TestHumanizeDuration_Now(t *testing.T) {
	now := time.Now()
	got := timeutil.HumanizeDuration(now, now)
	if got != "now" {
		t.Errorf("expected 'now', got %q", got)
	}
}

func TestHumanizeDuration_Past(t *testing.T) {
	now := time.Now()
	past := now.Add(-10 * time.Minute)
	got := timeutil.HumanizeDuration(now, past)
	if got != "now" {
		t.Errorf("expected 'now' for past time, got %q", got)
	}
}

func TestFormatTime(t *testing.T) {
	ts := time.Date(2024, 6, 15, 9, 5, 0, 0, time.UTC)
	got := timeutil.FormatTime(ts)
	want := "2024-06-15 09:05:00 UTC"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
