package env

import (
	"testing"
)

func TestProtectSnapshot_ReturnsClone(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2"}
	out, err := ProtectSnapshot(src, []string{"A"}, DefaultProtectOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("expected clone of src, got %v", out)
	}
	// Mutating original must not affect output.
	src["A"] = "changed"
	if out["A"] != "1" {
		t.Fatal("ProtectSnapshot returned a reference, not a clone")
	}
}

func TestProtectSnapshot_StrictMissingKey(t *testing.T) {
	src := Snapshot{"A": "1"}
	opts := DefaultProtectOptions()
	opts.Strict = true
	_, err := ProtectSnapshot(src, []string{"MISSING"}, opts)
	if err == nil {
		t.Fatal("expected error for missing protected key in strict mode")
	}
}

func TestProtectSnapshot_NonStrictMissingKey(t *testing.T) {
	src := Snapshot{"A": "1"}
	_, err := ProtectSnapshot(src, []string{"MISSING"}, DefaultProtectOptions())
	if err != nil {
		t.Fatalf("expected no error in non-strict mode, got %v", err)
	}
}

func TestApplyProtected_BlocksProtectedKeys(t *testing.T) {
	base := Snapshot{"A": "original", "B": "base"}
	patch := Snapshot{"A": "overwrite", "B": "patched"}
	out := ApplyProtected(base, patch, []string{"A"})
	if out["A"] != "original" {
		t.Errorf("protected key A should not be overwritten, got %q", out["A"])
	}
	if out["B"] != "patched" {
		t.Errorf("unprotected key B should be patched, got %q", out["B"])
	}
}

func TestApplyProtected_DoesNotMutateBase(t *testing.T) {
	base := Snapshot{"X": "x"}
	patch := Snapshot{"Y": "y"}
	out := ApplyProtected(base, patch, nil)
	if _, ok := base["Y"]; ok {
		t.Fatal("ApplyProtected mutated base snapshot")
	}
	if out["Y"] != "y" {
		t.Fatal("expected Y in output")
	}
}

func TestApplyProtected_EmptyProtectedList(t *testing.T) {
	base := Snapshot{"A": "1"}
	patch := Snapshot{"A": "2", "B": "3"}
	out := ApplyProtected(base, patch, []string{})
	if out["A"] != "2" {
		t.Errorf("expected A to be overwritten without protection, got %q", out["A"])
	}
}
