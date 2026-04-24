package env

import (
	"strings"
	"testing"
)

func TestLimitSnapshot_MaxKeys_Truncates(t *testing.T) {
	s := Snapshot{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}, {Key: "C", Value: "3"}}
	out, err := LimitSnapshot(s, LimitOptions{MaxKeys: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestLimitSnapshot_MaxKeys_Strict_Error(t *testing.T) {
	s := Snapshot{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}, {Key: "C", Value: "3"}}
	_, err := LimitSnapshot(s, LimitOptions{MaxKeys: 2, StrictKeys: true})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds max") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLimitSnapshot_MaxValueLen_DropsLong(t *testing.T) {
	s := Snapshot{
		{Key: "SHORT", Value: "hi"},
		{Key: "LONG", Value: strings.Repeat("x", 200)},
	}
	out, err := LimitSnapshot(s, LimitOptions{MaxValueLen: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "SHORT" {
		t.Errorf("expected only SHORT, got %+v", out)
	}
}

func TestLimitSnapshot_MaxValueLen_Strict_Error(t *testing.T) {
	s := Snapshot{
		{Key: "LONG", Value: strings.Repeat("x", 200)},
	}
	_, err := LimitSnapshot(s, LimitOptions{MaxValueLen: 10, StrictValues: true})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds max length") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLimitSnapshot_NoLimits_ReturnsAll(t *testing.T) {
	s := Snapshot{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	out, err := LimitSnapshot(s, DefaultLimitOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(s) {
		t.Errorf("expected %d entries, got %d", len(s), len(out))
	}
}

func TestLimitSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{
		{Key: "A", Value: "1"},
		{Key: "B", Value: strings.Repeat("z", 50)},
		{Key: "C", Value: "3"},
	}
	orig := make(Snapshot, len(s))
	copy(orig, s)

	_, _ = LimitSnapshot(s, LimitOptions{MaxKeys: 1, MaxValueLen: 5})

	for i, e := range s {
		if e != orig[i] {
			t.Errorf("original mutated at index %d", i)
		}
	}
}

func TestLimitSnapshot_EmptySnapshot(t *testing.T) {
	out, err := LimitSnapshot(Snapshot{}, LimitOptions{MaxKeys: 5, MaxValueLen: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty snapshot, got %d entries", len(out))
	}
}
