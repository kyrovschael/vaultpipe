package env

import "testing"

func TestSnapshotsEqual_DifferentLengths(t *testing.T) {
	a := Snapshot{"A": "1"}
	b := Snapshot{"A": "1", "B": "2"}
	if snapshotsEqual(a, b) {
		t.Fatal("different lengths should not be equal")
	}
}

func TestSnapshotsEqual_BothEmpty(t *testing.T) {
	if !snapshotsEqual(Snapshot{}, Snapshot{}) {
		t.Fatal("two empty snapshots should be equal")
	}
}

func TestSnapshotsEqual_NilAndEmpty(t *testing.T) {
	// nil map and empty map have len 0, treated as equal.
	if !snapshotsEqual(nil, Snapshot{}) {
		t.Fatal("nil and empty snapshot should be equal")
	}
}

func TestSnapshotsEqual_SameKeys_DifferentValues(t *testing.T) {
	a := Snapshot{"X": "foo"}
	b := Snapshot{"X": "bar"}
	if snapshotsEqual(a, b) {
		t.Fatal("different values should not be equal")
	}
}
