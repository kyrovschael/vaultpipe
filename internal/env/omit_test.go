package env

import (
	"testing"
)

func TestOmitSnapshot_ExactMatch(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2", "C": "3"}
	opts := DefaultOmitOptions()
	opts.Keys = []string{"B"}

	got := OmitSnapshot(src, opts)

	if _, ok := got["B"]; ok {
		t.Error("expected B to be omitted")
	}
	if got["A"] != "1" || got["C"] != "3" {
		t.Error("expected A and C to be preserved")
	}
}

func TestOmitSnapshot_PrefixMatch(t *testing.T) {
	src := Snapshot{"SECRET_TOKEN": "x", "SECRET_KEY": "y", "SAFE": "z"}
	opts := DefaultOmitOptions()
	opts.Keys = []string{"SECRET_"}
	opts.PrefixMatch = true

	got := OmitSnapshot(src, opts)

	if _, ok := got["SECRET_TOKEN"]; ok {
		t.Error("expected SECRET_TOKEN to be omitted")
	}
	if _, ok := got["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY to be omitted")
	}
	if got["SAFE"] != "z" {
		t.Error("expected SAFE to be preserved")
	}
}

func TestOmitSnapshot_CaseFold(t *testing.T) {
	src := Snapshot{"Database_URL": "postgres://", "PORT": "5432"}
	opts := DefaultOmitOptions()
	opts.Keys = []string{"database_url"}
	opts.CaseFold = true

	got := OmitSnapshot(src, opts)

	if _, ok := got["Database_URL"]; ok {
		t.Error("expected Database_URL to be omitted via case-fold")
	}
	if got["PORT"] != "5432" {
		t.Error("expected PORT to be preserved")
	}
}

func TestOmitSnapshot_DoesNotMutateOriginal(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2"}
	opts := DefaultOmitOptions()
	opts.Keys = []string{"A"}

	_ = OmitSnapshot(src, opts)

	if _, ok := src["A"]; !ok {
		t.Error("OmitSnapshot must not mutate the source snapshot")
	}
}

func TestOmitSnapshot_EmptyKeys(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2"}
	opts := DefaultOmitOptions()

	got := OmitSnapshot(src, opts)

	if len(got) != len(src) {
		t.Errorf("expected all %d keys preserved, got %d", len(src), len(got))
	}
}

func TestOmitSnapshot_NoMatch(t *testing.T) {
	src := Snapshot{"X": "1", "Y": "2"}
	opts := DefaultOmitOptions()
	opts.Keys = []string{"Z"}

	got := OmitSnapshot(src, opts)

	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}
