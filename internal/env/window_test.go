package env

import (
	"testing"
)

func TestWindowSnapshot_BothBounds(t *testing.T) {
	s := Snapshot{
		{Key: "ALPHA", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "DELTA", Value: "3"},
		{Key: "GAMMA", Value: "4"},
		{Key: "ZETA", Value: "5"},
	}
	out := WindowSnapshot(s, WindowOptions{From: "BETA", To: "GAMMA"})
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	keys := make([]string, len(out))
	for i, e := range out {
		keys[i] = e.Key
	}
	for _, want := range []string{"BETA", "DELTA", "GAMMA"} {
		found := false
		for _, k := range keys {
			if k == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected key %q in output", want)
		}
	}
}

func TestWindowSnapshot_OpenLowerBound(t *testing.T) {
	s := Snapshot{
		{Key: "A", Value: "1"},
		{Key: "M", Value: "2"},
		{Key: "Z", Value: "3"},
	}
	out := WindowSnapshot(s, WindowOptions{To: "M"})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestWindowSnapshot_OpenUpperBound(t *testing.T) {
	s := Snapshot{
		{Key: "A", Value: "1"},
		{Key: "M", Value: "2"},
		{Key: "Z", Value: "3"},
	}
	out := WindowSnapshot(s, WindowOptions{From: "M"})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestWindowSnapshot_NoBounds_ReturnsAll(t *testing.T) {
	s := Snapshot{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	out := WindowSnapshot(s, DefaultWindowOptions())
	if len(out) != len(s) {
		t.Fatalf("expected %d entries, got %d", len(s), len(out))
	}
}

func TestWindowSnapshot_CaseFold(t *testing.T) {
	s := Snapshot{
		{Key: "alpha", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "gamma", Value: "3"},
	}
	// range [alpha, beta] case-insensitively
	out := WindowSnapshot(s, WindowOptions{From: "ALPHA", To: "BETA", CaseFold: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestWindowSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3"},
	}
	origLen := len(s)
	_ = WindowSnapshot(s, WindowOptions{From: "B", To: "B"})
	if len(s) != origLen {
		t.Fatalf("original snapshot mutated")
	}
}
