package env

import (
	"testing"
	"time"
)

func TestCheckpoint_NameAndTime(t *testing.T) {
	before := time.Now()
	cp := NewCheckpoint("baseline", Snapshot{"A": "1"})
	after := time.Now()

	if cp.Name() != "baseline" {
		t.Fatalf("want baseline, got %s", cp.Name())
	}
	if cp.RecordedAt().Before(before) || cp.RecordedAt().After(after) {
		t.Fatal("recorded time out of range")
	}
}

func TestCheckpoint_SnapshotIsClone(t *testing.T) {
	orig := Snapshot{"X": "1"}
	cp := NewCheckpoint("test", orig)
	orig["X"] = "mutated"

	if cp.Snapshot()["X"] != "1" {
		t.Fatal("checkpoint should not reflect mutation of original")
	}
}

func TestCheckpoint_Update(t *testing.T) {
	cp := NewCheckpoint("c", Snapshot{"K": "old"})
	t1 := cp.RecordedAt()

	time.Sleep(2 * time.Millisecond)
	cp.Update(Snapshot{"K": "new"})

	if cp.Snapshot()["K"] != "new" {
		t.Fatal("expected updated value")
	}
	if !cp.RecordedAt().After(t1) {
		t.Fatal("recorded time should advance after Update")
	}
}

func TestCheckpoint_DiffFrom(t *testing.T) {
	cp := NewCheckpoint("base", Snapshot{"A": "1", "B": "2"})
	current := Snapshot{"A": "changed", "C": "3"}

	entries := cp.DiffFrom(current, DefaultDiffOptions())

	kinds := map[string]DiffKind{}
	for _, e := range entries {
		kinds[e.Key] = e.Kind
	}

	if kinds["A"] != DiffChanged {
		t.Errorf("expected A changed, got %v", kinds["A"])
	}
	if kinds["B"] != DiffRemoved {
		t.Errorf("expected B removed, got %v", kinds["B"])
	}
	if kinds["C"] != DiffAdded {
		t.Errorf("expected C added, got %v", kinds["C"])
	}
}
