package env

import "strings"

// DefaultClampOptions returns a ClampOptions with safe zero-value defaults.
func DefaultClampOptions() ClampOptions {
	return ClampOptions{
		PadChar: ' ',
	}
}

// ClampOptions controls the behaviour of ClampSnapshot.
type ClampOptions struct {
	// Keys restricts clamping to the listed keys. Empty means all keys.
	Keys []string
	// MinLen pads values shorter than this length. 0 disables padding.
	MinLen int
	// MaxLen truncates values longer than this length. 0 disables truncation.
	MaxLen int
	// PadChar is the character used when padding short values (default ' ').
	PadChar rune
}

// ClampSnapshot returns a new Snapshot where string values are constrained to
// [MinLen, MaxLen].  The original snapshot is never mutated.
func ClampSnapshot(src Snapshot, opts ClampOptions) Snapshot {
	if opts.PadChar == 0 {
		opts.PadChar = ' '
	}

	restricted := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		restricted[k] = struct{}{}
	}

	wantClamp := func(key string) bool {
		if len(restricted) == 0 {
			return true
		}
		_, ok := restricted[key]
		return ok
	}

	out := make(Snapshot, len(src))
	for k, v := range src {
		if !wantClamp(k) {
			out[k] = v
			continue
		}
		out[k] = clampValue(v, opts)
	}
	return out
}

func clampValue(v string, opts ClampOptions) string {
	if opts.MaxLen > 0 && len(v) > opts.MaxLen {
		v = v[:opts.MaxLen]
	}
	if opts.MinLen > 0 && len(v) < opts.MinLen {
		v += strings.Repeat(string(opts.PadChar), opts.MinLen-len(v))
	}
	return v
}
