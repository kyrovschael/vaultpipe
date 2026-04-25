// Package env provides utilities for managing process environment variables.
//
// # Parse
//
// ParseSlice and ParseMap convert raw environment data into a typed Snapshot.
//
//	snap, err := env.ParseSlice(os.Environ(), env.DefaultParseOptions())
//
// ParseOptions.TrimSpace removes leading/trailing whitespace from keys and
// values before validation. ParseOptions.SkipInvalid silently drops malformed
// entries rather than returning an error, which is useful when consuming
// untrusted input.
//
// # Snapshot
//
// A Snapshot is an immutable view of environment variables at a point in time.
// Use Snapshot.Get to retrieve a value by key, and Snapshot.Keys to iterate
// over all present keys in sorted order.
//
// # Validation
//
// Keys must be non-empty and must not contain the '=' character. Values may
// be empty strings. Entries that violate these constraints are considered
// malformed and will cause ParseSlice or ParseMap to return an error unless
// ParseOptions.SkipInvalid is set.
package env
