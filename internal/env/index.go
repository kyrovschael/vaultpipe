package env

import "strings"

// Index is an inverted mapping from value → keys that hold that value.
type Index map[string][]string

// Lookup returns the list of keys whose value equals v.
func (idx Index) Lookup(v string) ([]string, bool) {
	keys, ok := idx[v]
	return keys, ok
}

// IndexOptions controls the behaviour of IndexSnapshot.
type IndexOptions struct {
	// Keys restricts indexing to the listed keys. An empty slice means all keys.
	Keys []string
	// CaseFold normalises both keys and values to lower-case before indexing.
	CaseFold bool
}

// DefaultIndexOptions returns a sensible default configuration.
func DefaultIndexOptions() IndexOptions {
	return IndexOptions{}
}

// IndexSnapshot builds an inverted index from the values in snap to the keys
// that contain them. When multiple keys share the same value they are all
// recorded, sorted in the order they appear in the snapshot.
func IndexSnapshot(snap Snapshot, opts IndexOptions) (Index, error) {
	allowed := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		if opts.CaseFold {
			k = strings.ToLower(k)
		}
		allowed[k] = struct{}{}
	}

	idx := make(Index)

	for _, entry := range snap.ToSlice() {
		k, v, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		if opts.CaseFold {
			k = strings.ToLower(k)
			v = strings.ToLower(v)
		}
		if len(allowed) > 0 {
			if _, want := allowed[k]; !want {
				continue
			}
		}
		idx[v] = append(idx[v], k)
	}

	return idx, nil
}
