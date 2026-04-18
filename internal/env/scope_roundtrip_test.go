package env

import "testing"

// TestScopeRoundTrip verifies that Namespace ∘ Scope is the identity when
// StripPrefix is true and CaseFold is false.
func TestScopeRoundTrip(t *testing.T) {
	original := Snapshot{
		"SVC_ADDR": "0.0.0.0",
		"SVC_PORT": "9000",
		"UNRELATED": "drop",
	}

	const pfx = "SVC_"
	scoped := ScopeSnapshot(original, ScopeOptions{
		Prefix:      pfx,
		StripPrefix: true,
		CaseFold:    false,
	})
	renamed := NamespaceSnapshot(scoped, pfx)

	for _, k := range []string{"SVC_ADDR", "SVC_PORT"} {
		if renamed[k] != original[k] {
			t.Errorf("key %s: want %q got %q", k, original[k], renamed[k])
		}
	}
	if _, ok := renamed["UNRELATED"]; ok {
		t.Error("UNRELATED should not survive the round-trip")
	}
}

func TestNamespaceScope_EmptySource(t *testing.T) {
	out := NamespaceSnapshot(Snapshot{}, "X_")
	if len(out) != 0 {
		t.Errorf("expected empty snapshot, got %d keys", len(out))
	}
}
