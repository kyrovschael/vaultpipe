// Package env — SquashSnapshot
//
// SquashSnapshot merges a slice of snapshots into a single snapshot,
// applying a caller-supplied merge function to resolve key conflicts.
//
// Ordering:
//
//	Snapshots are processed left-to-right. When two snapshots share a
//	key the MergeFn decides which value survives (or produces a new one).
//
// Default behaviour (MergeFn == nil):
//
//	The last value seen for each key wins, equivalent to a right-fold
//	over the input slice.
//
// Example:
//
//	out, err := env.SquashSnapshot(env.DefaultSquashOptions(), snaps...)
package env
