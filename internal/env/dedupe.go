package env

// DedupeOptions controls how duplicate keys are resolved when
// flattening a slice of snapshots into one.
type DedupeOptions struct {
	// LastWins, when true, keeps the last occurrence of a key;
	// when false (default) the first occurrence is kept.
	LastWins bool
}

// DefaultDedupeOptions returns conservative defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{LastWins: false}
}

// DedupeSnapshot removes duplicate keys from s according to opts.
// The original snapshot is not mutated; a new one is returned.
func DedupeSnapshot(s Snapshot, opts DedupeOptions) Snapshot {
	out := make(Snapshot, len(s))
	seen := make(map[string]int, len(s)) // key -> index in out
	idx := 0

	for _, e := range s {
		if pos, exists := seen[e.Key]; exists {
			if opts.LastWins {
				out[pos] = e
			}
			continue
		}
		out[idx] = e
		seen[e.Key] = idx
		idx++
	}

	return out[:idx]
}

// DedupeSlice deduplicates a raw KEY=VALUE slice and returns a clean slice.
func DedupeSlice(pairs []string, opts DedupeOptions) ([]string, error) {
	s, err := ParseSlice(pairs, DefaultParseOptions())
	if err != nil {
		return nil, err
	}
	deduped := DedupeSnapshot(s, opts)
	return deduped.ToSlice(), nil
}
