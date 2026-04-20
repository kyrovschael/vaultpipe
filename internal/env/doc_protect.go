// Package env provides environment variable manipulation utilities.
//
// # Protect
//
// ProtectSnapshot marks a set of keys as immutable. Once protected,
// those keys cannot be overwritten by subsequent merge or overlay
// operations via ApplyProtected.
//
// This is useful when certain secrets or system variables must not be
// shadowed by user-supplied environment files or dynamic configuration.
//
// Example:
//
//	protected, err := env.ProtectSnapshot(ctx, source, env.DefaultProtectOptions())
//	result, err := env.ApplyProtected(protected, overlay)
package env
