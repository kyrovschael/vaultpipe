// Package env – subset
//
// SubsetSnapshot extracts a named subset of keys from a Snapshot.
//
// Example:
//
//	out, err := env.SubsetSnapshot(full, env.SubsetOptions{
//		Keys:          []string{"HOME", "PATH", "USER"},
//		IgnoreMissing: true,
//	})
package env
