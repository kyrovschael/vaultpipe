package env

// InheritConfig controls which environment variables from the host
// are passed through to the child process.
type InheritConfig struct {
	// AllowList is an explicit set of keys to always pass through.
	// If empty, all keys not blocked by the DenyList are passed through.
	AllowList []string

	// DenyList filters out sensitive or unwanted keys.
	DenyList *DenyList
}

// Apply returns a new Snapshot containing only the variables that should
// be inherited according to the config.
func (c *InheritConfig) Apply(base Snapshot) Snapshot {
	out := make(Snapshot, len(base))

	if len(c.AllowList) > 0 {
		allowed := make(map[string]struct{}, len(c.AllowList))
		for _, k := range c.AllowList {
			allowed[k] = struct{}{}
		}
		for k, v := range base {
			if _, ok := allowed[k]; ok {
				out[k] = v
			}
		}
		return out
	}

	for k, v := range base {
		if c.DenyList != nil && c.DenyList.Blocked(k) {
			continue
		}
		out[k] = v
	}
	return out
}

// DefaultInheritConfig returns an InheritConfig using the default deny list.
func DefaultInheritConfig() *InheritConfig {
	return &InheritConfig{
		DenyList: DefaultDenyList(),
	}
}
