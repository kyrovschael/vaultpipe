package env

import (
	"testing"
)

func TestPickSnapshot_ReturnsRequestedKeys(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2", "C": "3"}
	out, err := PickSnapshot(src, []string{"A", "C"}, DefaultPickOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Errorf("unexpected values: %v", out)
	}
}

func TestPickSnapshot_IgnoresMissingByDefault(t *testing.T) {
	src := Snapshot{"A": "1"}
	out, err := PickSnapshot(src, []string{"A", "MISSING"}, DefaultPickOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestPickSnapshot_StrictErrorOnMissing(t *testing.T) {
	src := Snapshot{"A": "1"}
	opts := DefaultPickOptions()
	opts.Strict = true
	_, err := PickSnapshot(src, []string{"A", "MISSING"}, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestPickSnapshot_CaseFold(t *testing.T) {
	src := Snapshot{"MY_KEY": "hello"}
	opts := DefaultPickOptions()
	opts.CaseFold = true
	out, err := PickSnapshot(src, []string{"my_key"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MY_KEY"] != "hello" {
		t.Errorf("expected MY_KEY=hello, got %v", out)
	}
}

func TestPickSnapshot_DoesNotMutateOriginal(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2"}
	out, _ := PickSnapshot(src, []string{"A"}, DefaultPickOptions())
	out["A"] = "mutated"
	if src["A"] != "1" {
		t.Error("PickSnapshot mutated the source snapshot")
	}
}

func TestPickSnapshot_EmptyKeys_ReturnsEmpty(t *testing.T) {
	src := Snapshot{"A": "1"}
	out, err := PickSnapshot(src, []string{}, DefaultPickOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty snapshot, got %v", out)
	}
}
