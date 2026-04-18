package env

import (
	"errors"
	"fmt"
	"regexp"
)

// validKeyRe matches POSIX-compliant environment variable names.
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds all invalid keys found during validation.
type ValidationError struct {
	InvalidKeys []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("env: invalid variable names: %v", e.InvalidKeys)
}

// ValidateSnapshot checks that all keys in the snapshot are valid POSIX
// environment variable names. Returns a *ValidationError listing bad keys,
// or nil if all keys are valid.
func ValidateSnapshot(s Snapshot) error {
	var bad []string
	for k := range s {
		if !validKeyRe.MatchString(k) {
			bad = append(bad, k)
		}
	}
	if len(bad) > 0 {
		return &ValidationError{InvalidKeys: bad}
	}
	return nil
}

// ValidateKey returns an error if the given key is not a valid POSIX
// environment variable name.
func ValidateKey(key string) error {
	if key == "" {
		return errors.New("env: key must not be empty")
	}
	if !validKeyRe.MatchString(key) {
		return fmt.Errorf("env: invalid variable name %q", key)
	}
	return nil
}
