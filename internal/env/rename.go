package env

// RenameOptions controls the behaviour of RenameSnapshot.
type RenameOptions struct {
	// DropUnmapped drops keys that are not present in the mapping.
	// When false, unmapped keys are passed through unchanged.
	DropUnmapped bool
}

// DefaultRenameOptions returns the default RenameOptions.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{DropUnmapped: false}
}

// RenameSnapshot returns a new Snapshot with keys renamed according to the
// provided mapping (oldKey → newKey). Unmapped keys are kept or dropped
// depending on opts.DropUnmapped. If two old keys map to the same new key the
// last one (in iteration order of the mapping) wins; callers should avoid
// ambiguous mappings.
func RenameSnapshot(s Snapshot, mapping map[string]string, opts RenameOptions) Snapshot {
	out := make(Snapshot, len(s))

	// Build reverse lookup so we can detect which originals were renamed.
	renamed := make(map[string]struct{}, len(mapping))
	for old := range mapping {
		renamed[old] = struct{}{}
	}

	for k, v := range s {
		if newKey, ok := mapping[k]; ok {
			out[newKey] = v
			continue
		}
		if !opts.DropUnmapped {
			out[k] = v
		}
	}
	return out
}
