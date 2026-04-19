package env

import "strings"

// MaskOptions controls how MaskSnapshot redacts values.
type MaskOptions struct {
	// Replacement is the string used in place of a masked value. Defaults to "***".
	Replacement string
	// Keys lists the exact keys whose values should be masked.
	Keys []string
}

// DefaultMaskOptions returns sensible defaults for MaskOptions.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Replacement: "***",
	}
}

// MaskSnapshot returns a new Snapshot where every key listed in opts.Keys has
// its value replaced with opts.Replacement. The original Snapshot is not
// mutated.
func MaskSnapshot(s Snapshot, opts MaskOptions) Snapshot {
	if opts.Replacement == "" {
		opts.Replacement = "***"
	}
	masked := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		masked[strings.ToUpper(k)] = true
	}
	out := s.Clone()
	for i, e := range out.entries {
		if masked[strings.ToUpper(e.Key)] {
			out.entries[i].Value = opts.Replacement
		}
	}
	return out
}
