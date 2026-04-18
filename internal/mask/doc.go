// Package mask provides secret-redaction utilities for vaultpipe.
//
// It exposes a Masker that tracks sensitive string values and can
// replace them with the literal token [REDACTED] wherever they appear.
// A masked io.Writer wrapper (Writer) allows stdout/stderr streams
// produced by child processes to be sanitised before they reach any
// log sink or terminal, preventing accidental secret leakage.
//
// Typical usage:
//
//	m := mask.New(secretValues)
//	maskedStdout := mask.NewWriter(os.Stdout, m)
//	maskedStderr := mask.NewWriter(os.Stderr, m)
package mask
