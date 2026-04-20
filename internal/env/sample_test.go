package env

import (
	"testing"
)

func buildSnapshot(pairs map[string]string) Snapshot {
	s := make(Snapshot, len(pairs))
	for k, v := range pairs {
		s[k] = v
	}
	return s
}

func TestSampleSnapshot_ReturnsNEntries(t *testing.T) {
	src := buildSnapshot(map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5"})
	out, err := SampleSnapshot(src, SampleOptions{N: 3, Seed: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestSampleSnapshot_NGreaterThanLen_ReturnsAll(t *testing.T) {
	src := buildSnapshot(map[string]string{"X": "1", "Y": "2"})
	out, err := SampleSnapshot(src, SampleOptions{N: 100, Seed: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestSampleSnapshot_ZeroN_ReturnsEmpty(t *testing.T) {
	src := buildSnapshot(map[string]string{"A": "1", "B": "2"})
	out, err := SampleSnapshot(src, SampleOptions{N: 0, Seed: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}

func TestSampleSnapshot_NegativeN_ReturnsError(t *testing.T) {
	src := buildSnapshot(map[string]string{"A": "1"})
	_, err := SampleSnapshot(src, SampleOptions{N: -1, Seed: 0})
	if err == nil {
		t.Fatal("expected error for negative N")
	}
}

func TestSampleSnapshot_Deterministic(t *testing.T) {
	src := buildSnapshot(map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"})
	opts := SampleOptions{N: 2, Seed: 99}

	out1, _ := SampleSnapshot(src, opts)
	out2, _ := SampleSnapshot(src, opts)

	for k := range out1 {
		if _, ok := out2[k]; !ok {
			t.Errorf("key %q present in first sample but not second", k)
		}
	}
}

func TestSampleSnapshot_DoesNotMutateSource(t *testing.T) {
	src := buildSnapshot(map[string]string{"A": "1", "B": "2", "C": "3"})
	origLen := len(src)
	_, _ = SampleSnapshot(src, SampleOptions{N: 1, Seed: 7})
	if len(src) != origLen {
		t.Errorf("source snapshot was mutated")
	}
}
