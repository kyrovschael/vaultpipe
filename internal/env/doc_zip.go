// Package env provides ZipSnapshot for combining two environment snapshots
// key-by-key with a user-supplied merge function.
//
// # Overview
//
// ZipSnapshot walks both snapshots and calls opts.Fn for every key that
// appears in both. Keys exclusive to one side are retained or dropped
// according to KeepLeft / KeepRight. The originals are never mutated.
//
// # Example
//
//	result := env.ZipSnapshot(base, overlay, env.ZipOptions{
//		KeepLeft:  true,
//		KeepRight: true,
//		Fn: func(key, left, right string) string {
//			return left + ":" + right
//		},
//	})
package env
