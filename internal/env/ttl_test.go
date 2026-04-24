package env

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestTTLSnapshot_StampsAllKeys(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	src := Snapshot{"A": "1", "B": "2"}
	entries, err := TTLSnapshot(src, TTLOptions{
		TTL: 10 * time.Second,
		Now: fixedClock(now),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	want := now.Add(10 * time.Second)
	for k, e := range entries {
		if !e.ExpiresAt.Equal(want) {
			t.Errorf("key %s: expiry = %v, want %v", k, e.ExpiresAt, want)
		}
	}
}

func TestTTLSnapshot_RestrictedKeys(t *testing.T) {
	now := time.Now()
	src := Snapshot{"A": "1", "B": "2", "C": "3"}
	entries, err := TTLSnapshot(src, TTLOptions{
		TTL:  time.Minute,
		Keys: []string{"A", "C"},
		Now:  fixedClock(now),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := entries["B"]; ok {
		t.Error("key B should not be present in restricted output")
	}
	if _, ok := entries["A"]; !ok {
		t.Error("key A should be present")
	}
}

func TestTTLSnapshot_InvalidTTL(t *testing.T) {
	_, err := TTLSnapshot(Snapshot{"X": "1"}, TTLOptions{TTL: 0})
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestTTLEntry_Expired(t *testing.T) {
	now := time.Now()
	e := TTLEntry{ExpiresAt: now.Add(-time.Second)}
	if !e.Expired(now) {
		t.Error("entry should be expired")
	}
	e2 := TTLEntry{ExpiresAt: now.Add(time.Hour)}
	if e2.Expired(now) {
		t.Error("entry should not be expired")
	}
}

func TestTTLSnapshot_EmptySnapshot(t *testing.T) {
	entries, err := TTLSnapshot(Snapshot{}, TTLOptions{TTL: time.Minute})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty map, got %d entries", len(entries))
	}
}
