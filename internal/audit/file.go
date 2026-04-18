package audit

import (
	"fmt"
	"os"
)

// FileLogger wraps Logger with a backing file that can be closed.
type FileLogger struct {
	*Logger
	f *os.File
}

// NewFileLogger opens (or creates) the file at path and returns a FileLogger.
func NewFileLogger(path string) (*FileLogger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &FileLogger{
		Logger: NewLogger(f),
		f:      f,
	}, nil
}

// Close flushes and closes the underlying file.
func (fl *FileLogger) Close() error {
	if fl.f != nil {
		return fl.f.Close()
	}
	return nil
}
