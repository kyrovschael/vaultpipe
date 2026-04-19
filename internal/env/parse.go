package env

import (
	"fmt"
	"strings"
)

// ParseOptions controls how raw environment strings are parsed.
type ParseOptions struct {
	// TrimSpace strips leading/trailing whitespace from keys and values.
	TrimSpace bool
	// SkipInvalid silently drops malformed entries instead of returning an error.
	SkipInvalid bool
}

// DefaultParseOptions returns sensible defaults.
func DefaultParseOptions() ParseOptions {
	return ParseOptions{
		TrimSpace:   true,
		SkipInvalid: false,
	}
}

// ParseSlice parses a slice of "KEY=VALUE" strings into a Snapshot.
// Entries without '=' are treated as keys with empty values when SkipInvalid
// is false, or dropped when SkipInvalid is true.
func ParseSlice(entries []string, opts ParseOptions) (Snapshot, error) {
	out := make(Snapshot, len(entries))
	for _, e := range entries {
		key, val, err := parsePair(e, opts.TrimSpace)
		if err != nil {
			if opts.SkipInvalid {
				continue
			}
			return nil, err
		}
		out[key] = val
	}
	return out, nil
}

// ParseMap parses a map[string]string applying optional trimming and key
// validation, returning a validated Snapshot.
func ParseMap(m map[string]string, opts ParseOptions) (Snapshot, error) {
	out := make(Snapshot, len(m))
	for k, v := range m {
		if opts.TrimSpace {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
		}
		if err := ValidateKey(k); err != nil {
			if opts.SkipInvalid {
				continue
			}
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}

func parsePair(entry string, trim bool) (string, string, error) {
	idx := strings.IndexByte(entry, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("env: missing '=' in entry %q", entry)
	}
	k := entry[:idx]
	v := entry[idx+1:]
	if trim {
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
	}
	if err := ValidateKey(k); err != nil {
		return "", "", err
	}
	return k, v, nil
}
