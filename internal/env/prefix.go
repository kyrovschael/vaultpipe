package env

import "strings"

// PrefixOptions controls how prefix stripping and filtering behaves.
type PrefixOptions struct {
	// StripPrefix removes the prefix from keys in the resulting snapshot.
	StripPrefix bool
}

// DefaultPrefixOptions returns sensible defaults.
func DefaultPrefixOptions() PrefixOptions {
	return PrefixOptions{StripPrefix: true}
}

// FilterByPrefix returns a new Snapshot containing only keys that start with
// the given prefix. If opts.StripPrefix is true the prefix is removed from
// each key in the returned snapshot.
func FilterByPrefix(snap Snapshot, prefix string, opts PrefixOptions) Snapshot {
	if prefix == "" {
		return snap.Clone()
	}

	out := make(Snapshot, len(snap))
	for k, v := range snap {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		key := k
		if opts.StripPrefix {
			key = strings.TrimPrefix(k, prefix)
			if key == "" {
				continue
			}
		}
		out[key] = v
	}
	return out
}

// AddPrefix returns a new Snapshot with prefix prepended to every key.
func AddPrefix(snap Snapshot, prefix string) Snapshot {
	if prefix == "" {
		return snap.Clone()
	}
	out := make(Snapshot, len(snap))
	for k, v := range snap {
		out[prefix+k] = v
	}
	return out
}
