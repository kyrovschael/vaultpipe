// Package env provides utilities for managing environment variable snapshots.
//
// # Transform
//
// TransformSnapshot applies caller-supplied functions to the keys and/or values
// of a Snapshot, returning a new Snapshot without modifying the original.
//
// Convenience constructors are provided for the most common transformations:
//
//	// Upper-case all keys
//	out := env.TransformSnapshot(s, env.UpperKeys())
//
//	// Trim whitespace from all values
//	out = env.TransformSnapshot(s, env.TrimValues())
package env
