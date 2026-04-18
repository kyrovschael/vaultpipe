package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func newHookLogger(buf *bytes.Buffer) *RedactHook {
	l := NewLogger(buf)
	return NewRedactHook(l, []string{"s3cr3t", "tok-abc"})
}

func TestRedactHook_RedactsEventString(t *testing.T) {
	var buf bytes.Buffer
	h := newHookLogger(&buf)
	h.Info("fetched secret s3cr3t from vault", nil)

	if strings.Contains(buf.String(), "s3cr3t") {
		t.Error("expected secret to be redacted in event string")
	}
	if !strings.Contains(buf.String(), "[REDACTED]") {
		t.Error("expected [REDACTED] placeholder in output")
	}
}

func TestRedactHook_RedactsMetaValues(t *testing.T) {
	var buf bytes.Buffer
	h := newHookLogger(&buf)
	h.Error("auth failed", map[string]string{"token": "tok-abc", "user": "alice"})

	if strings.Contains(buf.String(), "tok-abc") {
		t.Error("expected token to be redacted in meta")
	}
}

func TestRedactHook_PreservesNonSecretData(t *testing.T) {
	var buf bytes.Buffer
	h := newHookLogger(&buf)
	h.Info("process started", map[string]string{"pid": "1234"})

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry["event"] != "process started" {
		t.Errorf("unexpected event: %v", entry["event"])
	}
}

func TestRedactHook_EmptySecretsIgnored(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(&buf)
	h := NewRedactHook(l, []string{"", ""})
	h.Info("hello world", nil)

	if !strings.Contains(buf.String(), "hello world") {
		t.Error("expected message to be preserved when no secrets set")
	}
}
