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
package env
