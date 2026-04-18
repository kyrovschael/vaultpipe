package template

import "fmt"

// PlaceholderError describes a failure to resolve a single vault placeholder.
type PlaceholderError struct {
	Path  string
	Field string
	Cause string
}

func (e *PlaceholderError) Error() string {
	return fmt.Sprintf("cannot resolve {{vault:%s#%s}}: %s", e.Path, e.Field, e.Cause)
}

// newPathError returns a PlaceholderError for a missing secret path.
func newPathError(path string) *PlaceholderError {
	return &PlaceholderError{Path: path, Cause: "path not found in fetched secrets"}
}

// newFieldError returns a PlaceholderError for a missing field within a path.
func newFieldError(path, field string) *PlaceholderError {
	return &PlaceholderError{Path: path, Field: field, Cause: "field not present"}
}
