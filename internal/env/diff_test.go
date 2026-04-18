package env

import (
	"testing"
)

func snapshotOf(pairs ...string) Snapshot {
	s := make(Snapshot, len(pairs)/2)
	for i := 0; i < len(pairs)-1; i += 2 {
		s[pairs[i]] = pairs[i+1]
	}
	return s
}

func findEntry(entries []DiffEntry, key string) (DiffEntry, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e, true
		}
	}
	return DiffEntry{}, false
}

func TestDiffSnapshots_Added(t *testing.T) {
	base := snapshotOf("A", "1")
	overlay := snapshotOf("A", "1", "B", "2")
	entries := DiffSnapshots(base, overlay, DefaultDiffOptions())
	e, ok := findEntry(entries, "B")
	if !ok || e.Kind != DiffAdded || e.NewValue != "2" {
		t.Fatalf("expected B added, got %+v", entries)
	}
}

func TestDiffSnapshots_Removed(t *testing.T) {
	base := snapshotOf("A", "1", "B", "2")
	overlay := snapshotOf("A", "1")
	entries := DiffSnapshots(base, overlay, DefaultDiffOptions())
	e, ok := findEntry(entries, "B")
	if !ok || e.Kind != DiffRemoved || e.OldValue != "2" {
		t.Fatalf("expected B removed, got %+v", entries)
	}
}

func TestDiffSnapshots_Changed(t *testing.T) {
	base := snapshotOf("A", "old")
	overlay := snapshotOf("A", "new")
	entries := DiffSnapshots(base, overlay, DefaultDiffOptions())
	e, ok := findEntry(entries, "A")
	if !ok || e.Kind != DiffChanged || e.OldValue != "old" || e.NewValue != "new" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestDiffSnapshots_NoChange(t *testing.T) {
	base := snapshotOf("A", "1")
	overlay := snapshotOf("A", "1")
	entries := DiffSnapshots(base, overlay, DefaultDiffOptions())
	if len(entries) != 0 {
		t.Fatalf("expected no diff, got %+v", entries)
	}
}

func TestDiffSnapshots_RedactValues(t *testing.T) {
	base := snapshotOf("SECRET", "hunter2")
	overlay := snapshotOf("SECRET", "p@ssw0rd")
	opts := DiffOptions{RedactValues: true}
	entries := DiffSnapshots(base, overlay, opts)
	e, ok := findEntry(entries, "SECRET")
	if !ok || e.OldValue != "***" || e.NewValue != "***" {
		t.Fatalf("expected redacted values, got %+v", e)
	}
}

func TestDiffSnapshots_Empty(t *testing.T) {
	entries := DiffSnapshots(Snapshot{}, Snapshot{}, DefaultDiffOptions())
	if len(entries) != 0 {
		t.Fatalf("expected empty diff")
	}
}
