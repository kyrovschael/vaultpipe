// Package env provides environment variable manipulation utilities.
//
// # PartitionSnapshot
//
// PartitionSnapshot splits a [Snapshot] into two snapshots based on a
// predicate function. Keys for which the predicate returns true are placed
// in the first (matched) snapshot; all others go into the second (rest)
// snapshot.
//
// Neither output snapshot shares memory with the source.
//
//	matched, rest := env.PartitionSnapshot(snap, env.DefaultPartitionOptions(), func(k, v string) bool {
//	    return strings.HasPrefix(k, "DB_")
//	})
package env
