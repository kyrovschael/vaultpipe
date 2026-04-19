package env

import "strings"

// FlattenOptions controls how nested map keys are joined.
type FlattenOptions struct {
	// Separator is placed between key segments. Defaults to "_".
	Separator string
	// UpperKeys uppercases every resulting key.
	UpperKeys bool
}

// DefaultFlattenOptions returns sensible defaults.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "_",
		UpperKeys: true,
	}
}

// FlattenMap converts a nested map[string]any into a Snapshot by joining
// key segments with the configured separator. Only leaf string values are
// kept; non-string leaves are skipped.
func FlattenMap(src map[string]any, opts FlattenOptions) Snapshot {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	out := make(Snapshot)
	flattenRecurse(src, "", opts, out)
	return out
}

func flattenRecurse(m map[string]any, prefix string, opts FlattenOptions, out Snapshot) {
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + opts.Separator + k
		}
		if opts.UpperKeys {
			full = strings.ToUpper(full)
		}
		switch val := v.(type) {
		case map[string]any:
			flattenRecurse(val, full, opts, out)
		case string:
			out[full] = val
		}
	}
}
