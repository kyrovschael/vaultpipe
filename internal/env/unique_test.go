package env

import (
	"testing"
)

func TestUniqueSnapshot_DistinctValues(t *testing.T) {
	s := Snapshot{"A": "hello", "B": "world", "C": "foo"}
	out, err := UniqueSnapshot(s, DefaultUniqueOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestUniqueSnapshot_RemovesDuplicateValues(t *testing.T) {
	s := Snapshot{"A": "same", "B": "same", "C": "other"}
	out, err := UniqueSnapshot(s, DefaultUniqueOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Two keys share value "same"; only one should survive.
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(out), out)
	}
	if out["C"] != "other" {
		t.Errorf("expected C=other, got %q", out["C"])
	}
}

func TestUniqueSnapshot_CaseFold(t *testing.T) {
	s := Snapshot{"A": "Hello", "B": "hello", "C": "WORLD"}
	opts := DefaultUniqueOptions()
	opts.CaseFold = true
	out, err := UniqueSnapshot(s, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "Hello" and "hello" are duplicates under case-fold; only one survives.
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(out), out)
	}
}

func TestUniqueSnapshot_RestrictedKeys(t *testing.T) {
	s := Snapshot{"A": "dup", "B": "dup", "C": "dup"}
	opts := DefaultUniqueOptions()
	opts.Keys = []string{"A", "B"} // only dedupe within A and B; C always passes
	out, err := UniqueSnapshot(s, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// A or B survives (one removed), C always kept.
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(out), out)
	}
	if out["C"] != "dup" {
		t.Errorf("C should always be present, got %v", out)
	}
}

func TestUniqueSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{"X": "v", "Y": "v"}
	orig := Snapshot{"X": "v", "Y": "v"}
	_, _ = UniqueSnapshot(s, DefaultUniqueOptions())
	for k, v := range orig {
		if s[k] != v {
			t.Errorf("original mutated at key %q", k)
		}
	}
}

func TestUniqueSnapshot_EmptySnapshot(t *testing.T) {
	out, err := UniqueSnapshot(Snapshot{}, DefaultUniqueOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty snapshot, got %v", out)
	}
}
