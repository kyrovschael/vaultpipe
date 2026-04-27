// Package env provides utilities for managing environment variable snapshots.
//
// # Resolve
//
// ResolveSnapshot resolves a set of keys from multiple sources in priority order.
// The first source that returns a non-empty value for a given key wins.
//
// Example:
//
//	resolved, err := env.ResolveSnapshot(keys, env.DefaultResolveOptions(),
//		vaultSource,
//		env.OSSource(),
//	)
//
// Keys not found in any source are either skipped (lenient) or cause an error
// (strict mode).
package env
