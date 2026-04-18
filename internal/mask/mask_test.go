package mask_test

import (
	"testing"

	"github.com/yourorg/vaultpipe/internal/mask"
)

func TestRedact_SingleSecret(t *testing.T) {
	m := mask.New([]string{"s3cr3t"})
	got := m.Redact("password is s3cr3t ok")
	want := "password is [REDACTED] ok"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_MultipleOccurrences(t *testing.T) {
	m := mask.New([]string{"abc"})
	got := m.Redact("abc and abc again")
	want := "[REDACTED] and [REDACTED] again"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedact_NoMatch(t *testing.T) {
	m := mask.New([]string{"xyz"})
	got := m.Redact("nothing here")
	if got != "nothing here" {
		t.Errorf("unexpected redaction: %q", got)
	}
}

func TestRedact_EmptySecret_Ignored(t *testing.T) {
	m := mask.New([]string{"", "real"})
	if m.Len() != 1 {
		t.Errorf("expected 1 secret, got %d", m.Len())
	}
}

func TestAdd_AppendsSecrets(t *testing.T) {
	m := mask.New(nil)
	m.Add("token123", "")
	if m.Len() != 1 {
		t.Errorf("expected 1 secret after Add, got %d", m.Len())
	}
	got := m.Redact("bearer token123")
	if got != "bearer [REDACTED]" {
		t.Errorf("got %q", got)
	}
}

func TestRedact_MultipleSecrets(t *testing.T) {
	m := mask.New([]string{"foo", "bar"})
	got := m.Redact("foo meets bar")
	want := "[REDACTED] meets [REDACTED]"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
