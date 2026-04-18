package env

import (
	"os"
	"strings"
)

// ExpandOptions controls how variable expansion is performed.
type ExpandOptions struct {
	// FallbackToOS allows falling back to os.Getenv when a key is not in the snapshot.
	FallbackToOS bool
	// NoExpand contains keys whose values should not be expanded.
	NoExpand map[string]bool
}

// DefaultExpandOptions returns sensible defaults for expansion.
func DefaultExpandOptions() ExpandOptions {
	return ExpandOptions{
		FallbackToOS: false,
		NoExpand:     map[string]bool{},
	}
}

// Expand returns a new Snapshot where each value has $VAR and ${VAR}
// references resolved against the snapshot itself, with optional OS fallback.
func Expand(snap Snapshot, opts ExpandOptions) Snapshot {
	out := make(Snapshot, len(snap))

	mapping := func(key string) string {
		if opts.NoExpand[key] {
			return "$" + key
		}
		if v, ok := snap[key]; ok {
			return v
		}
		if opts.FallbackToOS {
			return os.Getenv(key)
		}
		return ""
	}

	for k, v := range snap {
		if opts.NoExpand[k] {
			out[k] = v
			continue
		}
		out[k] = strings.Expand(v, mapping)
	}
	return out
}
