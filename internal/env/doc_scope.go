// Package env — scope helpers.
//
// ScopeSnapshot narrows a snapshot to keys sharing a common prefix,
// optionally stripping that prefix from the result. This is useful when
// a single environment contains variables destined for multiple sub-systems
// (e.g. DB_*, CACHE_*) and you want to hand each sub-system only its own
// view.
//
// NamespaceSnapshot is the inverse: it takes a flat snapshot and re-exports
// every key under a chosen prefix, making it easy to merge scoped snapshots
// back into a unified environment.
package env
