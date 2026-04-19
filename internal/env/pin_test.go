package env

import (
	"testing"
)

func TestPinSnapshot_AllKeys(t *testing.T) {
	snap := snapshotFrom(map[string]string{"A": "1", "B": "2"})
	p := PinSnapshot(snap, DefaultPinOptions())

	got := p.Snapshot()
	if got.m["A"] != "1" || got.m["B"] != "2" {
		t.Fatalf("unexpected snapshot: %v", got.m)
	}
}

func TestPinSnapshot_FilteredKeys(t *testing.T) {
	snap := snapshotFrom(map[string]string{"A": "1", "B": "2", "C": "3"})
	p := PinSnapshot(snap, PinOptions{Keys: []string{"A", "C"}})

	got := p.Snapshot()
	if _, ok := got.m["B"]; ok {
		t.Fatal("B should not be pinned")
	}
	if got.m["A"] != "1" || got.m["C"] != "3" {
		t.Fatalf("unexpected snapshot: %v", got.m)
	}
}

func TestPinSnapshot_DoesNotMutateOriginal(t *testing.T) {
	snap := snapshotFrom(map[string]string{"X": "original"})
	p := PinSnapshot(snap, DefaultPinOptions())

	snap.Set("X", "mutated")
	got := p.Snapshot()
	if got.m["X"] != "original" {
		t.Fatalf("pin was mutated: got %q", got.m["X"])
	}
}

func TestPinned_Update(t *testing.T) {
	snap := snapshotFrom(map[string]string{"K": "v1"})
	p := PinSnapshot(snap, DefaultPinOptions())

	newSnap := snapshotFrom(map[string]string{"K": "v2", "Z": "extra"})
	p.Update(newSnap)

	got := p.Snapshot()
	if got.m["K"] != "v2" {
		t.Fatalf("expected v2, got %q", got.m["K"])
	}
}

func TestPinned_Has(t *testing.T) {
	snap := snapshotFrom(map[string]string{"PRESENT": "yes"})
	p := PinSnapshot(snap, DefaultPinOptions())

	if !p.Has("PRESENT") {
		t.Fatal("expected PRESENT to be present")
	}
	if p.Has("ABSENT") {
		t.Fatal("expected ABSENT to be absent")
	}
}

func snapshotFrom(m map[string]string) Snapshot {
	s := NewSnapshot()
	for k, v := range m {
		s.Set(k, v)
	}
	return s
}
