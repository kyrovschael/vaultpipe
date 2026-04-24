package env

// DefaultOmitOptions returns an OmitOptions with safe defaults.
func DefaultOmitOptions() OmitOptions {
	return OmitOptions{
		Keys:        nil,
		CaseFold:    false,
		PrefixMatch: false,
	}
}

// OmitOptions controls the behaviour of OmitSnapshot.
type OmitOptions struct {
	// Keys is the set of keys (or prefixes) to remove from the snapshot.
	Keys []string

	// CaseFold performs case-insensitive key comparison when true.
	CaseFold bool

	// PrefixMatch removes any key whose name starts with one of the listed
	// Keys rather than requiring an exact match.
	PrefixMatch bool
}

// OmitSnapshot returns a clone of src with the specified keys removed.
// The original snapshot is never mutated.
func OmitSnapshot(src Snapshot, opts OmitOptions) Snapshot {
	out := make(Snapshot, len(src))
	for k, v := range src {
		if omitKey(k, opts) {
			continue
		}
		out[k] = v
	}
	return out
}

func omitKey(key string, opts OmitOptions) bool {
	cmp := key
	if opts.CaseFold {
		cmp = toUpper(key)
	}
	for _, omit := range opts.Keys {
		pattern := omit
		if opts.CaseFold {
			pattern = toUpper(omit)
		}
		if opts.PrefixMatch {
			if len(cmp) >= len(pattern) && cmp[:len(pattern)] == pattern {
				return true
			}
		} else {
			if cmp == pattern {
				return true
			}
		}
	}
	return false
}
