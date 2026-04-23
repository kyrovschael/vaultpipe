// Package env — InvertSnapshot
//
// InvertSnapshot swaps keys and values in a Snapshot, producing a reverse
// look-up table. This is useful when secrets are stored under a canonical
// name but downstream processes need to reference them by their value.
//
// Example:
//
//	src := env.Snapshot{"DB_PASS": "s3cr3t", "API_KEY": "tok123"}
//	inv := env.InvertSnapshot(src, env.DefaultInvertOptions())
//	// inv == Snapshot{"s3cr3t": "DB_PASS", "tok123": "API_KEY"}
package env
