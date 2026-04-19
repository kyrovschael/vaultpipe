package env

import (
	"strings"
	"unicode"
)

// SanitizeOptions controls how SanitizeSnapshot behaves.
type SanitizeOptions struct {
	// ReplaceInvalidChars replaces invalid characters in keys with this string.
	// Defaults to "_".
	ReplaceInvalidChars string
	// SkipInvalidKeys drops keys that are still invalid after replacement.
	SkipInvalidKeys bool
}

// DefaultSanitizeOptions returns sensible defaults.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		ReplaceInvalidChars: "_",
		SkipInvalidKeys:     false,
	}
}

// SanitizeSnapshot returns a new Snapshot with all keys sanitized so they
// are valid POSIX environment variable names ([A-Za-z_][A-Za-z0-9_]*).
// Values are never modified.
func SanitizeSnapshot(src Snapshot, opts SanitizeOptions) Snapshot {
	out := make(Snapshot, len(src))
	for k, v := range src {
		sanitized := sanitizeKey(k, opts.ReplaceInvalidChars)
		if opts.SkipInvalidKeys && !isValidKey(sanitized) {
			continue
		}
		out[sanitized] = v
	}
	return out
}

func sanitizeKey(key, replacement string) string {
	if key == "" {
		return replacement
	}
	var b strings.Builder
	for i, r := range key {
		switch {
		case r == '_' || unicode.IsLetter(r):
			b.WriteRune(r)
		case unicode.IsDigit(r) && i > 0:
			b.WriteRune(r)
		default:
			b.WriteString(replacement)
		}
	}
	return b.String()
}

func isValidKey(key string) bool {
	if key == "" {
		return false
	}
	for i, r := range key {
		if r == '_' || unicode.IsLetter(r) {
			continue
		}
		if unicode.IsDigit(r) && i > 0 {
			continue
		}
		return false
	}
	return true
}
