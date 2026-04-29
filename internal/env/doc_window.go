// Package env – window
//
// WindowSnapshot returns a new snapshot containing only the entries whose
// keys sort between [From, To] (inclusive). Both bounds are optional; an
// empty string disables that side of the range.
//
// Options
//
//	From      – lower bound key (inclusive, empty = no lower bound)
//	To        – upper bound key (inclusive, empty = no upper bound)
//	CaseFold  – compare keys case-insensitively
package env
