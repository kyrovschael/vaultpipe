// Package env — ClampSnapshot
//
// ClampSnapshot constrains the values in a Snapshot to a bounded range of
// string length or numeric magnitude.
//
// # Length clamping
//
// When MinLen / MaxLen are non-zero, values shorter than MinLen are padded
// (right-padded with spaces by default) and values longer than MaxLen are
// truncated.  Set PadChar to override the padding rune.
//
// # Restricted keys
//
// When Keys is non-empty only those keys are clamped; all other entries are
// copied through unchanged.
package env
