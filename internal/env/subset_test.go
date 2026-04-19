package env

import (
	"testing"
)

func TestSubsetSnapshot_ReturnsRequestedKeys(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2", "C": "3"}
	out, err := SubsetSnapshot(src, SubsetOptions{Keys: []string{"A", "C"}, IgnoreMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Errorf("unexpected values: %v", out)
	}
}

func TestSubsetSnapshot_IgnoreMissing(t *testing.T) {
	src := Snapshot{"A": "1"}
	out, err := SubsetSnapshot(src, SubsetOptions{Keys: []string{"A", "MISSING"}, IgnoreMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MISSING"]; ok {
		t.Error("missing key should not appear in output")
	}
}

func TestSubsetSnapshot_ErrorOnMissing(t *testing.T) {
	src := Snapshot{"A": "1"}
	_, err := SubsetSnapshot(src, SubsetOptions{Keys: []string{"A", "MISSING"}, IgnoreMissing: false})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	me, ok := err.(*MissingKeyError)
	if !ok {
		t.Fatalf("expected *MissingKeyError, got %T", err)
	}
	if me.Key != "MISSING" {
		t.Errorf("expected key MISSING, got %q", me.Key)
	}
}

func TestSubsetSnapshot_DoesNotMutateSource(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2"}
	out, _ := SubsetSnapshot(src, SubsetOptions{Keys: []string{"A"}, IgnoreMissing: true})
	out["A"] = "changed"
	if src["A"] != "1" {
		t.Error("source snapshot was mutated")
	}
}

func TestSubsetSnapshot_EmptyKeys(t *testing.T) {
	src := Snapshot{"A": "1"}
	out, err := SubsetSnapshot(src, SubsetOptions{Keys: nil, IgnoreMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty snapshot, got %v", out)
	}
}
