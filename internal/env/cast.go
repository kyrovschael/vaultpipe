package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents a target type for casting an env value.
type CastType string

const (
	CastString CastType = "string"
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastBool   CastType = "bool"
)

// CastRule defines a key and the type to cast its value to.
type CastRule struct {
	Key  string
	Type CastType
}

// CastSnapshot applies cast rules to a snapshot, returning a new snapshot
// with normalised string representations of the typed values.
// Invalid conversions are collected and returned as a combined error.
func CastSnapshot(snap Snapshot, rules []CastRule) (Snapshot, error) {
	out := snap.Clone()
	var errs []string

	for _, r := range rules {
		v, ok := out[r.Key]
		if !ok {
			continue
		}
		normalised, err := castValue(v, r.Type)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", r.Key, err))
			continue
		}
		out[r.Key] = normalised
	}

	if len(errs) > 0 {
		return out, fmt.Errorf("cast errors: %s", strings.Join(errs, "; "))
	}
	return out, nil
}

func castValue(v string, t CastType) (string, error) {
	switch t {
	case CastString:
		return v, nil
	case CastInt:
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return v, fmt.Errorf("cannot cast %q to int", v)
		}
		return strconv.FormatInt(n, 10), nil
	case CastFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return v, fmt.Errorf("cannot cast %q to float", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case CastBool:
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return v, fmt.Errorf("cannot cast %q to bool", v)
		}
		return strconv.FormatBool(b), nil
	default:
		return v, fmt.Errorf("unknown cast type %q", t)
	}
}
