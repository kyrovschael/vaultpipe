package env

// DefaultEnv returns a Snapshot populated from the current OS environment.
func DefaultEnv() Snapshot {
	return OSSource().Snapshot()
}
