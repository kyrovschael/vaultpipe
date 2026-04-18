package env

import (
	"testing"
)

func TestExpand_SimpleReference(t *testing.T) {
	snap := Snapshot{"HOME": "/home/user", "PATH": "$HOME/bin"}
	out := Expand(snap, DefaultExpandOptions())
	if out["PATH"] != "/home/user/bin" {
		t.Fatalf("expected /home/user/bin, got %q", out["PATH"])
	}
}

func TestExpand_BraceReference(t *testing.T) {
	snap := Snapshot{"BASE": "/opt", "DIR": "${BASE}/app"}
	out := Expand(snap, DefaultExpandOptions())
	if out["DIR"] != "/opt/app" {
		t.Fatalf("expected /opt/app, got %q", out["DIR"])
	}
}

func TestExpand_MissingKeyNoFallback(t *testing.T) {
	snap := Snapshot{"VAL": "$UNDEFINED"}
	out := Expand(snap, DefaultExpandOptions())
	if out["VAL"] != "" {
		t.Fatalf("expected empty string, got %q", out["VAL"])
	}
}

func TestExpand_FallbackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	snap := Snapshot{"VAL": "$OS_VAR"}
	opts := DefaultExpandOptions()
	opts.FallbackToOS = true
	out := Expand(snap, opts)
	if out["VAL"] != "from-os" {
		t.Fatalf("expected from-os, got %q", out["VAL"])
	}
}

func TestExpand_NoExpandKey(t *testing.T) {
	snap := Snapshot{"HOME": "/home/user", "LITERAL": "$HOME"}
	opts := DefaultExpandOptions()
	opts.NoExpand = map[string]bool{"LITERAL": true}
	out := Expand(snap, opts)
	if out["LITERAL"] != "$HOME" {
		t.Fatalf("expected literal $HOME, got %q", out["LITERAL"])
	}
}

func TestExpand_DoesNotMutateOriginal(t *testing.T) {
	snap := Snapshot{"A": "1", "B": "$A"}
	_ = Expand(snap, DefaultExpandOptions())
	if snap["B"] != "$A" {
		t.Fatal("original snapshot was mutated")
	}
}
