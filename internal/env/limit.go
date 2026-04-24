package env

import (
	"errors"
	"fmt"
)

// DefaultLimitOptions returns a LimitOptions with sensible defaults.
func DefaultLimitOptions() LimitOptions {
	return LimitOptions{
		MaxKeys:      0,
		MaxValueLen:  0,
		StrictKeys:   false,
		StrictValues: false,
	}
}

// LimitOptions controls the behaviour of LimitSnapshot.
type LimitOptions struct {
	// MaxKeys is the maximum number of keys allowed. 0 means unlimited.
	MaxKeys int
	// MaxValueLen is the maximum byte length of any single value. 0 means unlimited.
	MaxValueLen int
	// StrictKeys causes an error when the snapshot exceeds MaxKeys.
	StrictKeys bool
	// StrictValues causes an error when a value exceeds MaxValueLen.
	StrictValues bool
}

// LimitSnapshot enforces upper bounds on the number of keys and the length of
// individual values within a Snapshot. When a limit is zero it is not applied.
// Entries that violate a non-strict limit are silently dropped; strict mode
// returns an error instead.
func LimitSnapshot(s Snapshot, opts LimitOptions) (Snapshot, error) {
	out := make(Snapshot, 0, len(s))

	for _, e := range s {
		if opts.MaxValueLen > 0 && len(e.Value) > opts.MaxValueLen {
			if opts.StrictValues {
				return nil, fmt.Errorf("env: value for key %q exceeds max length %d", e.Key, opts.MaxValueLen)
			}
			continue
		}
		out = append(out, e)
	}

	if opts.MaxKeys > 0 && len(out) > opts.MaxKeys {
		if opts.StrictKeys {
			return nil, fmt.Errorf("env: snapshot has %d keys, exceeds max %d", len(out), opts.MaxKeys)
		}
		out = out[:opts.MaxKeys]
	}

	if out == nil {
		return nil, errors.New("env: LimitSnapshot produced nil slice")
	}

	return out, nil
}
