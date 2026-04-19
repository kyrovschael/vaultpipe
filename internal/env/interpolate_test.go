package env

import (
	"testing"
)

func TestInterpolate_Simple(t *testing.T) {
	snap := Snapshot{"DSN": "postgres://${USER}:${PASS}@host/db"}
	secrets := Snapshot{"USER": "alice", "PASS": "s3cr3t"}
	out, err := InterpolateSnapshot(snap, secrets, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out["DSN"], "postgres://alice:s3cr3t@host/db"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolate_MissingKeyLenient(t *testing.T) {
	snap := Snapshot{"URL": "https://${HOST}/path"}
	out, err := InterpolateSnapshot(snap, Snapshot{}, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out["URL"], "https://${HOST}/path"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolate_MissingKeyStrict(t *testing.T) {
	snap := Snapshot{"URL": "https://${HOST}/path"}
	opts := DefaultInterpolateOptions()
	opts.Strict = true
	_, err := InterpolateSnapshot(snap, Snapshot{}, opts)
	if err == nil {
		t.Fatal("expected error for missing placeholder in strict mode")
	}
}

func TestInterpolate_NoPlaceholder(t *testing.T) {
	snap := Snapshot{"KEY": "plain-value"}
	out, err := InterpolateSnapshot(snap, Snapshot{"KEY": "other"}, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := out["KEY"], "plain-value"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInterpolate_DoesNotMutateOriginal(t *testing.T) {
	origVal := "hello ${NAME}"
	snap := Snapshot{"GREETING": origVal}
	secrets := Snapshot{"NAME": "world"}
	_, err := InterpolateSnapshot(snap, secrets, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap["GREETING"] != origVal {
		t.Errorf("original snapshot mutated")
	}
}

func TestInterpolate_UnclosedDelimiter(t *testing.T) {
	snap := Snapshot{"VAL": "prefix_${UNCLOSED"}
	out, err := InterpolateSnapshot(snap, Snapshot{"UNCLOSED": "x"}, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No closing brace — literal preserved.
	if got, want := out["VAL"], "prefix_${UNCLOSED"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
