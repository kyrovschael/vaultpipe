package env

import (
	"os"
	"testing"
)

func TestDefaultEnv_ContainsKnownKey(t *testing.T) {
	// PATH (or GOPATH) should be present in virtually every test environment.
	snap := DefaultEnv()
	if snap == nil {
		t.Fatal("DefaultEnv returned nil snapshot")
	}

	// Verify at least one key from the real environment is present.
	envSlice := os.Environ()
	if len(envSlice) == 0 {
		t.Skip("no environment variables available")
	}

	if len(snap) == 0 {
		t.Fatal("DefaultEnv returned empty snapshot but os.Environ is non-empty")
	}
}

func TestDefaultEnv_ReturnsClone(t *testing.T) {
	snap := DefaultEnv()

	// Mutating the snapshot must not affect os.Environ.
	before := len(os.Environ())
	snap["__VAULTPIPE_TEST_KEY__"] = "sentinel"
	after := len(os.Environ())

	if before != after {
		t.Errorf("mutating snapshot affected os.Environ: before=%d after=%d", before, after)
	}

	if _, ok := os.LookupEnv("__VAULTPIPE_TEST_KEY__"); ok {
		t.Error("sentinel key leaked into real environment")
	}
}

func TestDefaultEnv_MatchesFromSlice(t *testing.T) {
	expected := FromSlice(os.Environ())
	actual := DefaultEnv()

	if len(actual) != len(expected) {
		t.Errorf("length mismatch: got %d, want %d", len(actual), len(expected))
	}

	for k, v := range expected {
		if got, ok := actual[k]; !ok {
			t.Errorf("key %q missing from DefaultEnv snapshot", k)
		} else if got != v {
			t.Errorf("key %q: got %q, want %q", k, got, v)
		}
	}
}
