// Package mask provides utilities for redacting secret values
// from log output and error messages.
package mask

import "strings"

const redacted = "[REDACTED]"

// Masker holds a set of sensitive values and can redact them from strings.
type Masker struct {
	secrets []string
}

// New creates a Masker preloaded with the provided secret values.
// Empty strings are ignored.
func New(values []string) *Masker {
	filtered := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" {
			filtered = append(filtered, v)
		}
	}
	return &Masker{secrets: filtered}
}

// Redact replaces all known secret values in s with [REDACTED].
func (m *Masker) Redact(s string) string {
	for _, secret := range m.secrets {
		s = strings.ReplaceAll(s, secret, redacted)
	}
	return s
}

// Add appends additional secret values to the masker at runtime.
func (m *Masker) Add(values ...string) {
	for _, v := range values {
		if v != "" {
			m.secrets = append(m.secrets, v)
		}
	}
}

// Len returns the number of tracked secret values.
func (m *Masker) Len() int {
	return len(m.secrets)
}
