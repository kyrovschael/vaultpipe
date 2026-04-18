package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/audit"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	err := l.Info(audit.EventSecretsLoaded, map[string]string{"path": "secret/app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse log entry: %v", err)
	}

	if entry.Level != audit.LevelInfo {
		t.Errorf("expected level INFO, got %s", entry.Level)
	}
	if entry.Event != audit.EventSecretsLoaded {
		t.Errorf("unexpected event: %s", entry.Event)
	}
	if entry.Meta["path"] != "secret/app" {
		t.Errorf("meta path missing or wrong")
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	_ = l.Error(audit.EventAuthFailed, map[string]string{"reason": "invalid credentials"})

	if !strings.Contains(buf.String(), "auth.failed") {
		t.Error("expected event name in output")
	}
	if !strings.Contains(buf.String(), "ERROR") {
		t.Error("expected ERROR level in output")
	}
}

func TestLogger_NilWriter_UsesStderr(t *testing.T) {
	// Should not panic when w is nil.
	l := audit.NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogger_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	_ = l.Info(audit.EventAuthSuccess, nil)
	_ = l.Info(audit.EventProcessStart, map[string]string{"cmd": "env"})
	_ = l.Info(audit.EventProcessExit, map[string]string{"code": "0"})

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 log lines, got %d", len(lines))
	}
}
