package env

import (
	"testing"
)

func TestCoerceSnapshot_BoolNormalisation(t *testing.T) {
	snap := FromSlice([]string{"ENABLED=1", "VERBOSE=FALSE"})
	rules := []CoerceRule{
		{Key: "ENABLED", Type: CoerceBool},
		{Key: "VERBOSE", Type: CoerceBool},
	}
	out, errs := CoerceSnapshot(snap, rules)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if v, _ := out.Get("ENABLED"); v != "true" {
		t.Errorf("ENABLED: got %q, want \"true\"", v)
	}
	if v, _ := out.Get("VERBOSE"); v != "false" {
		t.Errorf("VERBOSE: got %q, want \"false\"", v)
	}
}

func TestCoerceSnapshot_IntPassthrough(t *testing.T) {
	snap := FromSlice([]string{"PORT=8080"})
	rules := []CoerceRule{{Key: "PORT", Type: CoerceInt}}
	out, errs := CoerceSnapshot(snap, rules)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if v, _ := out.Get("PORT"); v != "8080" {
		t.Errorf("PORT: got %q, want \"8080\"", v)
	}
}

func TestCoerceSnapshot_InvalidInt(t *testing.T) {
	snap := FromSlice([]string{"PORT=abc"})
	rules := []CoerceRule{{Key: "PORT", Type: CoerceInt}}
	_, errs := CoerceSnapshot(snap, rules)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestCoerceSnapshot_MissingKeySkipped(t *testing.T) {
	snap := FromSlice([]string{"FOO=bar"})
	rules := []CoerceRule{{Key: "MISSING", Type: CoerceBool}}
	_, errs := CoerceSnapshot(snap, rules)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestCoerceSnapshot_DoesNotMutateOriginal(t *testing.T) {
	snap := FromSlice([]string{"ENABLED=1"})
	rules := []CoerceRule{{Key: "ENABLED", Type: CoerceBool}}
	CoerceSnapshot(snap, rules)
	if v, _ := snap.Get("ENABLED"); v != "1" {
		t.Errorf("original mutated: got %q, want \"1\"", v)
	}
}
