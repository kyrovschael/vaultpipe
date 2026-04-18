package env

import "strings"

// ScopeOptions controls how a snapshot is scoped to a namespace.
type ScopeOptions struct {
	// Prefix is the namespace prefix, e.g. "APP_".
	Prefix string
	// StripPrefix removes the prefix from keys in the returned snapshot.
	StripPrefix bool
	// CaseFold normalises keys to upper-case before matching.
	CaseFold bool
}

// DefaultScopeOptions returns conservative defaults.
func DefaultScopeOptions(prefix string) ScopeOptions {
	return ScopeOptions{
		Prefix:      prefix,
		StripPrefix: true,
		CaseFold:    true,
	}
}

// ScopeSnapshot returns a new snapshot containing only keys that match the
// given prefix. Keys are optionally stripped of the prefix and/or uppercased.
func ScopeSnapshot(s Snapshot, opts ScopeOptions) Snapshot {
	out := make(Snapshot, len(s))
	pfx := opts.Prefix
	if opts.CaseFold {
		pfx = strings.ToUpper(pfx)
	}
	for k, v := range s {
		candidate := k
		if opts.CaseFold {
			candidate = strings.ToUpper(k)
		}
		if pfx == "" || strings.HasPrefix(candidate, pfx) {
			outKey := k
			if opts.StripPrefix && pfx != "" {
				outKey = k[len(pfx):]
				if outKey == "" {
					continue
				}
			}
			out[outKey] = v
		}
	}
	return out
}

// NamespaceSnapshot injects all keys from src into dst under the given prefix,
// returning a new snapshot without mutating either input.
func NamespaceSnapshot(src Snapshot, prefix string) Snapshot {
	out := make(Snapshot, len(src))
	for k, v := range src {
		out[prefix+k] = v
	}
	return out
}
