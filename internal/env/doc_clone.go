// Package env provides environment snapshot manipulation utilities.
//
// # Clone
//
// CloneSnapshot produces a deep copy of an environment snapshot,
// optionally restricting the result to a specific subset of keys.
//
// When Keys is non-empty only those keys are included; missing keys
// are silently skipped unless SkipMissing is false, in which case an
// error is returned for the first absent key.
//
// Example:
//
//	cloned, err := env.CloneSnapshot(snap, env.DefaultCloneOptions())
package env
