package env

import "strings"

// TransformOptions controls how values in a snapshot are transformed.
type TransformOptions struct {
	// KeyFn is applied to every key. Nil means no transformation.
	KeyFn func(string) string
	// ValueFn is applied to every value. Nil means no transformation.
	ValueFn func(string) string
}

// DefaultTransformOptions returns a no-op TransformOptions.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{}
}

// TransformSnapshot returns a new Snapshot with keys and/or values transformed
// by the functions in opts. The original snapshot is never mutated.
//
// If two keys collide after key transformation the last one (in iteration order)
// wins.
func TransformSnapshot(s Snapshot, opts TransformOptions) Snapshot {
	out := make(Snapshot, len(s))
	for k, v := range s {
		nk := k
		if opts.KeyFn != nil {
			nk = opts.KeyFn(k)
		}
		nv := v
		if opts.ValueFn != nil {
			nv = opts.ValueFn(v)
		}
		out[nk] = nv
	}
	return out
}

// UpperKeys returns a TransformOptions that upper-cases all keys.
func UpperKeys() TransformOptions {
	return TransformOptions{KeyFn: strings.ToUpper}
}

// LowerKeys returns a TransformOptions that lower-cases all keys.
func LowerKeys() TransformOptions {
	return TransformOptions{KeyFn: strings.ToLower}
}

// TrimValues returns a TransformOptions that trims whitespace from all values.
func TrimValues() TransformOptions {
	return TransformOptions{ValueFn: strings.TrimSpace}
}
