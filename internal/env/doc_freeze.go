// Package env provides FreezeSnapshot for creating immutable, concurrency-safe
// views of a Snapshot.
//
// A FrozenSnapshot captures the state of a Snapshot at a point in time and
// exposes read-only access. It is safe for concurrent use by multiple goroutines
// and is useful when secrets must be passed across goroutine boundaries without
// risk of accidental mutation.
//
// Example:
//
//	snap := env.Snapshot{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}
//	frozen := env.FreezeSnapshot(snap)
//
//	val, ok := frozen.Get("DB_PASS") // "s3cr3t", true
//	clone := frozen.Snapshot()        // returns a mutable copy
package env
