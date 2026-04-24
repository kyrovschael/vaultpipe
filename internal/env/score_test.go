package env

import (
	"testing"
)

func TestScoreSnapshot_DescendingOrder(t *testing.T) {
	s := Snapshot{Entries: []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3"},
	}}
	scores := map[string]float64{"A": 1.0, "B": 3.0, "C": 2.0}
	out, err := ScoreSnapshot(s, DefaultScoreOptions(), func(k, _ string) float64 {
		return scores[k]
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"B", "C", "A"}
	for i, e := range out.Entries {
		if e.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Key, want[i])
		}
	}
}

func TestScoreSnapshot_AscendingOrder(t *testing.T) {
	s := Snapshot{Entries: []Entry{
		{Key: "X", Value: "a"},
		{Key: "Y", Value: "b"},
	}}
	opts := DefaultScoreOptions()
	opts.Ascending = true
	out, err := ScoreSnapshot(s, opts, func(k, _ string) float64 {
		if k == "X" {
			return 5.0
		}
		return 1.0
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Key != "Y" {
		t.Errorf("expected Y first, got %q", out.Entries[0].Key)
	}
}

func TestScoreSnapshot_RestrictedKeys(t *testing.T) {
	s := Snapshot{Entries: []Entry{
		{Key: "VAULT_TOKEN", Value: "tok"},
		{Key: "HOME", Value: "/root"},
		{Key: "VAULT_ADDR", Value: "http://localhost"},
	}}
	opts := DefaultScoreOptions()
	opts.Keys = []string{"VAULT_TOKEN", "VAULT_ADDR"}
	out, err := ScoreSnapshot(s, opts, func(k, _ string) float64 {
		if k == "VAULT_TOKEN" {
			return 10.0
		}
		return 2.0
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// VAULT_TOKEN (10) > VAULT_ADDR (2), HOME appended last
	if out.Entries[0].Key != "VAULT_TOKEN" {
		t.Errorf("expected VAULT_TOKEN first, got %q", out.Entries[0].Key)
	}
	if out.Entries[2].Key != "HOME" {
		t.Errorf("expected HOME last, got %q", out.Entries[2].Key)
	}
}

func TestScoreSnapshot_NilFnError(t *testing.T) {
	s := Snapshot{Entries: []Entry{{Key: "K", Value: "v"}}}
	_, err := ScoreSnapshot(s, DefaultScoreOptions(), nil)
	if err == nil {
		t.Fatal("expected error for nil fn")
	}
}

func TestScoreSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{Entries: []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}}
	origOrder := []string{s.Entries[0].Key, s.Entries[1].Key}
	_, _ = ScoreSnapshot(s, DefaultScoreOptions(), func(k, _ string) float64 {
		if k == "B" {
			return 99.0
		}
		return 0.0
	})
	for i, e := range s.Entries {
		if e.Key != origOrder[i] {
			t.Errorf("original mutated at pos %d: got %q", i, e.Key)
		}
	}
}
