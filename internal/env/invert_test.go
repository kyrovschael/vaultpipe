package env

import (
	"testing"
)

func TestInvertSnapshot_Basic(t *testing.T) {
	src := Snapshot{"FOO": "bar", "BAZ": "qux"}
	got := InvertSnapshot(src, DefaultInvertOptions())

	if got["bar"] != "FOO" {
		t.Errorf("expected got[bar]=FOO, got %q", got["bar"])
	}
	if got["qux"] != "BAZ" {
		t.Errorf("expected got[qux]=BAZ, got %q", got["qux"])
	}
}

func TestInvertSnapshot_DropsEmptyValue(t *testing.T) {
	src := Snapshot{"EMPTY": "", "KEY": "val"}
	got := InvertSnapshot(src, DefaultInvertOptions())

	if _, ok := got[""); ok {
		t.Error("expected empty-value entry to be dropped")
	}
	if got["val"] != "KEY" {
		t.Errorf("expected got[val]=KEY, got %q", got["val"])
	}
}

func TestInvertSnapshot_SkipDuplicateValues_KeepsSmallestKey(t *testing.T) {
	// Both A_KEY and B_KEY map to the same value "shared".
	// With SkipDuplicateValues=true the smallest key (A_KEY) should win.
	src := Snapshot{"A_KEY": "shared", "B_KEY": "shared"}
	opts := DefaultInvertOptions()
	opts.SkipDuplicateValues = true
	got := InvertSnapshot(src, opts)

	if got["shared"] != "A_KEY" {
		t.Errorf("expected got[shared]=A_KEY, got %q", got["shared"])
	}
}

func TestInvertSnapshot_NoDuplicates_LastWins(t *testing.T) {
	// Unique values — inversion is unambiguous.
	src := Snapshot{"X": "one", "Y": "two"}
	opts := DefaultInvertOptions()
	opts.SkipDuplicateValues = false
	got := InvertSnapshot(src, opts)

	if got["one"] != "X" {
		t.Errorf("expected got[one]=X, got %q", got["one"])
	}
	if got["two"] != "Y" {
		t.Errorf("expected got[two]=Y, got %q", got["two"])
	}
}

func TestInvertSnapshot_DoesNotMutateOriginal(t *testing.T) {
	src := Snapshot{"K": "V"}
	_ = InvertSnapshot(src, DefaultInvertOptions())

	if src["K"] != "V" {
		t.Error("InvertSnapshot mutated the source snapshot")
	}
}

func TestInvertSnapshot_EmptySnapshot(t *testing.T) {
	got := InvertSnapshot(Snapshot{}, DefaultInvertOptions())
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d entries", len(got))
	}
}
