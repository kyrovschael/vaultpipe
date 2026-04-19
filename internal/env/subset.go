package env

// SubsetOptions controls how SubsetSnapshot behaves.
type SubsetOptions struct {
	// Keys is the explicit set of keys to retain.
	Keys []string
	// IgnoreMissing silently skips keys that are absent from the snapshot.
	IgnoreMissing bool
}

// DefaultSubsetOptions returns sensible defaults.
func DefaultSubsetOptions() SubsetOptions {
	return SubsetOptions{
		IgnoreMissing: true,
	}
}

// SubsetSnapshot returns a new Snapshot containing only the specified keys.
// If IgnoreMissing is false, an error is returned when a requested key is
// absent from src.
func SubsetSnapshot(src Snapshot, opts SubsetOptions) (Snapshot, error) {
	out := make(Snapshot, len(opts.Keys))
	for _, k := range opts.Keys {
		v, ok := src[k]
		if !ok {
			if opts.IgnoreMissing {
				continue
			}
			return nil, &MissingKeyError{Key: k}
		}
		out[k] = v
	}
	return out, nil
}

// MissingKeyError is returned when a required key is absent.
type MissingKeyError struct {
	Key string
}

func (e *MissingKeyError) Error() string {
	return "env: required key not found: " + e.Key
}
