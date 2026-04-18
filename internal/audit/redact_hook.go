package audit

import (
	"strings"
)

// RedactHook wraps a Logger and redacts secret values from log messages
// before they are written.
type RedactHook struct {
	inner   *Logger
	secrets []string
}

// NewRedactHook creates a RedactHook that wraps the given Logger.
func NewRedactHook(l *Logger, secrets []string) *RedactHook {
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return &RedactHook{inner: l, secrets: filtered}
}

// Info logs an info-level event with secrets redacted from the message.
func (r *RedactHook) Info(event string, meta map[string]string) {
	r.inner.Info(r.redact(event), r.redactMeta(meta))
}

// Error logs an error-level event with secrets redacted from the message.
func (r *RedactHook) Error(event string, meta map[string]string) {
	r.inner.Error(r.redact(event), r.redactMeta(meta))
}

func (r *RedactHook) redact(s string) string {
	for _, secret := range r.secrets {
		s = strings.ReplaceAll(s, secret, "[REDACTED]")
	}
	return s
}

func (r *RedactHook) redactMeta(meta map[string]string) map[string]string {
	if meta == nil {
		return nil
	}
	out := make(map[string]string, len(meta))
	for k, v := range meta {
		out[k] = r.redact(v)
	}
	return out
}
