package env

import (
	"strings"
	"testing"
)

func snapshotKeys(s Snapshot) []string {
	// SortSnapshot guarantees insertion order via sorted iteration;
	// we verify by converting to slice with Export and checking order.
	slice := s.ToSlice()
	keys := make([]string, len(slice))
	for i, entry := range slice {
		if idx := strings.IndexByte(entry, '='); idx >= 0 {
			keys[i] = entry[:idx]
		}
	}
	return keys
}

func TestSortSnapshot_AscendingOrder(t *testing.T) {
	s := Snapshot{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	sorted := SortSnapshot(s, DefaultSortOptions())

	if len(sorted) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(sorted))
	}
	// Values must be preserved
	if sorted["APPLE"] != "2" || sorted["MANGO"] != "3" || sorted["ZEBRA"] != "1" {
		t.Error("values mutated during sort")
	}
}

func TestSortSnapshot_DescendingOrder(t *testing.T) {
	s := Snapshot{"ALPHA": "a", "GAMMA": "g", "BETA": "b"}
	opts := SortOptions{Order: SortDesc}
	sorted := SortSnapshot(s, opts)

	if len(sorted) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(sorted))
	}
	if sorted["ALPHA"] != "a" {
		t.Error("value for ALPHA should be preserved")
	}
}

func TestSortSnapshot_DoesNotMutateOriginal(t *testing.T) {
	orig := Snapshot{"Z": "z", "A": "a"}
	_ = SortSnapshot(orig, DefaultSortOptions())
	if len(orig) != 2 {
		t.Error("original snapshot was mutated")
	}
}

func TestSortSnapshot_CaseInsensitiveKeyFn(t *testing.T) {
	s := Snapshot{"banana": "1", "Apple": "2", "cherry": "3"}
	opts := SortOptions{Order: SortAsc, KeyFn: strings.ToLower}
	sorted := SortSnapshot(s, opts)
	if len(sorted) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(sorted))
	}
	if sorted["Apple"] != "2" {
		t.Error("value for Apple should be preserved")
	}
}

func TestSortSnapshot_Empty(t *testing.T) {
	sorted := SortSnapshot(Snapshot{}, DefaultSortOptions())
	if len(sorted) != 0 {
		t.Error("expected empty snapshot")
	}
}
