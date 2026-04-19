package env

// DefaultEnv returns a Snapshot populated from the current OS environment,
// filtered through the default deny list and optionally extended with
// caller-supplied overrides.
func DefaultEnv(overrides ...Snapshot) (Snapshot, error) {
	base, err := OSSource()()
	if err != nil {
		return Snapshot{}, err
	}

	dl := DefaultDenyList()
	filtered := dl.Filter(base)

	if len(overrides) == 0 {
		return filtered, nil
	}

	opts := DefaultMergeOptions()
	result := filtered
	for _, ov := range overrides {
		result = MergeSnapshots(result, ov, opts)
	}
	return result, nil
}
