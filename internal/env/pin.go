package env

import "sync"

// PinOptions controls the behaviour of PinSnapshot.
type PinOptions struct {
	// Keys is the explicit set of keys to pin. If empty, all keys are pinned.
	Keys []string
}

// DefaultPinOptions returns a PinOptions that pins every key.
func DefaultPinOptions() PinOptions {
	return PinOptions{}
}

// Pinned holds an immutable copy of a Snapshot that can be safely shared
// across goroutines and compared against later snapshots.
type Pinned struct {
	mu   sync.RWMutex
	snap Snapshot
	keys map[string]struct{} // nil means "all keys"
}

// PinSnapshot captures snap according to opts and returns a *Pinned.
func PinSnapshot(snap Snapshot, opts PinOptions) *Pinned {
	keys := map[string]struct{}(nil)
	if len(opts.Keys) > 0 {
		keys = make(map[string]struct{}, len(opts.Keys))
		for _, k := range opts.Keys {
			keys[k] = struct{}{}
		}
	}

	clone := snap.Clone()
	if keys != nil {
		filtered := NewSnapshot()
		for k, v := range clone.m {
			if _, ok := keys[k]; ok {
				filtered.Set(k, v)
			}
		}
		clone = filtered
	}

	return &Pinned{snap: clone, keys: keys}
}

// Snapshot returns a clone of the pinned snapshot.
func (p *Pinned) Snapshot() Snapshot {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.snap.Clone()
}

// Update replaces the pinned snapshot with a new one derived from snap,
// respecting the original key filter.
func (p *Pinned) Update(snap Snapshot) {
	clone := snap.Clone()
	if p.keys != nil {
		filtered := NewSnapshot()
		for k, v := range clone.m {
			if _, ok := p.keys[k]; ok {
				filtered.Set(k, v)
			}
		}
		clone = filtered
	}
	p.mu.Lock()
	p.snap = clone
	p.mu.Unlock()
}

// Has reports whether key is present in the pinned snapshot.
func (p *Pinned) Has(key string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.snap.m[key]
	return ok
}
