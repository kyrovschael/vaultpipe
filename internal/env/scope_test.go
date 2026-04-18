package env

import (
	"testing"
)

func TestScopeSnapshot_StripsPrefix(t *testing.T) {
	s := Snapshot{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "x"}
	got := ScopeSnapshot(s, DefaultScopeOptions("APP_"))
	if _, ok := got["HOST"]; !ok {
		t.Error("expected HOST key")
	}
	if _, ok := got["PORT"]; !ok {
		t.Error("expected PORT key")
	}
	if _, ok := got["OTHER"]; ok {
		t.Error("OTHER should be excluded")
	}
}

func TestScopeSnapshot_KeepsPrefix(t *testing.T) {
	s := Snapshot{"APP_HOST": "localhost", "OTHER": "x"}
	opts := ScopeOptions{Prefix: "APP_", StripPrefix: false, CaseFold: true}
	got := ScopeSnapshot(s, opts)
	if _, ok := got["APP_HOST"]; !ok {
		t.Error("expected APP_HOST key when StripPrefix=false")
	}
}

func TestScopeSnapshot_EmptyPrefix_ReturnsAll(t *testing.T) {
	s := Snapshot{"A": "1", "B": "2"}
	got := ScopeSnapshot(s, ScopeOptions{Prefix: "", StripPrefix: false})
	if len(got) != len(s) {
		t.Errorf("expected %d keys, got %d", len(s), len(got))
	}
}

func TestScopeSnapshot_KeyEqualsPrefix_Skipped(t *testing.T) {
	s := Snapshot{"APP_": "orphan", "APP_X": "val"}
	got := ScopeSnapshot(s, DefaultScopeOptions("APP_"))
	if _, ok := got[""]; ok {
		t.Error("empty key after strip should be skipped")
	}
	if got["X"] != "val" {
		t.Errorf("expected val, got %q", got["X"])
	}
}

func TestScopeSnapshot_CaseFold(t *testing.T) {
	s := Snapshot{"app_host": "h"}
	got := ScopeSnapshot(s, DefaultScopeOptions("APP_"))
	if _, ok := got["host"]; !ok {
		t.Error("case-fold match should include app_host")
	}
}

func TestNamespaceSnapshot_PrependsPrefix(t *testing.T) {
	s := Snapshot{"HOST": "localhost", "PORT": "5432"}
	got := NamespaceSnapshot(s, "DB_")
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestNamespaceSnapshot_DoesNotMutateSource(t *testing.T) {
	s := Snapshot{"KEY": "val"}
	_ = NamespaceSnapshot(s, "NS_")
	if _, ok := s["NS_KEY"]; ok {
		t.Error("source snapshot should not be mutated")
	}
}
