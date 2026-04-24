package env

import (
	"fmt"
	"time"
)

// TTLEntry holds a snapshot value alongside its expiry time.
type TTLEntry struct {
	Snapshot Snapshot
	ExpiresAt time.Time
}

// Expired reports whether the entry has passed its expiry time.
func (e TTLEntry) Expired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// DefaultTTLOptions returns a TTLOptions with a 5-minute TTL and no
// restricted keys.
func DefaultTTLOptions() TTLOptions {
	return TTLOptions{
		TTL: 5 * time.Minute,
	}
}

// TTLOptions controls the behaviour of TTLSnapshot.
type TTLOptions struct {
	// TTL is the duration after which entries are considered expired.
	TTL time.Duration
	// Keys restricts expiry tracking to the listed keys. When empty all
	// keys are tracked.
	Keys []string
	// Now overrides the clock used to stamp entries. Defaults to
	// time.Now when nil.
	Now func() time.Time
}

// TTLSnapshot stamps every entry in src with an expiry time derived
// from opts.TTL and returns a map of key → TTLEntry. Keys absent from
// opts.Keys (when non-empty) are passed through without a TTL entry.
func TTLSnapshot(src Snapshot, opts TTLOptions) (map[string]TTLEntry, error) {
	if opts.TTL <= 0 {
		return nil, fmt.Errorf("env: TTL must be positive, got %v", opts.TTL)
	}
	now := time.Now
	if opts.Now != nil {
		now = opts.Now
	}
	filter := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		filter[k] = struct{}{}
	}
	out := make(map[string]TTLEntry, len(src))
	for k, v := range src {
		if len(filter) > 0 {
			if _, ok := filter[k]; !ok {
				continue
			}
		}
		out[k] = TTLEntry{
			Snapshot:  Snapshot{k: v},
			ExpiresAt: now().Add(opts.TTL),
		}
	}
	return out, nil
}
