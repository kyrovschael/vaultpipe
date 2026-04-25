// Package env provides utilities for managing environment variable snapshots.
//
// # Index
//
// IndexSnapshot builds an inverted index from environment variable values to
// their keys, optionally filtering by a set of allowed keys and normalising
// keys before insertion.
//
// Basic usage:
//
//	snap := env.FromSlice([]string{"FOO=bar", "BAZ=qux"})
//	idx, err := env.IndexSnapshot(snap, env.DefaultIndexOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	keys, ok := idx.Lookup("bar") // ["FOO"], true
package env
