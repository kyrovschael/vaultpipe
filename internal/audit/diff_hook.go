package audit

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/env"
)

// DiffHook is an audit hook that logs snapshot diffs as structured events.
type DiffHook struct {
	logger  *Logger
	redact  bool
}

// NewDiffHook returns a DiffHook that writes diff events to logger.
// When redact is true, env values are replaced with "***" in the log.
func NewDiffHook(logger *Logger, redact bool) *DiffHook {
	return &DiffHook{logger: logger, redact: redact}
}

// Record computes the diff between base and overlay and logs each change.
func (h *DiffHook) Record(base, overlay env.Snapshot) {
	opts := env.DiffOptions{RedactValues: h.redact}
	entries := env.DiffSnapshots(base, overlay, opts)
	for _, e := range entries {
		switch e.Kind {
		case env.DiffAdded:
			h.logger.Info(Event{
				Action: "env.diff",
				Meta: map[string]string{
					"kind":  "added",
					"key":   e.Key,
					"value": e.NewValue,
				},
			})
		case env.DiffRemoved:
			h.logger.Info(Event{
				Action: "env.diff",
				Meta: map[string]string{
					"kind":    "removed",
					"key":     e.Key,
					"old_val": e.OldValue,
				},
			})
		case env.DiffChanged:
			h.logger.Info(Event{
				Action: "env.diff",
				Meta: map[string]string{
					"kind":    "changed",
					"key":     e.Key,
					"old_val": e.OldValue,
					"new_val": e.NewValue,
				},
			})
		default:
			h.logger.Error(fmt.Errorf("unknown diff kind %q for key %s", e.Kind, e.Key))
		}
	}
}
