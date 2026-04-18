// Package env provides utilities for managing process environment variables.
//
// The watch sub-feature (watch.go) allows callers to poll an arbitrary
// Snapshot source and receive callbacks whenever the environment changes.
// This is used by vaultpipe to detect secret rotation and trigger a
// controlled process restart or in-place environment update.
package env
