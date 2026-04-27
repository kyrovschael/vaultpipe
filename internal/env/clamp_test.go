package env

import (
	"strings"
	"testing"
)

func TestClampSnapshot_TruncatesLongValues(t *testing.T) {
	src := Snapshot{"KEY": "hello world"}
	opts := DefaultClampOptions()
	opts.MaxLen = 5
	out := ClampSnapshot(src, opts)
	if got := out["KEY"]; got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestClampSnapshot_PadsShortValues(t *testing.T) {
	src := Snapshot{"KEY": "hi"}
	opts := DefaultClampOptions()
	opts.MinLen = 6
	out := ClampSnapshot(src, opts)
	if got := out["KEY"]; got != "hi    " {
		t.Fatalf("expected %q, got %q", "hi    ", got)
	}
}

func TestClampSnapshot_CustomPadChar(t *testing.T) {
	src := Snapshot{"KEY": "x"}
	opts := DefaultClampOptions()
	opts.MinLen = 4
	opts.PadChar = '-'
	out := ClampSnapshot(src, opts)
	if got := out["KEY"]; got != "x---" {
		t.Fatalf("expected %q, got %q", "x---", got)
	}
}

func TestClampSnapshot_ShortValueUnchanged(t *testing.T) {
	src := Snapshot{"KEY": "ok"}
	opts := DefaultClampOptions()
	opts.MaxLen = 10
	out := ClampSnapshot(src, opts)
	if got := out["KEY"]; got != "ok" {
		t.Fatalf("expected %q, got %q", "ok", got)
	}
}

func TestClampSnapshot_RestrictedKeys(t *testing.T) {
	src := Snapshot{"A": "toolong", "B": "toolong"}
	opts := DefaultClampOptions()
	opts.MaxLen = 3
	opts.Keys = []string{"A"}
	out := ClampSnapshot(src, opts)
	if got := out["A"]; got != "too" {
		t.Fatalf("A: expected %q, got %q", "too", got)
	}
	if got := out["B"]; got != "toolong" {
		t.Fatalf("B: expected %q, got %q", "toolong", got)
	}
}

func TestClampSnapshot_DoesNotMutateOriginal(t *testing.T) {
	src := Snapshot{"KEY": strings.Repeat("a", 20)}
	opts := DefaultClampOptions()
	opts.MaxLen = 5
	_ = ClampSnapshot(src, opts)
	if got := src["KEY"]; len(got) != 20 {
		t.Fatalf("original mutated: len=%d", len(got))
	}
}

func TestClampSnapshot_ZeroLimits_Passthrough(t *testing.T) {
	src := Snapshot{"KEY": "unchanged"}
	out := ClampSnapshot(src, DefaultClampOptions())
	if got := out["KEY"]; got != "unchanged" {
		t.Fatalf("expected %q, got %q", "unchanged", got)
	}
}
