package env

import "fmt"

// DefaultPickOptions returns a PickOptions with safe defaults.
func DefaultPickOptions() PickOptions {
	return PickOptions{
		Strict:   false,
		CaseFold: false,
	}
}

// PickOptions controls the behaviour of PickSnapshot.
type PickOptions struct {
	// Strict causes PickSnapshot to return an error if any requested key is
	// absent from the source snapshot.
	Strict bool

	// CaseFold performs case-insensitive key matching when true.
	CaseFold bool
}

// PickSnapshot returns a new Snapshot containing only the entries whose keys
// appear in keys. The order of entries in the result follows the order of keys.
//
// When opts.Strict is true an error is returned as soon as a requested key
// cannot be found. When opts.CaseFold is true the comparison is
// case-insensitive and the original casing from the source snapshot is
// preserved in the output.
func PickSnapshot(src Snapshot, keys []string, opts PickOptions) (Snapshot, error) {
	out := make(Snapshot, 0, len(keys))

	for _, want := range keys {
		found := false
		for k, v := range src {
			match := k == want
			if opts.CaseFold {
				match = equalFold(k, want)
			}
			if match {
				out[k] = v
				found = true
				break
			}
		}
		if !found && opts.Strict {
			return nil, fmt.Errorf("env: pick: key %q not found in snapshot", want)
		}
	}

	return out, nil
}

// equalFold is a local ASCII-only case-fold helper to avoid importing strings
// in a hot path.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'a' && ca <= 'z' {
			ca -= 32
		}
		if cb >= 'a' && cb <= 'z' {
			cb -= 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}
