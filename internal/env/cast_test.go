package env

import (
	"testing"
)

func TestCastSnapshot_IntNormalisation(t *testing.T) {
	snap := Snapshot{"PORT": "  8080  ", "NAME": "app"}
	out, err := CastSnapshot(snap, []CastRule{{Key: "PORT", Type: CastInt}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", out["PORT"])
	}
	if out["NAME"] != "app" {
		t.Errorf("NAME should be unchanged")
	}
}

func TestCastSnapshot_BoolNormalisation(t *testing.T) {
	snap := Snapshot{"DEBUG": "TRUE"}
	out, err := CastSnapshot(snap, []CastRule{{Key: "DEBUG", Type: CastBool}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DEBUG"] != "true" {
		t.Errorf("expected true, got %q", out["DEBUG"])
	}
}

func TestCastSnapshot_FloatNormalisation(t *testing.T) {
	snap := Snapshot{"RATIO": "3.14000"}
	out, err := CastSnapshot(snap, []CastRule{{Key: "RATIO", Type: CastFloat}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RATIO"] != "3.14" {
		t.Errorf("expected 3.14, got %q", out["RATIO"])
	}
}

func TestCastSnapshot_InvalidInt(t *testing.T) {
	snap := Snapshot{"PORT": "not-a-number"}
	_, err := CastSnapshot(snap, []CastRule{{Key: "PORT", Type: CastInt}})
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestCastSnapshot_MissingKeySkipped(t *testing.T) {
	snap := Snapshot{"A": "1"}
	out, err := CastSnapshot(snap, []CastRule{{Key: "MISSING", Type: CastInt}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MISSING"]; ok {
		t.Error("missing key should not be added")
	}
}

func TestCastSnapshot_DoesNotMutateOriginal(t *testing.T) {
	snap := Snapshot{"PORT": "8080"}
	_, _ = CastSnapshot(snap, []CastRule{{Key: "PORT", Type: CastInt}})
	if snap["PORT"] != "8080" {
		t.Error("original snapshot was mutated")
	}
}

func TestCastSnapshot_StringPassthrough(t *testing.T) {
	snap := Snapshot{"NAME": "vault"}
	out, err := CastSnapshot(snap, []CastRule{{Key: "NAME", Type: CastString}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NAME"] != "vault" {
		t.Errorf("expected vault, got %q", out["NAME"])
	}
}
