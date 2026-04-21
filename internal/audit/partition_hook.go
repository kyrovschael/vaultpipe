package audit

import (
	"fmt"

	"github.com/yourusername/vaultpipe/internal/env"
)

// NewPartitionHook returns a Hook that logs a summary whenever a snapshot
// is partitioned. It records the sizes of the matched and rest partitions
// and, if redact is false, a sample of the matched keys (up to maxSample).
func NewPartitionHook(logger *Logger, label string, redact bool, maxSample int) func(matched, rest env.Snapshot) {
	return func(matched, rest env.Snapshot) {
		meta := map[string]string{
			"label":   label,
			"matched": fmt.Sprintf("%d", len(matched)),
			"rest":    fmt.Sprintf("%d", len(rest)),
			"total":   fmt.Sprintf("%d", len(matched)+len(rest)),
		}

		if !redact && maxSample > 0 {
			sample := make([]string, 0, maxSample)
			for k := range matched {
				if len(sample) >= maxSample {
					break
				}
				sample = append(sample, k)
			}
			meta["sample_keys"] = fmt.Sprintf("%v", sample)
		}

		logger.Info(Event{
			Action: "partition",
			Status: "ok",
			Meta:   meta,
		})
	}
}
