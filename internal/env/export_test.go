package env

import (
	"strings"
	"testing"
)

func TestExport_ShellFormat(t *testing.T) {
	snap := Snapshot{"FOO": "bar", "BAZ": "qux"}
	out := Export(snap, FormatShell, false)

	if !strings.Contains(out, "export FOO='bar'") {
		t.Errorf("expected export FOO='bar' in output, got: %s", out)
	}
	if !strings.Contains(out, "export BAZ='qux'") {
		t.Errorf("expected export BAZ='qux' in output, got: %s", out)
	}
}

func TestExport_DotenvFormat(t *testing.T) {
	snap := Snapshot{"KEY": "value"}
	out := Export(snap, FormatDotenv, false)

	if !strings.Contains(out, "KEY=value") {
		t.Errorf("expected KEY=value in output, got: %s", out)
	}
}

func TestExport_Redacted(t *testing.T) {
	snap := Snapshot{"SECRET_TOKEN": "supersecret", "HOME": "/home/user"}
	out := Export(snap, FormatShell, true)

	if strings.Contains(out, "supersecret") {
		t.Errorf("expected secret to be redacted, got: %s", out)
	}
	if !strings.Contains(out, "HOME") {
		t.Errorf("expected HOME to be present in output")
	}
}

func TestExport_SortedOutput(t *testing.T) {
	snap := Snapshot{"Z": "last", "A": "first", "M": "mid"}
	out := Export(snap, FormatDotenv, false)
	lines := strings.Split(strings.TrimSpace(out), "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A=") {
		t.Errorf("expected first line to start with A=, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z=") {
		t.Errorf("expected last line to start with Z=, got %s", lines[2])
	}
}

func TestShellQuote_EscapesSingleQuote(t *testing.T) {
	result := shellQuote("it's")
	expected := "'it'\\'''"
	if result != expected {
		t.Errorf("shellQuote: expected %s, got %s", expected, result)
	}
}
