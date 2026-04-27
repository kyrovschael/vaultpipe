package audit

import (
	"fmt"
	"sort"
)

// NewResolveHook returns a Logger middleware that emits an audit event each
// time ResolveSnapshot is called, listing which keys were resolved from which
// source index and which keys were missing.
//
// sourceName is a human-readable label for the resolve operation (e.g. the
// Vault path or config section name).
func NewResolveHook(inner *Logger, sourceName string) *ResolveHook {
	return &ResolveHook{inner: inner, sourceName: sourceName}
}

// ResolveHook wraps a Logger and records resolve outcomes.
type ResolveHook struct {
	inner      *Logger
	sourceName string
}

// LogResolved emits an info event listing all keys that were successfully
// resolved and the zero-based source index they came from.
func (h *ResolveHook) LogResolved(resolved map[string]int) {
	if len(resolved) == 0 {
		return
	}
	keys := make([]string, 0, len(resolved))
	for k := range resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	meta := make(map[string]string, len(resolved)+1)
	meta["source"] = h.sourceName
	for _, k := range keys {
		meta[fmt.Sprintf("key.%s", k)] = fmt.Sprintf("source_index=%d", resolved[k])
	}
	h.inner.Info(Event{
		Action:  "resolve",
		Outcome: "ok",
		Message: fmt.Sprintf("resolved %d key(s) from %s", len(resolved), h.sourceName),
		Meta:    meta,
	})
}

// LogMissing emits a warning event listing keys that could not be resolved
// from any source.
func (h *ResolveHook) LogMissing(missing []string) {
	if len(missing) == 0 {
		return
	}
	sorted := make([]string, len(missing))
	copy(sorted, missing)
	sort.Strings(sorted)

	meta := map[string]string{"source": h.sourceName}
	for i, k := range sorted {
		meta[fmt.Sprintf("missing.%d", i)] = k
	}
	h.inner.Error(Event{
		Action:  "resolve",
		Outcome: "missing",
		Message: fmt.Sprintf("%d key(s) not found in any source for %s", len(missing), h.sourceName),
		Meta:    meta,
	})
}
