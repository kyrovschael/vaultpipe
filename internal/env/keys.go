package env

import (
	"sort"
	"strings"
)

// KeysOptions controls the behaviour of KeysSnapshot.
type KeysOptions struct {
	// Prefix restricts output to keys that start with the given prefix.
	// An empty string means all keys are included.
	Prefix string

	// CaseFold performs case-insensitive prefix matching when true.
	CaseFold bool

	// Sorted returns keys in lexicographic ascending order when true.
	Sorted bool
}

// DefaultKeysOptions returns a KeysOptions with sensible defaults.
func DefaultKeysOptions() KeysOptions {
	return KeysOptions{
		Prefix:   "",
		CaseFold: false,
		Sorted:   true,
	}
}

// KeysSnapshot returns all keys present in src as a string slice.
// Options control prefix filtering, case folding, and ordering.
func KeysSnapshot(src Snapshot, opts KeysOptions) []string {
	keys := make([]string, 0, len(src))

	for k := range src {
		if opts.Prefix != "" {
			if opts.CaseFold {
				if !strings.HasPrefix(strings.ToLower(k), strings.ToLower(opts.Prefix)) {
					continue
				}
			} else {
				if !strings.HasPrefix(k, opts.Prefix) {
					continue
				}
			}
		}
		keys = append(keys, k)
	}

	if opts.Sorted {
		sort.Strings(keys)
	}

	return keys
}
