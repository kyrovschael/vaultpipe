package env

import (
	"os"
	"strings"
)

// Snapshot captures the current process environment as a map.
type Snapshot struct {
	vars map[string]string
}

// Take captures os.Environ() into a Snapshot.
func Take() *Snapshot {
	return FromSlice(os.Environ())
}

// FromSlice builds a Snapshot from a KEY=VALUE slice.
func FromSlice(env []string) *Snapshot {
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		} else if len(parts) == 1 {
			m[parts[0]] = ""
		}
	}
	return &Snapshot{vars: m}
}

// Get returns the value for a key and whether it was present.
func (s *Snapshot) Get(key string) (string, bool) {
	v, ok := s.vars[key]
	return v, ok
}

// Keys returns all variable names in the snapshot.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.vars))
	for k := range s.vars {
		keys = append(keys, k)
	}
	return keys
}

// ToSlice converts the snapshot back to a KEY=VALUE slice.
func (s *Snapshot) ToSlice() []string {
	slice := make([]string, 0, len(s.vars))
	for k, v := range s.vars {
		slice = append(slice, k+"="+v)
	}
	return slice
}

// Merge returns a new Snapshot where values from overlay override s.
func (s *Snapshot) Merge(overlay map[string]string) *Snapshot {
	merged := make(map[string]string, len(s.vars)+len(overlay))
	for k, v := range s.vars {
		merged[k] = v
	}
	for k, v := range overlay {
		merged[k] = v
	}
	return &Snapshot{vars: merged}
}

// Diff returns the keys whose values differ between s and other,
// including keys present in one snapshot but not the other.
func (s *Snapshot) Diff(other *Snapshot) []string {
	var changed []string
	for k, v := range s.vars {
		if ov, ok := other.vars[k]; !ok || ov != v {
			changed = append(changed, k)
		}
	}
	for k := range other.vars {
		if _, ok := s.vars[k]; !ok {
			changed = append(changed, k)
		}
	}
	return changed
}
