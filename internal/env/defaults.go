package env

import "os"

// DefaultEnv returns a Snapshot populated from the current process environment.
func DefaultEnv() Snapshot {
	return FromSlice(os.Environ())
}
