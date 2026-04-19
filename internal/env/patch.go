package env

// PatchOptions controls the behaviour of PatchSnapshot.
type PatchOptions struct {
	// SkipEmpty causes entries with an empty value in the patch to be ignored
	// rather than overwriting the base value.
	SkipEmpty bool

	// Remove lists keys that should be deleted from the resulting snapshot.
	Remove []string
}

// DefaultPatchOptions returns conservative defaults.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		SkipEmpty: false,
	}
}

// PatchSnapshot applies a set of key/value overrides (patch) on top of base
// and removes any keys listed in opts.Remove. The base snapshot is never
// mutated; a new Snapshot is always returned.
//
// Precedence: patch wins over base unless opts.SkipEmpty is true and the
// patch value is the empty string.
func PatchSnapshot(base, patch Snapshot, opts PatchOptions) Snapshot {
	out := base.Clone()

	for k, v := range patch {
		if opts.SkipEmpty && v == "" {
			continue
		}
		out[k] = v
	}

	for _, k := range opts.Remove {
		delete(out, k)
	}

	return out
}
