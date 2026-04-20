package env

import (
	"context"
	"sort"
	"strings"
)

// GroupOptions controls how GroupSnapshot partitions a snapshot.
type GroupOptions struct {
	// KeyFn derives the group name from a key. Defaults to prefix before
	// the first underscore (e.g. "DB_HOST" -> "DB").
	KeyFn func(key string) string
}

// DefaultGroupOptions returns GroupOptions with the default prefix splitter.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		KeyFn: func(key string) string {
			if idx := strings.IndexByte(key, '_'); idx > 0 {
				return key[:idx]
			}
			return key
		},
	}
}

// GroupSnapshot partitions the snapshot returned by src into named
// sub-snapshots according to opts.KeyFn. Each returned map value is
// an independent clone; mutations do not affect each other.
func GroupSnapshot(ctx context.Context, src Source, opts GroupOptions) (map[string]Snapshot, error) {
	snap, err := src.Load(ctx)
	if err != nil {
		return nil, err
	}

	if opts.KeyFn == nil {
		opts = DefaultGroupOptions()
	}

	groups := make(map[string]Snapshot)
	for _, e := range snap {
		group := opts.KeyFn(e.Key)
		if group == "" {
			group = "_"
		}
		if groups[group] == nil {
			groups[group] = make(Snapshot)
		}
		groups[group][e.Key] = e.Value
	}
	return groups, nil
}

// GroupNames returns the sorted group names from a grouped snapshot map.
func GroupNames(groups map[string]Snapshot) []string {
	names := make([]string, 0, len(groups))
	for k := range groups {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
