// Package env provides utilities for working with process environment snapshots.
//
// # Cast
//
// CastSnapshot applies a set of CastRules to a Snapshot, normalising string
// representations of typed values (int, float, bool, string). This is useful
// when secrets retrieved from Vault contain numeric or boolean values stored
// as raw strings that need to be normalised before injection.
//
// Example:
//
//	rules := []env.CastRule{
//		{Key: "PORT",  Type: env.CastInt},
//		{Key: "DEBUG", Type: env.CastBool},
//	}
//	normalised, err := env.CastSnapshot(snap, rules)
package env
