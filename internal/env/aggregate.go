package env

import "sort"

// AggregateOptions controls how AggregateSnapshot combines values.
type AggregateOptions struct {
	// Separator is placed between values when joining. Defaults to ",".
	Separator string
	// Keys restricts aggregation to the given keys. Empty means all keys.
	Keys []string
	// SortValues sorts the individual values before joining.
	SortValues bool
}

// DefaultAggregateOptions returns sensible defaults for AggregateSnapshot.
func DefaultAggregateOptions() AggregateOptions {
	return AggregateOptions{
		Separator:  ",",
		SortValues: false,
	}
}

// AggregateSnapshot merges values for matching keys across multiple snapshots
// into a single snapshot. When the same key appears in more than one snapshot
// the values are joined with opts.Separator rather than one overwriting another.
//
// Snapshots are processed left-to-right; the first snapshot that contains a key
// establishes the canonical key casing in the result.
func AggregateSnapshot(snapshots []Snapshot, opts AggregateOptions) (Snapshot, error) {
	if opts.Separator == "" {
		opts.Separator = ","
	}

	allowedKeys := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		allowedKeys[k] = struct{}{}
	}

	// collected maps canonical key → slice of values
	collected := make(map[string][]string)
	// order preserves first-seen insertion order
	var order []string

	for _, snap := range snapshots {
		for k, v := range snap {
			if len(allowedKeys) > 0 {
				if _, ok := allowedKeys[k]; !ok {
					continue
				}
			}
			if _, seen := collected[k]; !seen {
				order = append(order, k)
			}
			collected[k] = append(collected[k], v)
		}
	}

	result := make(Snapshot, len(collected))
	for _, k := range order {
		vals := collected[k]
		if opts.SortValues {
			sorted := make([]string, len(vals))
			copy(sorted, vals)
			sort.Strings(sorted)
			vals = sorted
		}
		result[k] = joinStrings(vals, opts.Separator)
	}
	return result, nil
}

func joinStrings(ss []string, sep string) string {
	switch len(ss) {
	case 0:
		return ""
	case 1:
		return ss[0]
	}
	n := len(sep) * (len(ss) - 1)
	for _, s := range ss {
		n += len(s)
	}
	b := make([]byte, 0, n)
	for i, s := range ss {
		if i > 0 {
			b = append(b, sep...)
		}
		b = append(b, s...)
	}
	return string(b)
}
