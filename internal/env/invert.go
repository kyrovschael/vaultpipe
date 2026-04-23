package env

// DefaultInvertOptions returns a safe default InvertOptions.
func DefaultInvertOptions() InvertOptions {
	return InvertOptions{
		SkipDuplicateValues: true,
	}
}

// InvertOptions controls the behaviour of InvertSnapshot.
type InvertOptions struct {
	// SkipDuplicateValues drops entries whose value already appears as a key
	// in the inverted map, keeping only the first occurrence (sorted by key
	// for determinism).
	SkipDuplicateValues bool
}

// InvertSnapshot returns a new Snapshot in which every key becomes a value
// and every value becomes a key.
//
// When SkipDuplicateValues is true (the default) and multiple source keys
// share the same value, only the lexicographically smallest source key is
// kept as the new value. When SkipDuplicateValues is false the last
// encountered source key wins (iteration order is undefined for maps, so
// Sort the input first for a stable result).
//
// Entries whose value is an empty string are silently dropped because an
// empty string is not a valid environment variable name.
func InvertSnapshot(src Snapshot, opts InvertOptions) Snapshot {
	out := make(Snapshot, len(src))

	for k, v := range src {
		if v == "" {
			continue
		}
		if opts.SkipDuplicateValues {
			if existing, ok := out[v]; ok {
				// keep the lexicographically smaller source key
				if k >= existing {
					continue
				}
			}
		}
		out[v] = k
	}

	return out
}
