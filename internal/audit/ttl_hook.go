package audit

import (
	"fmt"
	"time"
)

// TTLEvent is emitted when a secret entry is found to be expired.
type TTLEvent struct {
	Key       string
	ExpiresAt time.Time
	DetectedAt time.Time
}

func (e TTLEvent) String() string {
	return fmt.Sprintf("secret expired: key=%s expired_at=%s detected_at=%s",
		e.Key, e.ExpiresAt.Format(time.RFC3339), e.DetectedAt.Format(time.RFC3339))
}

// TTLHook wraps a Logger and emits an audit event for every expired
// TTL entry supplied to Check.
type TTLHook struct {
	logger *Logger
	now    func() time.Time
}

// NewTTLHook creates a TTLHook that writes to l. When now is nil
// time.Now is used.
func NewTTLHook(l *Logger, now func() time.Time) *TTLHook {
	if now == nil {
		now = time.Now
	}
	return &TTLHook{logger: l, now: now}
}

// Check inspects entries (key → expiry time) and logs an audit event
// for each key that has expired. It returns the set of expired keys.
func (h *TTLHook) Check(entries map[string]time.Time) []string {
	detected := h.now()
	var expired []string
	for k, exp := range entries {
		if detected.After(exp) {
			expired = append(expired, k)
			h.logger.Info(TTLEvent{
				Key:        k,
				ExpiresAt:  exp,
				DetectedAt: detected,
			}.String(), map[string]string{
				"key":        k,
				"expires_at": exp.Format(time.RFC3339),
			})
		}
	}
	return expired
}
