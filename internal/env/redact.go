package env

import "strings"

// SensitiveKeyPrefixes lists common prefixes for sensitive environment variables.
var SensitiveKeyPrefixes = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
}

const redactedValue = "[REDACTED]"

// RedactSnapshot returns a copy of the snapshot with sensitive values replaced.
func RedactSnapshot(s Snapshot) Snapshot {
	out := make(Snapshot, len(s))
	for k, v := range s {
		if isSensitiveKey(k) {
			out[k] = redactedValue
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactSlice returns a copy of the env slice with sensitive values replaced.
func RedactSlice(env []string) []string {
	out := make([]string, 0, len(env))
	for _, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			out = append(out, entry)
			continue
		}
		if isSensitiveKey(parts[0]) {
			out = append(out, parts[0]+"="+redactedValue)
		} else {
			out = append(out, entry)
		}
	}
	return out
}

func isSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, prefix := range SensitiveKeyPrefixes {
		if strings.HasPrefix(upper, prefix) {
			return true
		}
	}
	return false
}
