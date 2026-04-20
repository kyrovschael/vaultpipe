package env

import (
	"testing"
)

func TestCloneSnapshot_AllKeys(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2", "C": "3"}
	dst := CloneSnapshot(src, DefaultCloneOptions())

	if len(dst) != len(src) {
		t.Fatalf("expected %d keys, got %d", len(src), len(dst))
	}
	for k, v := range src {
		if dst[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, dst[k])
		}
	}
}

func TestCloneSnapshot_SubsetKeys(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2", "C": "3"}
	opts := CloneOptions{Keys: []string{"A", "C"}}
	dst := CloneSnapshot(src, opts)

	if len(dst) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(dst))
	}
	if dst["A"] != "1" || dst["C"] != "3" {
		t.Errorf("unexpected values: %v", dst)
	}
	if _, ok := dst["B"]; ok {
		t.Error("key B should not be present")
	}
}

func TestCloneSnapshot_MissingKeySkipped(t *testing.T) {
	src := Snapshot{"A": "1"}
	opts := CloneOptions{Keys: []string{"A", "MISSING"}}
	dst := CloneSnapshot(src, opts)

	if len(dst) != 1 {
		t.Fatalf("expected 1 key, got %d", len(dst))
	}
}

func TestCloneSnapshot_DoesNotMutateOriginal(t *testing.T) {
	src := Snapshot{"X": "original"}
	dst := CloneSnapshot(src, DefaultCloneOptions())
	dst["X"] = "mutated"
	dst["NEW"] = "added"

	if src["X"] != "original" {
		t.Errorf("src mutated: got %q", src["X"])
	}
	if _, ok := src["NEW"]; ok {
		t.Error("NEW key leaked into src")
	}
}

func TestCloneSnapshot_DeepCopy(t *testing.T) {
	src := Snapshot{"K": "value"}
	opts := CloneOptions{Deep: true}
	dst := CloneSnapshot(src, opts)

	if dst["K"] != "value" {
		t.Errorf("expected %q, got %q", "value", dst["K"])
	}
	// Mutate src value; dst must be unaffected (strings are immutable in Go,
	// but we verify the map entry itself is independent).
	src["K"] = "changed"
	if dst["K"] != "value" {
		t.Errorf("deep clone was not independent: got %q", dst["K"])
	}
}

func TestCloneSnapshot_Empty(t *testing.T) {
	dst := CloneSnapshot(Snapshot{}, DefaultCloneOptions())
	if len(dst) != 0 {
		t.Errorf("expected empty snapshot, got %d keys", len(dst))
	}
}
