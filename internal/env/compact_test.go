package env

import (
	"testing"
)

func TestCompactSnapshot_DropsEmpty(t *testing.T) {
	s := Snapshot{"A": "hello", "B": "", "C": "world"}
	opts := DefaultCompactOptions()
	out := CompactSnapshot(s, opts)

	if _, ok := out["B"]; ok {
		t.Error("expected empty key B to be dropped")
	}
	if out["A"] != "hello" {
		t.Errorf("expected A=hello, got %q", out["A"])
	}
	if out["C"] != "world" {
		t.Errorf("expected C=world, got %q", out["C"])
	}
}

func TestCompactSnapshot_TrimsSpace(t *testing.T) {
	s := Snapshot{"X": "  trimmed  ", "Y": "\t\n"}
	opts := DefaultCompactOptions()
	out := CompactSnapshot(s, opts)

	if out["X"] != "trimmed" {
		t.Errorf("expected X=trimmed, got %q", out["X"])
	}
	// whitespace-only becomes empty after trim, so dropped
	if _, ok := out["Y"]; ok {
		t.Error("expected whitespace-only key Y to be dropped")
	}
}

func TestCompactSnapshot_DropWhitespace_NoTrim(t *testing.T) {
	s := Snapshot{"A": "   ", "B": "keep"}
	opts := CompactOptions{
		TrimSpace:      false,
		DropEmpty:      true,
		DropWhitespace: true,
	}
	out := CompactSnapshot(s, opts)

	if _, ok := out["A"]; ok {
		t.Error("expected whitespace-only key A to be dropped")
	}
	if out["B"] != "keep" {
		t.Errorf("expected B=keep, got %q", out["B"])
	}
}

func TestCompactSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{"K": "  value  "}
	opts := DefaultCompactOptions()
	_ = CompactSnapshot(s, opts)

	if s["K"] != "  value  " {
		t.Error("original snapshot was mutated")
	}
}

func TestCompactSnapshot_EmptySnapshot(t *testing.T) {
	out := CompactSnapshot(Snapshot{}, DefaultCompactOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}

func TestCompactSnapshot_NoDropOptions(t *testing.T) {
	s := Snapshot{"A": "", "B": "  "}
	opts := CompactOptions{
		TrimSpace:      false,
		DropEmpty:      false,
		DropWhitespace: false,
	}
	out := CompactSnapshot(s, opts)

	if len(out) != 2 {
		t.Errorf("expected 2 entries preserved, got %d", len(out))
	}
}
