package env

import (
	"testing"
)

func TestSquashSnapshot_BasicMerge(t *testing.T) {
	snap := Snapshot{
		"DB_HOSTS_0": "alpha",
		"DB_HOSTS_1": "beta",
		"DB_HOSTS_2": "gamma",
	}
	out, err := SquashSnapshot(snap, DefaultSquashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := out["DB_HOSTS"]
	if !ok {
		t.Fatal("expected key DB_HOSTS in output")
	}
	// values must be joined; order is determined by suffix sort
	if got != "alpha,beta,gamma" {
		t.Errorf("got %q, want %q", got, "alpha,beta,gamma")
	}
}

func TestSquashSnapshot_CustomSeparator(t *testing.T) {
	snap := Snapshot{
		"ADDRS_0": "10.0.0.1",
		"ADDRS_1": "10.0.0.2",
	}
	opts := DefaultSquashOptions()
	opts.Separator = ";"
	out, err := SquashSnapshot(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["ADDRS"]; got != "10.0.0.1;10.0.0.2" {
		t.Errorf("got %q, want %q", got, "10.0.0.1;10.0.0.2")
	}
}

func TestSquashSnapshot_NoSuffixedKeys_Passthrough(t *testing.T) {
	snap := Snapshot{
		"PLAIN": "value",
		"OTHER": "data",
	}
	out, err := SquashSnapshot(snap, DefaultSquashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PLAIN"] != "value" {
		t.Errorf("expected PLAIN to pass through")
	}
	if out["OTHER"] != "data" {
		t.Errorf("expected OTHER to pass through")
	}
}

func TestSquashSnapshot_DoesNotMutateOriginal(t *testing.T) {
	snap := Snapshot{
		"X_0": "one",
		"X_1": "two",
	}
	origLen := len(snap)
	_, err := SquashSnapshot(snap, DefaultSquashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap) != origLen {
		t.Error("original snapshot was mutated")
	}
	if _, ok := snap["X"]; ok {
		t.Error("squashed key leaked into original")
	}
}

func TestSquashSnapshot_MixedKeys(t *testing.T) {
	snap := Snapshot{
		"HOSTS_0":  "h1",
		"HOSTS_1":  "h2",
		"TIMEOUT":  "30s",
		"RETRIES_0": "3",
	}
	out, err := SquashSnapshot(snap, DefaultSquashOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOSTS"] != "h1,h2" {
		t.Errorf("HOSTS: got %q", out["HOSTS"])
	}
	if out["TIMEOUT"] != "30s" {
		t.Errorf("TIMEOUT: got %q", out["TIMEOUT"])
	}
	if out["RETRIES"] != "3" {
		t.Errorf("RETRIES: got %q", out["RETRIES"])
	}
}
