// Package env provides utilities for filtering and sanitising environment
// variables before they are passed to a child process.
package env

import (
	"strings"
)

// DenyList holds variable name prefixes and exact names that must never be
// forwarded to the child process.
type DenyList struct {
	prefixes []string
	exact    map[string]struct{}
}

// NewDenyList constructs a DenyList from a slice of patterns.
// Patterns ending with "*" are treated as prefixes; all others are exact.
func NewDenyList(patterns []string) *DenyList {
	d := &DenyList{exact: make(map[string]struct{})}
	for _, p := range patterns {
		if strings.HasSuffix(p, "*") {
			d.prefixes = append(d.prefixes, strings.TrimSuffix(p, "*"))
		} else {
			d.exact[p] = struct{}{}
		}
	}
	return d
}

// Blocked reports whether the given variable name should be blocked.
func (d *DenyList) Blocked(name string) bool {
	if _, ok := d.exact[name]; ok {
		return true
	}
	for _, prefix := range d.prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// Filter returns a copy of env with any blocked variables removed.
// env entries are expected in KEY=VALUE form.
func (d *DenyList) Filter(env []string) []string {
	out := make([]string, 0, len(env))
	for _, entry := range env {
		name, _, _ := strings.Cut(entry, "=")
		if !d.Blocked(name) {
			out = append(out, entry)
		}
	}
	return out
}
