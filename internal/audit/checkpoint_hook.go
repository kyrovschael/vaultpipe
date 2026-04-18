package audit

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/env"
)

// CheckpointHook emits audit events whenever a named checkpoint detects
// environment drift relative to a saved baseline.
type CheckpointHook struct {
	logger *Logger
	opts   env.DiffOptions
}

// NewCheckpointHook creates a hook that logs diffs through logger.
func NewCheckpointHook(logger *Logger, opts env.DiffOptions) *CheckpointHook {
	return &CheckpointHook{logger: logger, opts: opts}
}

// Observe compares current against cp and emits one audit event per diff entry.
// It returns the number of drift entries detected.
func (h *CheckpointHook) Observe(cp *env.Checkpoint, current env.Snapshot) int {
	entries := cp.DiffFrom(current, h.opts)
	for _, e := range entries {
		meta := map[string]string{
			"checkpoint": cp.Name(),
			"key":        e.Key,
			"kind":       string(e.Kind),
		}
		if e.OldValue != "" {
			meta["old"] = e.OldValue
		}
		if e.NewValue != "" {
			meta["new"] = e.NewValue
		}
		h.logger.Info(fmt.Sprintf("env drift detected [%s] %s", e.Kind, e.Key), meta)
	}
	return len(entries)
}
