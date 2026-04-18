package mask_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/vaultpipe/internal/mask"
)

func TestWriter_RedactsOutput(t *testing.T) {
	var buf bytes.Buffer
	m := mask.New([]string{"supersecret"})
	w := mask.NewWriter(&buf, m)

	input := "token=supersecret\n"
	n, err := w.Write([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Errorf("expected n=%d, got %d", len(input), n)
	}
	got := buf.String()
	if got != "token=[REDACTED]\n" {
		t.Errorf("got %q", got)
	}
}

func TestWriter_PassthroughWhenNoMatch(t *testing.T) {
	var buf bytes.Buffer
	m := mask.New([]string{"hidden"})
	w := mask.NewWriter(&buf, m)

	_, _ = w.Write([]byte("hello world"))
	if buf.String() != "hello world" {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWriter_MultipleWrites(t *testing.T) {
	var buf bytes.Buffer
	m := mask.New([]string{"pwd"})
	w := mask.NewWriter(&buf, m)

	_, _ = w.Write([]byte("a=pwd "))
	_, _ = w.Write([]byte("b=pwd"))

	want := "a=[REDACTED] b=[REDACTED]"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}
