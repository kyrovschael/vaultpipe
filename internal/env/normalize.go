package env

// NormalizeOptions controls the behaviour of NormalizeSnapshot.
type NormalizeOptions struct {
	// SanitizeOpts is forwarded to SanitizeSnapshot.
	SanitizeOpts SanitizeOptions
	// DedupeOpts is forwarded to DedupeSnapshot.
	DedupeOpts DedupeOptions
	// CompactOpts is forwarded to CompactSnapshot.
	CompactOpts CompactOptions
	// TransformOpts is forwarded to TransformSnapshot.
	TransformOpts TransformOptions
	// SkipUpperKeys disables the final upper-case key transformation.
	SkipUpperKeys bool
}

// DefaultNormalizeOptions returns a sensible default configuration.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		SanitizeOpts:  DefaultSanitizeOptions(),
		DedupeOpts:    DefaultDedupeOptions(),
		CompactOpts:   DefaultCompactOptions(),
		TransformOpts: DefaultTransformOptions(),
	}
}

// NormalizeSnapshot sanitizes, deduplicates, compacts, and upper-cases keys in
// src, returning a clean canonical Snapshot. The source is not mutated.
func NormalizeSnapshot(src Snapshot, opts NormalizeOptions) (Snapshot, error) {
	steps := []func(Snapshot) (Snapshot, error){
		func(s Snapshot) (Snapshot, error) {
			return SanitizeSnapshot(s, opts.SanitizeOpts)
		},
		func(s Snapshot) (Snapshot, error) {
			return DedupeSnapshot(s, opts.DedupeOpts)
		},
		func(s Snapshot) (Snapshot, error) {
			return CompactSnapshot(s, opts.CompactOpts)
		},
	}

	if !opts.SkipUpperKeys {
		steps = append(steps, func(s Snapshot) (Snapshot, error) {
			topts := opts.TransformOpts
			topts.KeyFns = []func(string) string{UpperKeys(s)[""[0:0:0]].Key}
			// Use the exported helper directly instead.
			return transformUpperKeys(s)
		})
	}

	p := NewPipeline(src)
	for _, step := range steps[:len(steps)-1] {
		p.Add(Lift(step))
	}
	// Final step: upper-case keys unless skipped.
	if !opts.SkipUpperKeys {
		p.Add(Lift(func(s Snapshot) (Snapshot, error) {
			return transformUpperKeys(s)
		}))
	}
	return p.Run()
}

// transformUpperKeys is an internal helper that upper-cases all keys.
func transformUpperKeys(s Snapshot) (Snapshot, error) {
	opts := DefaultTransformOptions()
	opts.KeyFns = []func(string) string{strings_toUpper}
	return TransformSnapshot(s, opts)
}

// strings_toUpper is an alias so we avoid importing "strings" at package level
// (it is already imported by other files in the package).
var strings_toUpper = func(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		out[i] = c
	}
	return string(out)
}
