package env

// DiffKind describes the type of change between two snapshots.
type DiffKind string

const (
	DiffAdded   DiffKind = "added"
	DiffRemoved DiffKind = "removed"
	DiffChanged DiffKind = "changed"
)

// DiffEntry represents a single key-level change between two snapshots.
type DiffEntry struct {
	Key      string
	Kind     DiffKind
	OldValue string
	NewValue string
}

// DiffOptions controls behaviour of DiffSnapshots.
type DiffOptions struct {
	// RedactValues replaces values with "***" in the returned entries.
	RedactValues bool
}

// DefaultDiffOptions returns sensible defaults.
func DefaultDiffOptions() DiffOptions {
	return DiffOptions{RedactValues: false}
}

// DiffSnapshots returns the set of changes that transform base into overlay.
func DiffSnapshots(base, overlay Snapshot, opts DiffOptions) []DiffEntry {
	var entries []DiffEntry

	for k, newVal := range overlay {
		oldVal, exists := base[k]
		if !exists {
			entries = append(entries, DiffEntry{
				Key:      k,
				Kind:     DiffAdded,
				NewValue: redactIf(opts.RedactValues, newVal),
			})
		} else if oldVal != newVal {
			entries = append(entries, DiffEntry{
				Key:      k,
				Kind:     DiffChanged,
				OldValue: redactIf(opts.RedactValues, oldVal),
				NewValue: redactIf(opts.RedactValues, newVal),
			})
		}
	}

	for k, oldVal := range base {
		if _, exists := overlay[k]; !exists {
			entries = append(entries, DiffEntry{
				Key:      k,
				Kind:     DiffRemoved,
				OldValue: redactIf(opts.RedactValues, oldVal),
			})
		}
	}

	return entries
}

func redactIf(redact bool, v string) string {
	if redact {
		return "***"
	}
	return v
}
