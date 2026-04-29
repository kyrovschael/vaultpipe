package audit

import (
	"fmt"

	"github.com/your-org/vaultpipe/internal/env"
)

// NewWindowHook returns a Hook that logs the keys retained and discarded when
// WindowSnapshot is applied. It is intended to be used alongside
// env.WindowSnapshot to provide an audit trail of range-filtering decisions.
func NewWindowHook(logger *Logger, opts env.WindowOptions) func(before, after env.Snapshot) {
	return func(before, after env.Snapshot) {
		afterSet := make(map[string]struct{}, len(after))
		for _, e := range after {
			afterSet[e.Key] = struct{}{}
		}

		var kept, dropped []string
		for _, e := range before {
			if _, ok := afterSet[e.Key]; ok {
				kept = append(kept, e.Key)
			} else {
				dropped = append(dropped, e.Key)
			}
		}

		logger.Info("window_snapshot", map[string]string{
			"from":    opts.From,
			"to":      opts.To,
			"kept":    fmt.Sprintf("%d", len(kept)),
			"dropped": fmt.Sprintf("%d", len(dropped)),
		})
	}
}
