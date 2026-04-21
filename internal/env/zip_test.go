package env

import (
	"testing"
)

func TestZipSnapshot_BothSides(t *testing.T) {
	left := Snapshot{"A": "1", "B": "2"}
	right := Snapshot{"B": "3", "C": "4"}

	opts := DefaultZipOptions()
	opts.Fn = func(_, l, r string) string { return l + "+" + r }

	out := ZipSnapshot(left, right, opts)

	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
	if out["B"] != "2+3" {
		t.Errorf("expected B=2+3, got %q", out["B"])
	}
	if out["C"] != "4" {
		t.Errorf("expected C=4, got %q", out["C"])
	}
}

func TestZipSnapshot_KeepLeftFalse(t *testing.T) {
	left := Snapshot{"A": "1", "B": "2"}
	right := Snapshot{"B": "3"}

	opts := DefaultZipOptions()
	opts.KeepLeft = false

	out := ZipSnapshot(left, right, opts)

	if _, ok := out["A"]; ok {
		t.Error("expected A to be dropped when KeepLeft=false")
	}
	if out["B"] != "3" {
		t.Errorf("expected B=3, got %q", out["B"])
	}
}

func TestZipSnapshot_KeepRightFalse(t *testing.T) {
	left := Snapshot{"A": "1"}
	right := Snapshot{"A": "2", "C": "99"}

	opts := DefaultZipOptions()
	opts.KeepRight = false

	out := ZipSnapshot(left, right, opts)

	if _, ok := out["C"]; ok {
		t.Error("expected C to be dropped when KeepRight=false")
	}
}

func TestZipSnapshot_NilFnRightWins(t *testing.T) {
	left := Snapshot{"X": "left"}
	right := Snapshot{"X": "right"}

	out := ZipSnapshot(left, right, DefaultZipOptions())

	if out["X"] != "right" {
		t.Errorf("expected right value to win, got %q", out["X"])
	}
}

func TestZipSnapshot_DoesNotMutateOriginals(t *testing.T) {
	left := Snapshot{"K": "v1"}
	right := Snapshot{"K": "v2"}

	_ = ZipSnapshot(left, right, DefaultZipOptions())

	if left["K"] != "v1" {
		t.Error("left snapshot was mutated")
	}
	if right["K"] != "v2" {
		t.Error("right snapshot was mutated")
	}
}

func TestZipSnapshot_EmptySnapshots(t *testing.T) {
	out := ZipSnapshot(Snapshot{}, Snapshot{}, DefaultZipOptions())
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}
