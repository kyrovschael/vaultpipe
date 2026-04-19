// Package env provides utilities for working with process environment snapshots.
//
// # Interpolation
//
// InterpolateSnapshot resolves placeholder tokens embedded inside snapshot
// values using a separate secrets snapshot as the source of truth.
//
// The default placeholder syntax is ${KEY}, where KEY must exactly match a key
// present in the secrets snapshot. When Strict mode is enabled any unresolved
// placeholder causes an error; otherwise the original token is preserved.
//
// Example:
//
//	snap    := env.Snapshot{"DSN": "postgres://user:${DB_PASS}@host/db"}
//	secrets := env.Snapshot{"DB_PASS": "s3cr3t"}
//	out, _  := env.InterpolateSnapshot(snap, secrets, env.DefaultInterpolateOptions())
//	// out["DSN"] == "postgres://user:s3cr3t@host/db"
package env
