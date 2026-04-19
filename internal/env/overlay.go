package env

// DefaultOverlayOptions returns sensible defaults for OverlaySnapshot.
func DefaultOverlayOptions() OverlayOptions {
	return OverlayOptions{
		OverwriteExisting: true,
		SkipEmpty:         false,
	}
}

// OverlayOptions controls how layers are merged in OverlaySnapshot.
type OverlayOptions struct {
	// OverwriteExisting allows overlay values to replace base values.
	OverwriteExisting bool
	// SkipEmpty ignores overlay entries whose value is the empty string.
	SkipEmpty bool
}

// OverlaySnapshot merges one or more overlay snapshots onto base.
// Layers are applied left-to-right; later layers win over earlier ones
// when OverwriteExisting is true.
func OverlaySnapshot(base Snapshot, opts OverlayOptions, layers ...Snapshot) Snapshot {
	out := base.Clone()
	for _, layer := range layers {
		for _, e := range layer.Entries() {
			if opts.SkipEmpty && e.Value == "" {
				continue
			}
			_, exists := out.Lookup(e.Key)
			if exists && !opts.OverwriteExisting {
				continue
			}
			out.Set(e.Key, e.Value)
		}
	}
	return out
}
