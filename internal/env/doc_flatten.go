// Package env provides utilities for working with process environment
// snapshots.
//
// # Flatten
//
// FlattenMap converts a nested map (as returned by Vault's KV v2 data field)
// into a flat Snapshot by joining key path segments with a separator.
//
// Example:
//
//	raw := map[string]any{
//	    "database": map[string]any{"password": "s3cr3t"},
//	}
//	snap := env.FlattenMap(raw, env.DefaultFlattenOptions())
//	// snap["DATABASE_PASSWORD"] == "s3cr3t"
package env
