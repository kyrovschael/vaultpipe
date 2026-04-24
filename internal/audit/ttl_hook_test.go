package audit

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func newTTLLogger(buf *bytes.Buffer) *Logger {
	return NewLogger(buf)
}

func TestTTLHook_ReportsExpiredKeys(t *testing.T) {
	var buf bytes.Buffer
	hook := NewTTLHook(newTTLLogger(&buf), func() time.Time {
		return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	})
	entries := map[string]time.Time{
		"DB_PASS": time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC), // expired
		"API_KEY": time.Date(2024, 6, 1, 13, 0, 0, 0, time.UTC), // valid
	}
	expired := hook.Check(entries)
	if len(expired) != 1 || expired[0] != "DB_PASS" {
		t.Errorf("expected [DB_PASS], got %v", expired)
	}
	if !strings.Contains(buf.String(), "DB_PASS") {
		t.Error("audit log should mention DB_PASS")
	}
	if strings.Contains(buf.String(), "API_KEY") {
		t.Error("audit log should not mention API_KEY")
	}
}

func TestTTLHook_NoExpiredKeys(t *testing.T) {
	var buf bytes.Buffer
	hook := NewTTLHook(newTTLLogger(&buf), func() time.Time {
		return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	})
	entries := map[string]time.Time{
		"SECRET": time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	expired := hook.Check(entries)
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %v", expired)
	}
	if buf.Len() != 0 {
		t.Error("expected empty audit log")
	}
}

func TestTTLHook_NilNowDefaultsToTimeNow(t *testing.T) {
	var buf bytes.Buffer
	hook := NewTTLHook(newTTLLogger(&buf), nil)
	// All entries far in the future – nothing should expire.
	entries := map[string]time.Time{
		"K": time.Now().Add(24 * time.Hour),
	}
	expired := hook.Check(entries)
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %v", expired)
	}
}

func TestTTLHook_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	hook := NewTTLHook(newTTLLogger(&buf), nil)
	expired := hook.Check(nil)
	if len(expired) != 0 {
		t.Errorf("expected empty result, got %v", expired)
	}
}
