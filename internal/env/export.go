package env

import (
	"fmt"
	"sort"
	"strings"
)

// ExportFormat controls the output format of Export.
type ExportFormat int

const (
	// FormatShell produces POSIX shell export statements.
	FormatShell ExportFormat = iota
	// FormatDotenv produces KEY=VALUE lines compatible with .env files.
	FormatDotenv
)

// Export serialises a Snapshot into a human-readable string.
// Secret values are redacted before serialisation when redact is true.
func Export(snap Snapshot, format ExportFormat, redact bool) string {
	src := snap
	if redact {
		src = RedactSnapshot(snap)
	}

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := src[k]
		switch format {
		case FormatShell:
			fmt.Fprintf(&sb, "export %s=%s\n", k, shellQuote(v))
		case FormatDotenv:
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return sb.String()
}

// shellQuote wraps v in single quotes, escaping any existing single quotes.
func shellQuote(v string) string {
	return "'" + strings.ReplaceAll(v, "'", "'\\'''") + "'"
}
