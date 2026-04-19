package env

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// SourceFunc is a function that returns a Snapshot or an error.
type SourceFunc func(ctx context.Context) (Snapshot, error)

// OSSource returns a SourceFunc that reads the current process environment.
func OSSource() SourceFunc {
	return func(_ context.Context) (Snapshot, error) {
		return FromSlice(os.Environ()), nil
	}
}

// StaticSource returns a SourceFunc that always returns the given snapshot.
func StaticSource(s Snapshot) SourceFunc {
	return func(_ context.Context) (Snapshot, error) {
		return s.Clone(), nil
	}
}

// SliceSource returns a SourceFunc backed by a raw KEY=VALUE slice.
func SliceSource(pairs []string) SourceFunc {
	return func(_ context.Context) (Snapshot, error) {
		return FromSlice(pairs), nil
	}
}

// ChainSource returns a SourceFunc that merges multiple sources left-to-right,
// with later sources winning on key conflicts.
func ChainSource(sources ...SourceFunc) SourceFunc {
	return func(ctx context.Context) (Snapshot, error) {
		result := Snapshot{}
		for i, src := range sources {
			s, err := src(ctx)
			if err != nil {
				return nil, fmt.Errorf("source %d: %w", i, err)
			}
			for k, v := range s {
				result[k] = v
			}
		}
		return result, nil
	}
}

// MapSource returns a SourceFunc backed by a plain map.
func MapSource(m map[string]string) SourceFunc {
	return func(_ context.Context) (Snapshot, error) {
		s := make(Snapshot, len(m))
		for k, v := range m {
			s[k] = v
		}
		return s, nil
	}
}

// PrefixedOSSource reads OS env and returns only keys with the given prefix,
// stripping the prefix from the returned keys.
func PrefixedOSSource(prefix string) SourceFunc {
	return func(_ context.Context) (Snapshot, error) {
		s := make(Snapshot)
		for _, pair := range os.Environ() {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) != 2 {
				continue
			}
			if strings.HasPrefix(parts[0], prefix) {
				s[strings.TrimPrefix(parts[0], prefix)] = parts[1]
			}
		}
		return s, nil
	}
}
