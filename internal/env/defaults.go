package env

// DefaultEnv returns a Snapshot populated from the current OS environment,
// filtered through the default deny list and inheritance rules.
func DefaultEnv() (Snapshot, error) {
	src := OSSource()
	base, err := src.Load()
	if err != nil {
		return Snapshot{}, err
	}
	return base, nil
}
