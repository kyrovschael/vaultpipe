package env

import (
	"testing"
)

func TestFilterByPrefix_StripsPrefix(t *testing.T) {
	snap := Snapshot{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "val"}
	got := FilterByPrefix(snap, "APP_", DefaultPrefixOptions())
	if _, ok := got["OTHER"]; ok {
		t.Error("unexpected key OTHER")
	}
	if got["HOST"] != "localhost" {
		t.Errorf("HOST = %q, want localhost", got["HOST"])
	}
	if got["PORT"] != "8080" {
		t.Errorf("PORT = %q, want 8080", got["PORT"])
	}
}

func TestFilterByPrefix_KeepsPrefix(t *testing.T) {
	snap := Snapshot{"APP_HOST": "localhost", "OTHER": "val"}
	opts := PrefixOptions{StripPrefix: false}
	got := FilterByPrefix(snap, "APP_", opts)
	if _, ok := got["APP_HOST"]; !ok {
		t.Error("expected APP_HOST to be present")
	}
	if _, ok := got["OTHER"]; ok {
		t.Error("unexpected key OTHER")
	}
}

func TestFilterByPrefix_EmptyPrefix_ReturnsClone(t *testing.T) {
	snap := Snapshot{"A": "1", "B": "2"}
	got := FilterByPrefix(snap, "", DefaultPrefixOptions())
	if len(got) != len(snap) {
		t.Errorf("len = %d, want %d", len(got), len(snap))
	}
}

func TestFilterByPrefix_KeyEqualsPrefix_Skipped(t *testing.T) {
	snap := Snapshot{"APP_": "ghost", "APP_KEY": "val"}
	got := FilterByPrefix(snap, "APP_", DefaultPrefixOptions())
	if _, ok := got[""]; ok {
		t.Error("empty key should be skipped")
	}
	if got["KEY"] != "val" {
		t.Errorf("KEY = %q, want val", got["KEY"])
	}
}

func TestAddPrefix_PrependsTAllKeys(t *testing.T) {
	snap := Snapshot{"HOST": "localhost", "PORT": "8080"}
	got := AddPrefix(snap, "APP_")
	if got["APP_HOST"] != "localhost" {
		t.Errorf("APP_HOST = %q, want localhost", got["APP_HOST"])
	}
	if got["APP_PORT"] != "8080" {
		t.Errorf("APP_PORT = %q, want 8080", got["APP_PORT"])
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
}

func TestAddPrefix_EmptyPrefix_ReturnsClone(t *testing.T) {
	snap := Snapshot{"X": "1"}
	got := AddPrefix(snap, "")
	if got["X"] != "1" {
		t.Errorf("X = %q, want 1", got["X"])
	}
}
