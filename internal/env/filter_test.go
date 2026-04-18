package env

import (
	"testing"
)

func TestBlocked_ExactMatch(t *testing.T) {
	d := NewDenyList([]string{"AWS_SECRET_ACCESS_KEY", "VAULT_TOKEN"})
	if !d.Blocked("AWS_SECRET_ACCESS_KEY") {
		t.Fatal("expected AWS_SECRET_ACCESS_KEY to be blocked")
	}
	if !d.Blocked("VAULT_TOKEN") {
		t.Fatal("expected VAULT_TOKEN to be blocked")
	}
	if d.Blocked("HOME") {
		t.Fatal("expected HOME to be allowed")
	}
}

func TestBlocked_PrefixMatch(t *testing.T) {
	d := NewDenyList([]string{"AWS_*", "VAULT_*"})
	for _, name := range []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "VAULT_ADDR", "VAULT_TOKEN"} {
		if !d.Blocked(name) {
			t.Fatalf("expected %s to be blocked", name)
		}
	}
	if d.Blocked("HOME") {
		t.Fatal("expected HOME to be allowed")
	}
}

func TestFilter_RemovesBlocked(t *testing.T) {
	d := NewDenyList([]string{"AWS_*", "VAULT_TOKEN"})
	input := []string{
		"HOME=/root",
		"AWS_ACCESS_KEY_ID=AKIA123",
		"VAULT_TOKEN=s.abc",
		"PATH=/usr/bin",
	}
	got := d.Filter(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(got), got)
	}
	for _, entry := range got {
		if entry == "AWS_ACCESS_KEY_ID=AKIA123" || entry == "VAULT_TOKEN=s.abc" {
			t.Fatalf("blocked entry leaked: %s", entry)
		}
	}
}

func TestFilter_EmptyDenyList(t *testing.T) {
	d := NewDenyList(nil)
	input := []string{"HOME=/root", "PATH=/usr/bin"}
	got := d.Filter(input)
	if len(got) != len(input) {
		t.Fatalf("expected %d entries, got %d", len(input), len(got))
	}
}

func TestFilter_AllBlocked(t *testing.T) {
	d := NewDenyList([]string{"HOME", "PATH"})
	got := d.Filter([]string{"HOME=/root", "PATH=/usr/bin"})
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}
