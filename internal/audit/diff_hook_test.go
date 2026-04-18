package audit

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

func newDiffLogger(buf *bytes.Buffer) *Logger {
	return NewLogger(buf)
}

func collectEvents(buf *bytes.Buffer) []map[string]interface{} {
	var events []map[string]interface{}
	dec := json.NewDecoder(buf)
	for dec.More() {
		var m map[string]interface{}
		if err := dec.Decode(&m); err == nil {
			events = append(events, m)
		}
	}
	return events
}

func TestDiffHook_LogsAdded(t *testing.T) {
	var buf bytes.Buffer
	h := NewDiffHook(newDiffLogger(&buf), false)
	h.Record(env.Snapshot{}, env.Snapshot{"NEW_KEY": "val"})
	events := collectEvents(&buf)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestDiffHook_LogsRemoved(t *testing.T) {
	var buf bytes.Buffer
	h := NewDiffHook(newDiffLogger(&buf), false)
	h.Record(env.Snapshot{"OLD": "v"}, env.Snapshot{})
	events := collectEvents(&buf)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestDiffHook_LogsChanged(t *testing.T) {
	var buf bytes.Buffer
	h := NewDiffHook(newDiffLogger(&buf), false)
	h.Record(env.Snapshot{"K": "old"}, env.Snapshot{"K": "new"})
	events := collectEvents(&buf)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestDiffHook_RedactsValues(t *testing.T) {
	var buf bytes.Buffer
	h := NewDiffHook(newDiffLogger(&buf), true)
	h.Record(env.Snapshot{"SECRET": "old"}, env.Snapshot{"SECRET": "new"})
	if bytes.Contains(buf.Bytes(), []byte("old")) || bytes.Contains(buf.Bytes(), []byte("new")) {
		t.Fatal("expected values to be redacted")
	}
}

func TestDiffHook_NoEventsWhenEqual(t *testing.T) {
	var buf bytes.Buffer
	h := NewDiffHook(newDiffLogger(&buf), false)
	h.Record(env.Snapshot{"A": "1"}, env.Snapshot{"A": "1"})
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %s", buf.String())
	}
}
