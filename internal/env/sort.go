package env

import (
	"sort"
)

// SortOrder controls the direction of key sorting.
type SortOrder int

const (
	SortAsc  SortOrder = iota // ascending (default)
	SortDesc                  // descending
)

// SortOptions configures SortSnapshot behaviour.
type SortOptions struct {
	Order SortOrder
	// KeyFn transforms a key before comparison (e.g. strings.ToLower).
	// If nil, the raw key is used.
	KeyFn func(string) string
}

// DefaultSortOptions returns options that sort keys in ascending order.
func DefaultSortOptions() SortOptions {
	return SortOptions{Order: SortAsc}
}

// SortSnapshot returns a new Snapshot whose keys are ordered according to opts.
// The original snapshot is not mutated.
func SortSnapshot(s Snapshot, opts SortOptions) Snapshot {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}

	cmp := func(a, b string) bool {
		ka, kb := a, b
		if opts.KeyFn != nil {
			ka = opts.KeyFn(a)
			kb = opts.KeyFn(b)
		}
		if opts.Order == SortDesc {
			return ka > kb
		}
		return ka < kb
	}

	sort.Slice(keys, func(i, j int) bool {
		return cmp(keys[i], keys[j])
	})

	out := make(Snapshot, len(s))
	for _, k := range keys {
		out[k] = s[k]
	}
	return out
}
