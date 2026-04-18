package env

import (
	"testing"
)

func TestRedactSnapshot_SensitiveKeys(t *testing.T) {
	s := Snapshot{
		"SECRET_KEY":   "abc123",
		"HOME":         "/home/user",
		"TOKEN_VALUE":  "tok-xyz",
		"API_KEY_PROD": "prod-key",
		"PATH":         "/usr/bin",
	}
	redacted := RedactSnapshot(s)

	if redacted["SECRET_KEY"] != "[REDACTED]" {
		t.Errorf("expected SECRET_KEY to be redacted, got %q", redacted["SECRET_KEY"])
	}
	if redacted["TOKEN_VALUE"] != "[REDACTED]" {
		t.Errorf("expected TOKEN_VALUE to be redacted, got %q", redacted["TOKEN_VALUE"])
	}
	if redacted["API_KEY_PROD"] != "[REDACTED]" {
		t.Errorf("expected API_KEY_PROD to be redacted, got %q", redacted["API_KEY_PROD"])
	}
	if redacted["HOME"] != "/home/user" {
		t.Errorf("expected HOME to be unchanged, got %q", redacted["HOME"])
	}
	if redacted["PATH"] != "/usr/bin" {
		t.Errorf("expected PATH to be unchanged, got %q", redacted["PATH"])
	}
}

func TestRedactSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{"SECRET_KEY": "original"}
	_ = RedactSnapshot(s)
	if s["SECRET_KEY"] != "original" {
		t.Error("original snapshot was mutated")
	}
}

func TestRedactSlice_MixedEntries(t *testing.T) {
	env := []string{
		"PASSWORD=hunter2",
		"USER=alice",
		"AUTH_TOKEN=bearer-xyz",
		"NOEQUALS",
	}
	out := RedactSlice(env)

	expected := map[string]string{
		"PASSWORD":   "[REDACTED]",
		"USER":       "alice",
		"AUTH_TOKEN": "[REDACTED]",
	}
	for _, entry := range out {
		parts := splitEntry(entry)
		if len(parts) != 2 {
			continue
		}
		if want, ok := expected[parts[0]]; ok {
			if parts[1] != want {
				t.Errorf("key %s: got %q, want %q", parts[0], parts[1], want)
			}
		}
	}
}

func TestRedactSlice_PreservesNoEqualsEntries(t *testing.T) {
	env := []string{"NOEQUALS"}
	out := RedactSlice(env)
	if len(out) != 1 || out[0] != "NOEQUALS" {
		t.Errorf("expected NOEQUALS preserved, got %v", out)
	}
}

func splitEntry(entry string) []string {
	for i := 0; i < len(entry); i++ {
		if entry[i] == '=' {
			return []string{entry[:i], entry[i+1:]}
		}
	}
	return []string{entry}
}
