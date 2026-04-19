package env

import "fmt"

// RequiredError is returned when one or more required keys are absent from a snapshot.
type RequiredError struct {
	Missing []string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("env: missing required keys: %v", e.Missing)
}

// RequireOptions controls the behaviour of RequireKeys.
type RequireOptions struct {
	// AllowEmpty permits keys that are present but have an empty value.
	AllowEmpty bool
}

// DefaultRequireOptions returns the default options for RequireKeys.
func DefaultRequireOptions() RequireOptions {
	return RequireOptions{AllowEmpty: false}
}

// RequireKeys checks that every key in keys exists in snap.
// If AllowEmpty is false, keys whose value is the empty string are also
// treated as missing.
// Returns a *RequiredError listing all absent keys, or nil on success.
func RequireKeys(snap Snapshot, keys []string, opts RequireOptions) error {
	var missing []string
	for _, k := range keys {
		v, ok := snap[k]
		if !ok {
			missing = append(missing, k)
			continue
		}
		if !opts.AllowEmpty && v == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return &RequiredError{Missing: missing}
	}
	return nil
}
