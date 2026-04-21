package env

// DefaultZipOptions returns a ZipOptions with safe defaults.
func DefaultZipOptions() ZipOptions {
	return ZipOptions{
		KeepLeft:  true,
		KeepRight: true,
	}
}

// ZipOptions controls how two snapshots are combined key-by-key.
type ZipOptions struct {
	// KeepLeft retains keys that exist only in the left snapshot.
	KeepLeft bool
	// KeepRight retains keys that exist only in the right snapshot.
	KeepRight bool
	// Fn is called for each key present in both snapshots and returns the
	// merged value. When nil the right value wins.
	Fn func(key, left, right string) string
}

// ZipSnapshot merges two snapshots key-by-key using opts.Fn to resolve
// conflicts. Keys exclusive to either side are included according to
// KeepLeft / KeepRight. Neither source snapshot is mutated.
func ZipSnapshot(left, right Snapshot, opts ZipOptions) Snapshot {
	merge := opts.Fn
	if merge == nil {
		merge = func(_, _, r string) string { return r }
	}

	out := make(Snapshot)

	for k, lv := range left {
		if rv, ok := right[k]; ok {
			out[k] = merge(k, lv, rv)
		} else if opts.KeepLeft {
			out[k] = lv
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok && opts.KeepRight {
			out[k] = rv
		}
	}

	return out
}
