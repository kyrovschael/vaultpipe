// Package env provides environment snapshot manipulation utilities.
//
// # Group
//
// GroupSnapshot partitions an environment snapshot into named buckets
// using a caller-supplied key function. Keys that produce an empty
// group name are collected under the fallback name "_".
//
// GroupNames returns the sorted list of group names present in the
// result of a GroupSnapshot call.
//
// Typical usage:
//
//	groups := env.GroupSnapshot(snap, env.DefaultGroupOptions())
//	for _, name := range env.GroupNames(groups) {
//		fmt.Println(name, groups[name])
//	}
package env
