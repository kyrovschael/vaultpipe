package env

// TruncateOptions controls how values are truncated within a snapshot.
type TruncateOptions struct {
	// MaxLen is the maximum byte length of any single value.
	// Values longer than this are trimmed and optionally suffixed.
	MaxLen int

	// Suffix is appended to any value that was truncated.
	// Defaults to "..." when empty and a truncation actually occurs.
	Suffix string

	// Keys restricts truncation to the named keys only.
	// When empty, all keys are subject to truncation.
	Keys []string
}

// DefaultTruncateOptions returns conservative defaults: 256-byte limit with
// an ellipsis suffix applied to every key.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxLen: 256,
		Suffix: "...",
	}
}

// TruncateSnapshot returns a new Snapshot where every value (or only those
// matching opts.Keys) is capped at opts.MaxLen bytes. Values that are already
// within the limit are copied verbatim.
func TruncateSnapshot(s Snapshot, opts TruncateOptions) Snapshot {
	if opts.MaxLen <= 0 {
		opts.MaxLen = DefaultTruncateOptions().MaxLen
	}
	suffix := opts.Suffix

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make(Snapshot, len(s))
	for k, v := range s {
		if len(opts.Keys) > 0 {
			if _, ok := keySet[k]; !ok {
				out[k] = v
				continue
			}
		}
		if len(v) > opts.MaxLen {
			truncated := v[:opts.MaxLen]
			if suffix != "" {
				truncated += suffix
			}
			out[k] = truncated
		} else {
			out[k] = v
		}
	}
	return out
}
