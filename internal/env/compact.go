package env

// DefaultCompactOptions returns a CompactOptions with sensible defaults.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{
		TrimSpace:    true,
		DropEmpty:    true,
		DropWhitespace: true,
	}
}

// CompactOptions controls the behaviour of CompactSnapshot.
type CompactOptions struct {
	// TrimSpace trims leading and trailing whitespace from every value.
	TrimSpace bool

	// DropEmpty removes entries whose value is empty (after optional trimming).
	DropEmpty bool

	// DropWhitespace removes entries whose value is entirely whitespace.
	// Only evaluated when TrimSpace is false; when TrimSpace is true a
	// whitespace-only value becomes empty and is caught by DropEmpty.
	DropWhitespace bool
}

// CompactSnapshot returns a new Snapshot with values cleaned according to
// opts. The original snapshot is never mutated.
func CompactSnapshot(s Snapshot, opts CompactOptions) Snapshot {
	out := make(Snapshot, len(s))
	for k, v := range s {
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}
		if opts.DropEmpty && v == "" {
			continue
		}
		if !opts.TrimSpace && opts.DropWhitespace && strings.TrimSpace(v) == "" {
			continue
		}
		out[k] = v
	}
	return out
}
