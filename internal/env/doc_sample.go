// Package env provides environment variable snapshot utilities.
//
// # Sample
//
// SampleSnapshot returns a deterministic pseudo-random subset of entries
// from a snapshot. It is useful for logging, debugging, or producing
// representative previews without exposing the full environment.
//
// Example:
//
//	snap := env.StaticSource(map[string]string{
//		"A": "1", "B": "2", "C": "3", "D": "4",
//	})
//	result, _ := env.SampleSnapshot(snap, env.SampleOptions{N: 2, Seed: 42})
//	// result contains exactly 2 entries chosen deterministically.
package env
