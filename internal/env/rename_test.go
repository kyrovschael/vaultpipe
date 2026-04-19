package env

import (
	"testing"
)

func TestRenameSnapshot_BasicRename(t *testing.T) {
	s := Snapshot{"OLD_KEY": "val", "KEEP": "kept"}
	got := RenameSnapshot(s, map[string]string{"OLD_KEY": "NEW_KEY"}, DefaultRenameOptions())
	if got["NEW_KEY"] != "val" {
		t.Fatalf("expected NEW_KEY=val, got %q", got["NEW_KEY"])
	}
	if got["KEEP"] != "kept" {
		t.Fatalf("expected KEEP to pass through")
	}
	if _, ok := got["OLD_KEY"]; ok {
		t.Fatal("OLD_KEY should not be present after rename")
	}
}

func TestRenameSnapshot_DropUnmapped(t *testing.T) {
	s := Snapshot{"A": "1", "B": "2"}
	opts := RenameOptions{DropUnmapped: true}
	got := RenameSnapshot(s, map[string]string{"A": "ALPHA"}, opts)
	if got["ALPHA"] != "1" {
		t.Fatalf("expected ALPHA=1")
	}
	if _, ok := got["B"]; ok {
		t.Fatal("B should be dropped when DropUnmapped=true")
	}
}

func TestRenameSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{"X": "x"}
	_ = RenameSnapshot(s, map[string]string{"X": "Y"}, DefaultRenameOptions())
	if _, ok := s["X"]; !ok {
		t.Fatal("original snapshot was mutated")
	}
}

func TestRenameSnapshot_EmptyMapping(t *testing.T) {
	s := Snapshot{"A": "1", "B": "2"}
	got := RenameSnapshot(s, map[string]string{}, DefaultRenameOptions())
	if len(got) != len(s) {
		t.Fatalf("expected all keys to pass through, got %d", len(got))
	}
}

func TestRenameSnapshot_EmptySnapshot(t *testing.T) {
	got := RenameSnapshot(Snapshot{}, map[string]string{"A": "B"}, DefaultRenameOptions())
	if len(got) != 0 {
		t.Fatalf("expected empty snapshot")
	}
}
