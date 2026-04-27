package env

import (
	"testing"
)

func TestKeysSnapshot_AllKeys(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2", "C": "3"}
	opts := DefaultKeysOptions()
	keys := KeysSnapshot(src, opts)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
}

func TestKeysSnapshot_SortedAscending(t *testing.T) {
	src := Snapshot{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	opts := DefaultKeysOptions()
	keys := KeysSnapshot(src, opts)
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Fatalf("unexpected order: %v", keys)
	}
}

func TestKeysSnapshot_Unsorted(t *testing.T) {
	src := Snapshot{"B": "2", "A": "1"}
	opts := DefaultKeysOptions()
	opts.Sorted = false
	keys := KeysSnapshot(src, opts)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestKeysSnapshot_PrefixFilter(t *testing.T) {
	src := Snapshot{"APP_HOST": "h", "APP_PORT": "p", "DB_HOST": "d"}
	opts := DefaultKeysOptions()
	opts.Prefix = "APP_"
	keys := KeysSnapshot(src, opts)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
	for _, k := range keys {
		if k == "DB_HOST" {
			t.Fatalf("DB_HOST should have been filtered out")
		}
	}
}

func TestKeysSnapshot_CaseFoldPrefix(t *testing.T) {
	src := Snapshot{"APP_HOST": "h", "app_port": "p", "DB_HOST": "d"}
	opts := DefaultKeysOptions()
	opts.Prefix = "app_"
	opts.CaseFold = true
	keys := KeysSnapshot(src, opts)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys with case-fold, got %d: %v", len(keys), keys)
	}
}

func TestKeysSnapshot_EmptySnapshot(t *testing.T) {
	src := Snapshot{}
	keys := KeysSnapshot(src, DefaultKeysOptions())
	if len(keys) != 0 {
		t.Fatalf("expected empty slice, got %v", keys)
	}
}

func TestKeysSnapshot_DoesNotMutateSource(t *testing.T) {
	src := Snapshot{"X": "1", "Y": "2"}
	origLen := len(src)
	_ = KeysSnapshot(src, DefaultKeysOptions())
	if len(src) != origLen {
		t.Fatal("source snapshot was mutated")
	}
}
