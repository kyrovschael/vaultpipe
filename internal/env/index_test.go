package env

import (
	"testing"
)

func TestIndexSnapshot_Basic(t *testing.T) {
	snap := FromSlice([]string{"FOO=bar", "BAZ=qux", "QUX=bar"})
	idx, err := IndexSnapshot(snap, DefaultIndexOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, ok := idx.Lookup("bar")
	if !ok {
		t.Fatal("expected 'bar' to be in index")
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys for 'bar', got %d", len(keys))
	}
}

func TestIndexSnapshot_SingleValue(t *testing.T) {
	snap := FromSlice([]string{"ONLY=unique"})
	idx, err := IndexSnapshot(snap, DefaultIndexOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, ok := idx.Lookup("unique")
	if !ok || len(keys) != 1 || keys[0] != "ONLY" {
		t.Fatalf("unexpected result: ok=%v keys=%v", ok, keys)
	}
}

func TestIndexSnapshot_RestrictedKeys(t *testing.T) {
	snap := FromSlice([]string{"FOO=shared", "BAR=shared", "BAZ=other"})
	opts := DefaultIndexOptions()
	opts.Keys = []string{"FOO"}
	idx, err := IndexSnapshot(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, ok := idx.Lookup("shared")
	if !ok || len(keys) != 1 || keys[0] != "FOO" {
		t.Fatalf("expected only FOO, got %v", keys)
	}
	if _, ok := idx.Lookup("other"); ok {
		t.Fatal("BAZ should have been excluded")
	}
}

func TestIndexSnapshot_CaseFold(t *testing.T) {
	snap := FromSlice([]string{"FOO=Hello", "BAR=HELLO"})
	opts := DefaultIndexOptions()
	opts.CaseFold = true
	idx, err := IndexSnapshot(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, ok := idx.Lookup("hello")
	if !ok || len(keys) != 2 {
		t.Fatalf("expected 2 keys under 'hello', got %v", keys)
	}
}

func TestIndexSnapshot_MissingValueNotFound(t *testing.T) {
	snap := FromSlice([]string{"FOO=bar"})
	idx, err := IndexSnapshot(snap, DefaultIndexOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := idx.Lookup("nonexistent")
	if ok {
		t.Fatal("expected lookup to return false for missing value")
	}
}

func TestIndexSnapshot_EmptySnapshot(t *testing.T) {
	snap := FromSlice(nil)
	idx, err := IndexSnapshot(snap, DefaultIndexOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx) != 0 {
		t.Fatalf("expected empty index, got %d entries", len(idx))
	}
}
