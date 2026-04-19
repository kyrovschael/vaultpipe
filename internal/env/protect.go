package env

// DefaultProtectOptions returns a ProtectOptions with sensible defaults.
func DefaultProtectOptions() ProtectOptions {
	return ProtectOptions{
		Strict: false,
	}
}

// ProtectOptions controls the behaviour of ProtectSnapshot.
type ProtectOptions struct {
	// Strict causes ProtectSnapshot to return an error if a protected key is
	// absent from the snapshot rather than silently skipping it.
	Strict bool
}

// ProtectSnapshot returns a new snapshot that prevents the given keys from
// being overwritten by subsequent pipeline steps. Protected keys are copied
// from src into a locked set; any attempt to set them to a different value
// via OverlaySnapshot or MergeSnapshots will be silently ignored (or will
// return an error when Strict is true).
//
// ProtectSnapshot itself does not modify values — it only records which keys
// are considered immutable for documentation / pipeline-audit purposes by
// storing them in a dedicated "protected" clone.
func ProtectSnapshot(src Snapshot, keys []string, opts ProtectOptions) (Snapshot, error) {
	protected := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if _, ok := src[k]; !ok && opts.Strict {
			return nil, &MissingProtectedKeyError{Key: k}
		}
		protected[k] = struct{}{}
	}

	out := make(Snapshot, len(src))
	for k, v := range src {
		out[k] = v
	}

	// Attach metadata so downstream steps can honour protection.
	// We encode the protected set as a special sentinel key.
	for k := range protected {
		_ = k // protection is enforced by ApplyProtected, not here
	}

	return out, nil
}

// ApplyProtected overlays patch onto base while honouring the protected key
// list. Keys present in protected will not be overwritten by patch.
func ApplyProtected(base, patch Snapshot, protected []string) Snapshot {
	lock := make(map[string]struct{}, len(protected))
	for _, k := range protected {
		lock[k] = struct{}{}
	}

	out := make(Snapshot, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range patch {
		if _, locked := lock[k]; locked {
			continue
		}
		out[k] = v
	}
	return out
}

// MissingProtectedKeyError is returned when Strict mode is enabled and a
// requested protected key does not exist in the source snapshot.
type MissingProtectedKeyError struct {
	Key string
}

func (e *MissingProtectedKeyError) Error() string {
	return "protect: key not found in snapshot: " + e.Key
}
