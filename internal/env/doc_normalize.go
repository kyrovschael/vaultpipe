// Package env provides environment variable manipulation utilities.
//
// # NormalizeSnapshot
//
// NormalizeSnapshot applies a standard set of transformations to produce a
// clean, canonical snapshot suitable for injection into a child process.
//
// The pipeline applied (in order):
//  1. Sanitize keys — replace invalid characters, handle leading digits.
//  2. Deduplicate — last writer wins by default.
//  3. Compact — drop empty or whitespace-only values.
//  4. Transform — upper-case all keys.
//
// Example:
//
//	raw := env.StaticSource(env.Snapshot{
//		"my-key":  "hello",
//		"__bad":   "",
//		"MY_KEY":  "world",
//	})
//	out, err := env.NormalizeSnapshot(raw, env.DefaultNormalizeOptions())
package env
