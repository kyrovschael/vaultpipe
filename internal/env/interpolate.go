package env

import (
	"fmt"
	"strings"
)

// InterpolateOptions controls how secret placeholders are resolved in values.
type InterpolateOptions struct {
	// Delimiters wrapping the placeholder key, e.g. "${" and "}".
	Open  string
	Close string
	// Strict causes an error when a placeholder key is not found in secrets.
	Strict bool
}

// DefaultInterpolateOptions returns sensible defaults using ${KEY} syntax.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		Open:   "${",
		Close:  "}",
		Strict: false,
	}
}

// InterpolateSnapshot replaces placeholders in snapshot values with values
// drawn from secrets. The snapshot itself is not mutated; a new one is returned.
func InterpolateSnapshot(snap Snapshot, secrets Snapshot, opts InterpolateOptions) (Snapshot, error) {
	out := make(Snapshot, len(snap))
	for k, v := range snap {
		replaced, err := interpolateValue(v, secrets, opts)
		if err != nil {
			return nil, fmt.Errorf("env: interpolate key %q: %w", k, err)
		}
		out[k] = replaced
	}
	return out, nil
}

func interpolateValue(v string, secrets Snapshot, opts InterpolateOptions) (string, error) {
	var sb strings.Builder
	s := v
	for {
		start := strings.Index(s, opts.Open)
		if start == -1 {
			sb.WriteString(s)
			break
		}
		sb.WriteString(s[:start])
		rest := s[start+len(opts.Open):]
		end := strings.Index(rest, opts.Close)
		if end == -1 {
			// No closing delimiter — treat remainder as literal.
			sb.WriteString(opts.Open)
			sb.WriteString(rest)
			break
		}
		key := rest[:end]
		val, ok := secrets[key]
		if !ok {
			if opts.Strict {
				return "", fmt.Errorf("placeholder %q not found in secrets", key)
			}
			// Leave placeholder intact.
			sb.WriteString(opts.Open)
			sb.WriteString(key)
			sb.WriteString(opts.Close)
		} else {
			sb.WriteString(val)
		}
		s = rest[end+len(opts.Close):]
	}
	return sb.String(), nil
}
