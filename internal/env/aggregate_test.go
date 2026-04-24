package env

import (
	"testing"
)

func TestAggregateSnapshot_JoinsMatchingKeys(t *testing.T) {
	snaps := []Snapshot{
		{"TAGS": "alpha", "REGION": "us-east-1"},
		{"TAGS": "beta"},
		{"TAGS": "gamma", "ENV": "prod"},
	}
	got, err := AggregateSnapshot(snaps, DefaultAggregateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["TAGS"] != "alpha,beta,gamma" {
		t.Errorf("TAGS = %q, want %q", got["TAGS"], "alpha,beta,gamma")
	}
	if got["REGION"] != "us-east-1" {
		t.Errorf("REGION = %q, want %q", got["REGION"], "us-east-1")
	}
	if got["ENV"] != "prod" {
		t.Errorf("ENV = %q, want %q", got["ENV"], "prod")
	}
}

func TestAggregateSnapshot_SortValues(t *testing.T) {
	snaps := []Snapshot{
		{"TAGS": "gamma"},
		{"TAGS": "alpha"},
		{"TAGS": "beta"},
	}
	got, err := AggregateSnapshot(snaps, AggregateOptions{Separator: ",", SortValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["TAGS"] != "alpha,beta,gamma" {
		t.Errorf("TAGS = %q, want sorted %q", got["TAGS"], "alpha,beta,gamma")
	}
}

func TestAggregateSnapshot_RestrictedKeys(t *testing.T) {
	snaps := []Snapshot{
		{"TAGS": "a", "SECRET": "x"},
		{"TAGS": "b", "SECRET": "y"},
	}
	got, err := AggregateSnapshot(snaps, AggregateOptions{Separator: "|", Keys: []string{"TAGS"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["TAGS"] != "a|b" {
		t.Errorf("TAGS = %q, want %q", got["TAGS"], "a|b")
	}
	if _, ok := got["SECRET"]; ok {
		t.Error("SECRET should have been excluded by key restriction")
	}
}

func TestAggregateSnapshot_EmptySnapshots(t *testing.T) {
	got, err := AggregateSnapshot(nil, DefaultAggregateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty snapshot, got %v", got)
	}
}

func TestAggregateSnapshot_SingleSnapshot(t *testing.T) {
	snaps := []Snapshot{{"A": "1", "B": "2"}}
	got, err := AggregateSnapshot(snaps, DefaultAggregateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "1" || got["B"] != "2" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestAggregateSnapshot_DoesNotMutateInputs(t *testing.T) {
	original := Snapshot{"K": "v1"}
	snaps := []Snapshot{original, {"K": "v2"}}
	_, err := AggregateSnapshot(snaps, DefaultAggregateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if original["K"] != "v1" {
		t.Error("input snapshot was mutated")
	}
}

func TestAggregateSnapshot_CustomSeparator(t *testing.T) {
	snaps := []Snapshot{{"X": "foo"}, {"X": "bar"}}
	got, err := AggregateSnapshot(snaps, AggregateOptions{Separator: " "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["X"] != "foo bar" {
		t.Errorf("X = %q, want %q", got["X"], "foo bar")
	}
}
