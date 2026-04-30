// Package env provides utilities for manipulating environment variable snapshots.
//
// SquashSnapshot merges all values for keys sharing a common prefix into a
// single entry, joining them with a configurable separator. This is useful
// when multiple secret paths contribute fragments that belong together under
// one environment variable.
//
// Example:
//
//	snap := env.Snapshot{"DB_HOSTS_0": "a", "DB_HOSTS_1": "b", "DB_HOSTS_2": "c"}
//	out, _ := env.SquashSnapshot(snap, env.DefaultSquashOptions())
//	// out["DB_HOSTS"] == "a,b,c"
package env
