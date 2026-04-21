package env

// PartitionOptions controls the behaviour of [PartitionSnapshot].
type PartitionOptions struct {
	// CaseFold makes key comparisons inside the predicate case-insensitive
	// by normalising keys to upper-case before passing them to fn.
	CaseFold bool
}

// DefaultPartitionOptions returns a PartitionOptions with sensible defaults.
func DefaultPartitionOptions() PartitionOptions {
	return PartitionOptions{
		CaseFold: false,
	}
}

// PartitionSnapshot splits src into two snapshots.
// Keys for which fn(key, value) returns true go into the first (matched)
// snapshot; all remaining keys go into the second (rest) snapshot.
// fn receives the key exactly as stored unless opts.CaseFold is true, in
// which case the key is upper-cased before the call.
func PartitionSnapshot(src Snapshot, opts PartitionOptions, fn func(key, value string) bool) (matched Snapshot, rest Snapshot) {
	matched = make(Snapshot)
	rest = make(Snapshot)

	for k, v := range src {
		lookupKey := k
		if opts.CaseFold {
			lookupKey = toUpper(k)
		}
		if fn(lookupKey, v) {
			matched[k] = v
		} else {
			rest[k] = v
		}
	}
	return matched, rest
}

// toUpper is a minimal ASCII upper-case helper to avoid importing strings
// in the hot path.
func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
