// Package env provides environment variable utilities.
package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CoerceType describes how a raw string value should be interpreted.
type CoerceType string

const (
	CoerceString CoerceType = "string"
	CoerceBool   CoerceType = "bool"
	CoerceInt    CoerceType = "int"
)

// CoerceRule maps an environment variable key to an expected type.
type CoerceRule struct {
	Key  string
	Type CoerceType
}

// CoerceSnapshot normalises values in snap according to rules.
// It returns a new Snapshot with coerced values and any validation errors.
func CoerceSnapshot(snap Snapshot, rules []CoerceRule) (Snapshot, []error) {
	out := snap.Clone()
	var errs []error

	for _, r := range rules {
		v, ok := out.Get(r.Key)
		if !ok {
			continue
		}
		coerced, err := coerceValue(v, r.Type)
		if err != nil {
			errs = append(errs, fmt.Errorf("coerce %s: %w", r.Key, err))
			continue
		}
		out.Set(r.Key, coerced)
	}
	return out, errs
}

func coerceValue(v string, t CoerceType) (string, error) {
	switch t {
	case CoerceString:
		return v, nil
	case CoerceBool:
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as bool", v)
		}
		if b {
			return "true", nil
		}
		return "false", nil
	case CoerceInt:
		_, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as int", v)
		}
		return strings.TrimSpace(v), nil
	default:
		return v, nil
	}
}
