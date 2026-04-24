// Package env provides environment variable manipulation utilities.
//
// # LimitSnapshot
//
// LimitSnapshot enforces upper bounds on a Snapshot:
//
//   - MaxKeys caps the total number of entries retained.
//   - MaxValueLen caps the byte length of any individual value.
//
// By default both limits are disabled (zero value). When StrictKeys or
// StrictValues is true the function returns an error instead of silently
// dropping offending entries.
//
// Example:
//
//	out, err := env.LimitSnapshot(snap, env.LimitOptions{
//		MaxKeys:     100,
//		MaxValueLen: 4096,
//		StrictValues: true,
//	})
package env
