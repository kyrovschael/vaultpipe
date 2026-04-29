package env

import (
	"strings"
)

// DefaultWindowOptions returns WindowOptions with no bounds and case-sensitive
// comparison.
func DefaultWindowOptions() WindowOptions {
	return WindowOptions{}
}

// WindowOptions controls the behaviour of WindowSnapshot.
type WindowOptions struct {
	// From is the inclusive lower bound key. Empty means no lower bound.
	From string
	// To is the inclusive upper bound key. Empty means no upper bound.
	To string
	// CaseFold performs case-insensitive key comparison when true.
	CaseFold bool
}

// WindowSnapshot returns a new snapshot containing only entries whose keys
// fall within the [From, To] range (inclusive). Either bound may be left
// empty to indicate an open range on that side.
func WindowSnapshot(s Snapshot, opts WindowOptions) Snapshot {
	out := make(Snapshot, 0, len(s))
	for _, e := range s {
		key := e.Key
		if opts.CaseFold {
			key = strings.ToLower(key)
		}
		if opts.From != "" {
			lo := opts.From
			if opts.CaseFold {
				lo = strings.ToLower(lo)
			}
			if key < lo {
				continue
			}
		}
		if opts.To != "" {
			hi := opts.To
			if opts.CaseFold {
				hi = strings.ToLower(hi)
			}
			if key > hi {
				continue
			}
		}
		out = append(out, e)
	}
	return out
}
