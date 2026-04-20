// Package env provides environment variable manipulation utilities.
//
// # Patch
//
// PatchSnapshot applies a set of changes to an existing snapshot.
// Changes can overwrite existing keys, add new ones, or remove keys
// by setting their value to the sentinel empty string (when SkipEmpty
// is false) or by listing them in the Remove field.
//
// Example:
//
//	base := env.StaticSource(env.Snapshot{"FOO": "bar", "BAZ": "qux"})
//	patch := env.PatchOptions{Remove: []string{"BAZ"}, Overlay: env.Snapshot{"FOO": "new"}}
//	result, err := env.PatchSnapshot(ctx, base, patch)
package env
