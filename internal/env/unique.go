package env

// DefaultUniqueOptions returns a UniqueOptions with sensible defaults.
func DefaultUniqueOptions() UniqueOptions {
	return UniqueOptions{
		CaseFold: false,
	}
}

// UniqueOptions controls the behaviour of UniqueSnapshot.
type UniqueOptions struct {
	// CaseFold performs case-insensitive value comparison when true.
	CaseFold bool

	// Keys restricts deduplication to only the listed keys.
	// An empty slice means all keys are considered.
	Keys []string
}

// UniqueSnapshot returns a new Snapshot containing only entries whose values
// have not been seen before. When multiple keys share the same value the key
// that appears first in the sorted snapshot order is retained.
//
// This is distinct from DedupeSnapshot, which operates on duplicate keys;
// UniqueSnapshot operates on duplicate values.
func UniqueSnapshot(s Snapshot, opts UniqueOptions) (Snapshot, error) {
	seen := make(map[string]struct{})
	restricted := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		restricted[k] = struct{}{}
	}

	out := make(Snapshot)
	for k, v := range s {
		if len(restricted) > 0 {
			if _, ok := restricted[k]; !ok {
				out[k] = v
				continue
			}
		}

		norm := v
		if opts.CaseFold {
			norm = toLower(v)
		}

		if _, exists := seen[norm]; exists {
			continue
		}
		seen[norm] = struct{}{}
		out[k] = v
	}
	return out, nil
}

// toLower is a local helper to avoid importing strings in a hot path.
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
