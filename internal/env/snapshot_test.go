package env

import (
	"sort"
	"testing"
)

func TestFromSlice_Basic(t *testing.T) {
	s := FromSlice([]string{"FOO=bar", "BAZ=qux"})
	if v, ok := s.Get("FOO"); !ok || v != "bar" {
		t.Fatalf("expected FOO=bar, got %q ok=%v", v, ok)
	}
	if v, ok := s.Get("BAZ"); !ok || v != "qux" {
		t.Fatalf("expected BAZ=qux, got %q ok=%v", v, ok)
	}
}

func TestFromSlice_EmptyValue(t *testing.T) {
	s := FromSlice([]string{"EMPTY="})
	v, ok := s.Get("EMPTY")
	if !ok || v != "" {
		t.Fatalf("expected empty string, got %q ok=%v", v, ok)
	}
}

func TestFromSlice_NoEquals(t *testing.T) {
	s := FromSlice([]string{"STANDALONE"})
	v, ok := s.Get("STANDALONE")
	if !ok || v != "" {
		t.Fatalf("expected standalone key with empty value, got %q ok=%v", v, ok)
	}
}

func TestSnapshot_ToSlice_RoundTrip(t *testing.T) {
	input := []string{"A=1", "B=2", "C=3"}
	s := FromSlice(input)
	out := s.ToSlice()
	sort.Strings(input)
	sort.Strings(out)
	if len(out) != len(input) {
		t.Fatalf("length mismatch: want %d got %d", len(input), len(out))
	}
	for i := range input {
		if input[i] != out[i] {
			t.Errorf("mismatch at %d: want %q got %q", i, input[i], out[i])
		}
	}
}

func TestSnapshot_Merge_OverlayWins(t *testing.T) {
	base := FromSlice([]string{"FOO=original", "KEEP=yes"})
	merged := base.Merge(map[string]string{"FOO": "overridden", "NEW": "value"})

	if v, _ := merged.Get("FOO"); v != "overridden" {
		t.Errorf("expected FOO=overridden, got %q", v)
	}
	if v, _ := merged.Get("KEEP"); v != "yes" {
		t.Errorf("expected KEEP=yes, got %q", v)
	}
	if v, _ := merged.Get("NEW"); v != "value" {
		t.Errorf("expected NEW=value, got %q", v)
	}
}

func TestSnapshot_Merge_DoesNotMutateBase(t *testing.T) {
	base := FromSlice([]string{"FOO=original"})
	base.Merge(map[string]string{"FOO": "changed"})
	if v, _ := base.Get("FOO"); v != "original" {
		t.Errorf("base was mutated: got %q", v)
	}
}
