package mask

import "io"

// Writer wraps an io.Writer and redacts secrets before writing.
type Writer struct {
	inner  io.Writer
	masker *Masker
}

// NewWriter returns a Writer that redacts secrets from all writes.
func NewWriter(w io.Writer, m *Masker) *Writer {
	return &Writer{inner: w, masker: m}
}

// Write redacts p before forwarding to the underlying writer.
func (w *Writer) Write(p []byte) (n int, err error) {
	redacted := w.masker.Redact(string(p))
	_, err = w.inner.Write([]byte(redacted))
	if err != nil {
		return 0, err
	}
	// Report original length so callers don't see a short-write error.
	return len(p), nil
}
