package env

import (
	"testing"
)

func TestSanitizeSnapshot_ValidKeysUnchanged(t *testing.T) {
	src := Snapshot{"VALID_KEY": "value", "_ALSO_VALID": "x"}
	out := SanitizeSnapshot(src, DefaultSanitizeOptions())
	for k := range src {
		if out[k] != src[k] {
			t.Errorf("expected key %q unchanged", k)
		}
	}
}

func TestSanitizeSnapshot_ReplacesInvalidChars(t *testing.T) {
	src := Snapshot{"my-key": "val", "key.name": "v2"}
	out := SanitizeSnapshot(src, DefaultSanitizeOptions())
	if out["my_key"] != "val" {
		t.Errorf("expected my-key -> my_key, got %v", out)
	}
	if out["key_name"] != "v2" {
		t.Errorf("expected key.name -> key_name, got %v", out)
	}
}

func TestSanitizeSnapshot_LeadingDigitReplaced(t *testing.T) {
	src := Snapshot{"1bad": "v"}
	out := SanitizeSnapshot(src, DefaultSanitizeOptions())
	if _, ok := out["1bad"]; ok {
		t.Error("original key should not exist")
	}
	if out["_bad"] != "v" {
		t.Errorf("expected _bad, got %v", out)
	}
}

func TestSanitizeSnapshot_SkipInvalidKeys(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.ReplaceInvalidChars = ""
	opts.SkipInvalidKeys = true
	src := Snapshot{"bad key": "v", "GOOD": "ok"}
	out := SanitizeSnapshot(src, opts)
	if _, ok := out["bad key"]; ok {
		t.Error("bad key should be skipped")
	}
	if out["GOOD"] != "ok" {
		t.Error("GOOD should be retained")
	}
}

func TestSanitizeSnapshot_DoesNotMutateSource(t *testing.T) {
	src := Snapshot{"my-key": "val"}
	SanitizeSnapshot(src, DefaultSanitizeOptions())
	if _, ok := src["my-key"]; !ok {
		t.Error("source snapshot was mutated")
	}
}

func TestIsValidKey(t *testing.T) {
	cases := []struct {
		key   string
		valid bool
	}{
		{"VALID", true},
		{"_VALID", true},
		{"valid_123", true},
		{"1invalid", false},
		{"has space", false},
		{"", false},
	}
	for _, c := range cases {
		if got := isValidKey(c.key); got != c.valid {
			t.Errorf("isValidKey(%q) = %v, want %v", c.key, got, c.valid)
		}
	}
}
