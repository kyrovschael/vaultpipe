package env

// CloneOptions controls the behaviour of CloneSnapshot.
type CloneOptions struct {
	// Keys limits the clone to only the listed keys.
	// When empty all keys are copied.
	Keys []string

	// Deep forces value strings to be copied via a string builder so that
	// no underlying memory is shared with the source snapshot.
	Deep bool
}

// DefaultCloneOptions returns a CloneOptions that copies every key shallowly.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{}
}

// CloneSnapshot returns a new Snapshot that is an independent copy of src.
//
// When opts.Keys is non-empty only those keys are included in the result;
// keys that are absent from src are silently skipped.
//
// When opts.Deep is true each value is rebuilt character-by-character so that
// the returned snapshot shares no backing memory with src.
func CloneSnapshot(src Snapshot, opts CloneOptions) Snapshot {
	wantAll := len(opts.Keys) == 0

	var keys []string
	if wantAll {
		keys = make([]string, 0, len(src))
		for k := range src {
			keys = append(keys, k)
		}
	} else {
		keys = opts.Keys
	}

	out := make(Snapshot, len(keys))
	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		if opts.Deep {
			buf := make([]byte, len(v))
			copy(buf, v)
			v = string(buf)
		}
		out[k] = v
	}
	return out
}
