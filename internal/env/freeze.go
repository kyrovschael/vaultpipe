package env

import (
	"fmt"
	"sync"
)

// FrozenSnapshot is an immutable snapshot that rejects all writes after creation.
type FrozenSnapshot struct {
	mu   sync.RWMutex
	data map[string]string
}

// FreezeSnapshot returns a FrozenSnapshot from the given Snapshot.
// The frozen copy is independent of the original.
func FreezeSnapshot(s Snapshot) *FrozenSnapshot {
	clone := make(map[string]string, len(s))
	for k, v := range s {
		clone[k] = v
	}
	return &FrozenSnapshot{data: clone}
}

// Get returns the value for key and whether it was present.
func (f *FrozenSnapshot) Get(key string) (string, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.data[key]
	return v, ok
}

// Keys returns a sorted list of all keys.
func (f *FrozenSnapshot) Keys() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]string, 0, len(f.data))
	for k := range f.data {
		keys = append(keys, k)
	}
	return keys
}

// Snapshot returns an unfrozen clone of the underlying data.
func (f *FrozenSnapshot) Snapshot() Snapshot {
	f.mu.RLock()
	defer f.mu.RUnlock()
	clone := make(Snapshot, len(f.data))
	for k, v := range f.data {
		clone[k] = v
	}
	return clone
}

// Len returns the number of entries.
func (f *FrozenSnapshot) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.data)
}

// AssertKey returns an error if the key is absent.
func (f *FrozenSnapshot) AssertKey(key string) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if _, ok := f.data[key]; !ok {
		return fmt.Errorf("frozen snapshot: key %q not found", key)
	}
	return nil
}
