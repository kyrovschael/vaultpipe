// Package env – rename
//
// RenameSnapshot applies a key-rename mapping to a Snapshot, returning a new
// Snapshot with the requested keys renamed. Unmapped keys are either passed
// through unchanged or dropped, depending on RenameOptions.DropUnmapped.
//
// Example:
//
//	out := env.RenameSnapshot(snap,
//		map[string]string{"DB_PASS": "DATABASE_PASSWORD"},
//		env.DefaultRenameOptions(),
//	)
package env
