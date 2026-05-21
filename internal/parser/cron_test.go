package parser

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	expr := "0 12 * * 1-5"
	c, err := Parse(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Raw != expr {
		t.Errorf("expected Raw %q, got %q", expr, c.Raw)
	}
	if len(c.Fields) != 5 {
		t.Errorf("expected 5 fields, got %d", len(c.Fields))
	}
}

func TestParse_Empty(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected error for empty expression")
	}
}

func TestParse_WrongFieldCount(t *testing.T) {
	_, err := Parse("0 12 *")
	if err == nil {
		t.Fatal("expected error for wrong field count")
	}
}

func TestValidate_AllWildcards(t *testing.T) {
	c, _ := Parse("* * * * *")
	errs := Validate(c)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_ValidStep(t *testing.T) {
	c, _ := Parse("*/15 * * * *")
	errs := Validate(c)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_InvalidMinute(t *testing.T) {
	c, _ := Parse("99 12 * * *")
	errs := Validate(c)
	if len(errs) == 0 {
		t.Error("expected validation error for minute=99")
	}
}

func TestValidate_ValidRange(t *testing.T) {
	c, _ := Parse("0 9-17 * * 1-5")
	errs := Validate(c)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	c, _ := Parse("0 25-30 * * *")
	errs := Validate(c)
	if len(errs) == 0 {
		t.Error("expected validation error for hour range 25-30")
	}
}

func TestValidate_ValidList(t *testing.T) {
	c, _ := Parse("0,15,30,45 * * * *")
	errs := Validate(c)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_InvalidList(t *testing.T) {
	c, _ := Parse("0,15,99 * * * *")
	errs := Validate(c)
	if len(errs) == 0 {
		t.Error("expected validation error for minute list containing 99")
	}
}
