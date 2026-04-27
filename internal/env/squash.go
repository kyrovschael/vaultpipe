package env

// MergeFn resolves a conflict between two values for the same key.
// It receives the key, the current accumulated value, and the incoming
// value, and returns the value that should be stored.
type MergeFn func(key, current, incoming string) string

// SquashOptions controls the behaviour of SquashSnapshot.
type SquashOptions struct {
	// MergeFn is called when two snapshots share a key.
	// When nil the incoming (right-hand) value wins.
	MergeFn MergeFn

	// Keys restricts squashing to the listed keys.
	// When empty all keys are processed.
	Keys []string
}

// DefaultSquashOptions returns a SquashOptions where the last value wins.
func DefaultSquashOptions() SquashOptions {
	return SquashOptions{}
}

// SquashSnapshot merges one or more snapshots into a single Snapshot.
// Snapshots are processed left-to-right; conflicts are resolved by
// opts.MergeFn (defaults to last-wins when nil).
func SquashSnapshot(opts SquashOptions, snapshots ...Snapshot) (Snapshot, error) {
	if len(snapshots) == 0 {
		return Snapshot{}, nil
	}

	allowed := make(map[string]struct{}, len(opts.Keys))
	restricted := len(opts.Keys) > 0
	for _, k := range opts.Keys {
		allowed[k] = struct{}{}
	}

	mergeFn := opts.MergeFn
	if mergeFn == nil {
		mergeFn = func(_, _, incoming string) string { return incoming }
	}

	acc := make(map[string]string)

	for _, snap := range snapshots {
		for k, v := range snap.m {
			if restricted {
				if _, ok := allowed[k]; !ok {
					continue
				}
			}
			if cur, exists := acc[k]; exists {
				acc[k] = mergeFn(k, cur, v)
			} else {
				acc[k] = v
			}
		}
	}

	return FromMap(acc), nil
}
