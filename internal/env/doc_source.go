// Package env provides utilities for working with process environment snapshots.
//
// # Sources
//
// A SourceFunc is a context-aware function that produces an env.Snapshot.
// Several constructors are provided:
//
//   - OSSource        — reads os.Environ at call time
//   - StaticSource    — returns a fixed snapshot (cloned on each call)
//   - SliceSource     — parses a KEY=VALUE slice
//   - MapSource       — wraps a plain map[string]string
//   - ChainSource     — merges multiple sources; later sources win
//   - PrefixedOSSource — filters OS env by prefix, stripping it from keys
package env
