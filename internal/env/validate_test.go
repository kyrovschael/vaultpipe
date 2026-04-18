package env

import (
	"testing"
)

func TestValidateKey_Valid(t *testing.T) {
	valid := []string{"HOME", "_VAR", "MY_VAR_1", "a", "z9"}
	for _, k := range valid {
		if err := ValidateKey(k); err != nil {
			t.Errorf("expected %q to be valid, got: %v", k, err)
		}
	}
}

func TestValidateKey_Invalid(t *testing.T) {
	invalid := []string{"", "1START", "has-dash", "has space", "has=eq"}
	for _, k := range invalid {
		if err := ValidateKey(k); err == nil {
			t.Errorf("expected %q to be invalid, got nil", k)
		}
	}
}

func TestValidateSnapshot_AllValid(t *testing.T) {
	s := Snapshot{"HOME": "/home/user", "PATH": "/usr/bin", "_PRIVATE": "val"}
	if err := ValidateSnapshot(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSnapshot_SomeInvalid(t *testing.T) {
	s := Snapshot{
		"VALID_KEY": "ok",
		"bad-key":   "nope",
		"1INVALID":  "nope",
	}
	err := ValidateSnapshot(s)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.InvalidKeys) != 2 {
		t.Errorf("expected 2 invalid keys, got %d: %v", len(ve.InvalidKeys), ve.InvalidKeys)
	}
}

func TestValidateSnapshot_Empty(t *testing.T) {
	if err := ValidateSnapshot(Snapshot{}); err != nil {
		t.Fatalf("unexpected error on empty snapshot: %v", err)
	}
}
