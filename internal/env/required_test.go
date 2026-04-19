package env

import (
	"errors"
	"testing"
)

func TestRequireKeys_AllPresent(t *testing.T) {
	snap := Snapshot{"A": "1", "B": "2"}
	if err := RequireKeys(snap, []string{"A", "B"}, DefaultRequireOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireKeys_MissingKey(t *testing.T) {
	snap := Snapshot{"A": "1"}
	err := RequireKeys(snap, []string{"A", "B"}, DefaultRequireOptions())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var re *RequiredError
	if !errors.As(err, &re) {
		t.Fatalf("expected *RequiredError, got %T", err)
	}
	if len(re.Missing) != 1 || re.Missing[0] != "B" {
		t.Fatalf("unexpected missing keys: %v", re.Missing)
	}
}

func TestRequireKeys_EmptyValueDisallowed(t *testing.T) {
	snap := Snapshot{"A": ""}
	err := RequireKeys(snap, []string{"A"}, DefaultRequireOptions())
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestRequireKeys_EmptyValueAllowed(t *testing.T) {
	snap := Snapshot{"A": ""}
	opts := RequireOptions{AllowEmpty: true}
	if err := RequireKeys(snap, []string{"A"}, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireKeys_MultipleMissing(t *testing.T) {
	snap := Snapshot{}
	err := RequireKeys(snap, []string{"X", "Y", "Z"}, DefaultRequireOptions())
	var re *RequiredError
	if !errors.As(err, &re) {
		t.Fatalf("expected *RequiredError")
	}
	if len(re.Missing) != 3 {
		t.Fatalf("expected 3 missing, got %d", len(re.Missing))
	}
}

func TestRequireKeys_EmptyKeyList(t *testing.T) {
	snap := Snapshot{}
	if err := RequireKeys(snap, nil, DefaultRequireOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
