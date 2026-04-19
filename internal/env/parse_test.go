package env

import (
	"testing"
)

func TestParseSlice_Basic(t *testing.T) {
	snap, err := ParseSlice([]string{"FOO=bar", "BAZ=qux"}, DefaultParseOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap["FOO"] != "bar" || snap["BAZ"] != "qux" {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestParseSlice_TrimSpace(t *testing.T) {
	opts := DefaultParseOptions()
	snap, err := ParseSlice([]string{" FOO = bar "}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap["FOO"] != "bar" {
		t.Fatalf("expected trimmed value, got %q", snap["FOO"])
	}
}

func TestParseSlice_MissingEquals_Error(t *testing.T) {
	opts := DefaultParseOptions()
	_, err := ParseSlice([]string{"NOEQUALS"}, opts)
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseSlice_MissingEquals_SkipInvalid(t *testing.T) {
	opts := DefaultParseOptions()
	opts.SkipInvalid = true
	snap, err := ParseSlice([]string{"NOEQUALS", "GOOD=val"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := snap["NOEQUALS"]; ok {
		t.Fatal("invalid entry should have been skipped")
	}
	if snap["GOOD"] != "val" {
		t.Fatalf("expected GOOD=val, got %v", snap)
	}
}

func TestParseMap_Basic(t *testing.T) {
	snap, err := ParseMap(map[string]string{"KEY": "value"}, DefaultParseOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap["KEY"] != "value" {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestParseMap_InvalidKey_SkipInvalid(t *testing.T) {
	opts := DefaultParseOptions()
	opts.SkipInvalid = true
	snap, err := ParseMap(map[string]string{"1INVALID": "v", "VALID": "ok"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := snap["1INVALID"]; ok {
		t.Fatal("invalid key should have been skipped")
	}
	if snap["VALID"] != "ok" {
		t.Fatalf("expected VALID=ok, got %v", snap)
	}
}

func TestParseMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"A": "1"}
	snap, _ := ParseMap(input, DefaultParseOptions())
	snap["A"] = "mutated"
	if input["A"] != "1" {
		t.Fatal("original map was mutated")
	}
}
