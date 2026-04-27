package env

import "fmt"

// DefaultResolveOptions returns a ResolveOptions with sensible defaults.
func DefaultResolveOptions() ResolveOptions {
	return ResolveOptions{
		Strict: false,
	}
}

// ResolveOptions controls the behaviour of ResolveSnapshot.
type ResolveOptions struct {
	// Strict causes ResolveSnapshot to return an error when a key cannot be
	// found in any source. When false, missing keys are silently skipped.
	Strict bool
}

// ResolveSnapshot resolves each key in keys from the provided sources in order.
// The first source that returns a non-empty value for a key wins. Subsequent
// sources are not consulted for that key.
//
// Sources are queried via their Snapshot method. An error from a source causes
// ResolveSnapshot to return immediately.
func ResolveSnapshot(keys []string, opts ResolveOptions, sources ...Source) (Snapshot, error) {
	if len(sources) == 0 {
		if opts.Strict && len(keys) > 0 {
			return Snapshot{}, fmt.Errorf("resolve: no sources provided, cannot resolve %d key(s)", len(keys))
		}
		return Snapshot{}, nil
	}

	// Collect all source snapshots up front so we only call each source once.
	snapshots := make([]Snapshot, 0, len(sources))
	for _, src := range sources {
		snap, err := src.Snapshot()
		if err != nil {
			return Snapshot{}, fmt.Errorf("resolve: source error: %w", err)
		}
		snapshots = append(snapshots, snap)
	}

	out := make(map[string]string, len(keys))
	for _, key := range keys {
		resolved := false
		for _, snap := range snapshots {
			if val, ok := snap.data[key]; ok && val != "" {
				out[key] = val
			resolved = true
				break
			}
		}
		if !resolved && opts.Strict {
			return Snapshot{}, fmt.Errorf("resolve: key %q not found in any source", key)
		}
	}

	return Snapshot{data: out}, nil
}
