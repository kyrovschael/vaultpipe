package audit_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/internal/audit"
)

func TestFileLogger_WritesAndCloses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	fl, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("NewFileLogger: %v", err)
	}

	if err := fl.Info(audit.EventProcessStart, map[string]string{"cmd": "ls"}); err != nil {
		t.Fatalf("Info: %v", err)
	}

	if err := fl.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if entry.Event != audit.EventProcessStart {
		t.Errorf("expected event %s, got %s", audit.EventProcessStart, entry.Event)
	}
	if entry.Meta["cmd"] != "ls" {
		t.Errorf("expected cmd=ls in meta")
	}
}

func TestFileLogger_InvalidPath(t *testing.T) {
	_, err := audit.NewFileLogger("/nonexistent-dir/audit.log")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
