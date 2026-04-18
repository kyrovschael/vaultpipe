package env

// MergeStrategy controls how conflicting keys are resolved during a merge.
type MergeStrategy int

const (
	// OverlayWins means the overlay value replaces the base value.
	OverlayWins MergeStrategy = iota
	// BaseWins means the base value is preserved when a key exists in both.
	BaseWins
)

// MergeOptions configures the behaviour of MergeSnapshots.
type MergeOptions struct {
	Strategy MergeStrategy
	// DenyList, if non-nil, filters keys from the overlay before merging.
	DenyList *DenyList
}

// DefaultMergeOptions returns options that match the most common use-case:
// overlay (secrets) win over base (process environment).
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{Strategy: OverlayWins}
}

// MergeSnapshots combines base and overlay Snapshots according to opts.
// The original snapshots are not mutated; a new Snapshot is returned.
func MergeSnapshots(base, overlay Snapshot, opts MergeOptions) Snapshot {
	out := make(Snapshot, len(base))
	for k, v := range base {
		out[k] = v
	}

	for k, v := range overlay {
		if opts.DenyList != nil && opts.DenyList.Blocked(k) {
			continue
		}
		if opts.Strategy == BaseWins {
			if _, exists := out[k]; exists {
				continue
			}
		}
		out[k] = v
	}
	return out
}
