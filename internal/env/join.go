package env

import "strings"

// DefaultJoinOptions returns a JoinOptions with sensible defaults.
func DefaultJoinOptions() JoinOptions {
	return JoinOptions{
		Separator: ",",
	}
}

// JoinOptions controls how JoinSnapshot combines values.
type JoinOptions struct {
	// Separator is placed between values when multiple snapshots define the
	// same key. Defaults to ",".
	Separator string

	// Keys restricts joining to the listed keys. When empty, all keys that
	// appear in more than one snapshot are joined.
	Keys []string
}

// JoinSnapshot merges multiple snapshots into one. Where a key exists in more
// than one snapshot its values are concatenated with opts.Separator. Keys that
// appear in only one snapshot are copied as-is. Later snapshots take
// precedence for keys that are NOT joined (i.e. not in opts.Keys when the
// list is non-empty).
func JoinSnapshot(snapshots []Snapshot, opts JoinOptions) Snapshot {
	result := Snapshot{}

	// Collect all keys in encounter order so output is deterministic.
	seen := map[string]bool{}
	var order []string
	for _, s := range snapshots {
		for _, e := range s {
			if !seen[e.Key] {
				seen[e.Key] = true
				order = append(order, e.Key)
			}
		}
	}

	joinSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		joinSet[k] = true
	}

	for _, key := range order {
		var parts []string
		for _, s := range snapshots {
			for _, e := range s {
				if e.Key == key {
					parts = append(parts, e.Value)
				}
			}
		}

		shouldJoin := len(joinSet) == 0 || joinSet[key]
		var val string
		if shouldJoin && len(parts) > 1 {
			val = strings.Join(parts, opts.Separator)
		} else {
			val = parts[len(parts)-1]
		}
		result = append(result, Entry{Key: key, Value: val})
	}

	return result
}
