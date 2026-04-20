package env

import (
	"strings"
	"testing"
)

func TestTruncateSnapshot_ShortValuesUnchanged(t *testing.T) {
	s := Snapshot{"KEY": "hello", "OTHER": "world"}
	out := TruncateSnapshot(s, TruncateOptions{MaxLen: 10, Suffix: "..."})
	if out["KEY"] != "hello" {
		t.Fatalf("expected 'hello', got %q", out["KEY"])
	}
}

func TestTruncateSnapshot_LongValueTruncated(t *testing.T) {
	long := strings.Repeat("a", 300)
	s := Snapshot{"BIG": long}
	out := TruncateSnapshot(s, TruncateOptions{MaxLen: 10, Suffix: "..."})
	if got := out["BIG"]; got != "aaaaaaaaaa..." {
		t.Fatalf("unexpected value: %q", got)
	}
}

func TestTruncateSnapshot_NoSuffix(t *testing.T) {
	long := strings.Repeat("b", 50)
	s := Snapshot{"X": long}
	out := TruncateSnapshot(s, TruncateOptions{MaxLen: 5, Suffix: ""})
	if got := out["X"]; got != "bbbbb" {
		t.Fatalf("expected 5 b's, got %q", got)
	}
}

func TestTruncateSnapshot_RestrictedKeys(t *testing.T) {
	long := strings.Repeat("c", 100)
	s := Snapshot{"A": long, "B": long}
	out := TruncateSnapshot(s, TruncateOptions{MaxLen: 10, Suffix: "...", Keys: []string{"A"}})
	if len(out["A"]) > 13 { // 10 chars + 3 suffix
		t.Fatalf("A should be truncated, got len %d", len(out["A"]))
	}
	if out["B"] != long {
		t.Fatal("B should be unchanged")
	}
}

func TestTruncateSnapshot_DoesNotMutateOriginal(t *testing.T) {
	long := strings.Repeat("d", 200)
	s := Snapshot{"K": long}
	_ = TruncateSnapshot(s, DefaultTruncateOptions())
	if s["K"] != long {
		t.Fatal("original snapshot was mutated")
	}
}

func TestTruncateSnapshot_ZeroMaxLenUsesDefault(t *testing.T) {
	long := strings.Repeat("e", 300)
	s := Snapshot{"Z": long}
	out := TruncateSnapshot(s, TruncateOptions{MaxLen: 0, Suffix: ""})
	def := DefaultTruncateOptions()
	if len(out["Z"]) > def.MaxLen {
		t.Fatalf("expected at most %d bytes, got %d", def.MaxLen, len(out["Z"]))
	}
}

func TestDefaultTruncateOptions_Sensible(t *testing.T) {
	opts := DefaultTruncateOptions()
	if opts.MaxLen <= 0 {
		t.Fatal("default MaxLen should be positive")
	}
	if opts.Suffix == "" {
		t.Fatal("default Suffix should be non-empty")
	}
}
