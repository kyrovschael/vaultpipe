package env

import (
	"testing"
)

func TestDedupeSnapshot_FirstWins(t *testing.T) {
	s := Snapshot{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "A", Value: "99"},
	}
	out := DedupeSnapshot(s, DefaultDedupeOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Value != "1" {
		t.Errorf("expected first A=1, got %s", out[0].Value)
	}
}

func TestDedupeSnapshot_LastWins(t *testing.T) {
	s := Snapshot{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "A", Value: "99"},
	}
	out := DedupeSnapshot(s, DedupeOptions{LastWins: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Key == "A" && e.Value != "99" {
			t.Errorf("expected last A=99, got %s", e.Value)
		}
	}
}

func TestDedupeSnapshot_NoDuplicates(t *testing.T) {
	s := Snapshot{{Key: "X", Value: "1"}, {Key: "Y", Value: "2"}}
	out := DedupeSnapshot(s, DefaultDedupeOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestDedupeSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{{Key: "A", Value: "1"}, {Key: "A", Value: "2"}}
	_ = DedupeSnapshot(s, DefaultDedupeOptions())
	if s[1].Value != "2" {
		t.Error("original snapshot was mutated")
	}
}

func TestDedupeSlice_RemovesDuplicates(t *testing.T) {
	pairs := []string{"FOO=bar", "BAZ=qux", "FOO=override"}
	out, err := DedupeSlice(pairs, DefaultDedupeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count := 0
	for _, p := range out {
		if len(p) >= 3 && p[:3] == "FOO" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 FOO entry, got %d", count)
	}
}
