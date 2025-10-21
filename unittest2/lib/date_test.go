package lib_test

import (
	"testing"
	"unit-test-demo/unittest2/lib"
)

func TestFormatDateLong_Basic(t *testing.T) {
	got, err := lib.FormatDateLong("2025-10-20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "20 October 2025"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatDateLong_LeadingZeros(t *testing.T) {
	got, err := lib.FormatDateLong("2024-01-05")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "5 January 2024"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatDateLong_InvalidInput(t *testing.T) {
	if _, err := lib.FormatDateLong("2024-02-30"); err == nil {
		t.Fatal("expected error for invalid date, got nil")
	}
}