package env

import (
	"testing"
)

func TestPatchSnapshot_OverwritesBase(t *testing.T) {
	base := Snapshot{"A": "1", "B": "2"}
	patch := Snapshot{"B": "99", "C": "3"}
	out := PatchSnapshot(base, patch, DefaultPatchOptions())
	if out["B"] != "99" {
		t.Fatalf("expected B=99, got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Fatalf("expected C=3, got %s", out["C"])
	}
	if out["A"] != "1" {
		t.Fatalf("expected A=1, got %s", out["A"])
	}
}

func TestPatchSnapshot_SkipEmpty(t *testing.T) {
	base := Snapshot{"A": "original"}
	patch := Snapshot{"A": ""}
	opts := DefaultPatchOptions()
	opts.SkipEmpty = true
	out := PatchSnapshot(base, patch, opts)
	if out["A"] != "original" {
		t.Fatalf("expected A=original, got %s", out["A"])
	}
}

func TestPatchSnapshot_RemovesKeys(t *testing.T) {
	base := Snapshot{"A": "1", "B": "2", "C": "3"}
	opts := DefaultPatchOptions()
	opts.Remove = []string{"B", "C"}
	out := PatchSnapshot(base, Snapshot{}, opts)
	if _, ok := out["B"]; ok {
		t.Fatal("expected B to be removed")
	}
	if _, ok := out["C"]; ok {
		t.Fatal("expected C to be removed")
	}
	if out["A"] != "1" {
		t.Fatalf("expected A=1, got %s", out["A"])
	}
}

func TestPatchSnapshot_DoesNotMutateBase(t *testing.T) {
	base := Snapshot{"A": "1"}
	patch := Snapshot{"A": "2", "B": "3"}
	PatchSnapshot(base, patch, DefaultPatchOptions())
	if base["A"] != "1" {
		t.Fatal("base was mutated")
	}
	if _, ok := base["B"]; ok {
		t.Fatal("base gained unexpected key")
	}
}

func TestPatchSnapshot_EmptyPatch(t *testing.T) {
	base := Snapshot{"X": "hello"}
	out := PatchSnapshot(base, Snapshot{}, DefaultPatchOptions())
	if out["X"] != "hello" {
		t.Fatalf("expected X=hello, got %s", out["X"])
	}
}
