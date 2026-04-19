package audit

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/env"
)

// RenameHook logs every key rename operation performed by env.RenameSnapshot.
type RenameHook struct {
	logger *Logger
}

// NewRenameHook creates a RenameHook that writes events to logger.
func NewRenameHook(l *Logger) *RenameHook {
	return &RenameHook{logger: l}
}

// Apply calls env.RenameSnapshot and emits an audit event describing what was
// renamed and whether unmapped keys were dropped.
func (h *RenameHook) Apply(s env.Snapshot, mapping map[string]string, opts env.RenameOptions) env.Snapshot {
	out := env.RenameSnapshot(s, mapping, opts)

	renamed := 0
	for old := range mapping {
		if _, ok := s[old]; ok {
			renamed++
		}
	}

	dropped := 0
	if opts.DropUnmapped {
		for k := range s {
			if _, mapped := mapping[k]; !mapped {
				dropped++
			}
		}
	}

	h.logger.Info(Event{
		Action: "rename_snapshot",
		Meta: map[string]string{
			"renamed": fmt.Sprintf("%d", renamed),
			"dropped": fmt.Sprintf("%d", dropped),
			"drop_unmapped": fmt.Sprintf("%v", opts.DropUnmapped),
		},
	})
	return out
}
