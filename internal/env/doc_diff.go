// Package env provides utilities for working with environment variable
// snapshots.
//
// # Diff
//
// DiffSnapshots compares two Snapshots and returns a slice of DiffEntry
// values describing keys that were added, removed, or changed when moving
// from base to overlay.
//
// Values can be redacted in the output by setting DiffOptions.RedactValues,
// which is useful when logging secret changes to an audit trail.
package env
