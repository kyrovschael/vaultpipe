package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Level represents the audit log level.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Event     string            `json:"event"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Logger writes structured audit entries.
type Logger struct {
	out io.Writer
}

// NewLogger creates a Logger writing to w. Pass nil to use stderr.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{out: w}
}

// Log writes an audit entry at the given level.
func (l *Logger) Log(level Level, event string, meta map[string]string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Event:     event,
		Meta:      meta,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = l.out.Write(data)
	return err
}

// Info is a convenience wrapper for LevelInfo.
func (l *Logger) Info(event string, meta map[string]string) error {
	return l.Log(LevelInfo, event, meta)
}

// Warn is a convenience wrapper for LevelWarn.
func (l *Logger) Warn(event string, meta map[string]string) error {
	return l.Log(LevelWarn, event, meta)
}

// Error is a convenience wrapper for LevelError.
func (l *Logger) Error(event string, meta map[string]string) error {
	return l.Log(LevelError, event, meta)
}
