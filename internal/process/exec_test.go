package process

import (
	"testing"
)

func TestMergeEnv_NoConflict(t *testing.T) {
	base := []string{"HOME=/root", "PATH=/usr/bin"}
	secrets := map[string]string{"DB_PASSWORD": "s3cr3t"}

	result := mergeEnv(base, secrets)

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}

	if !contains(result, "DB_PASSWORD=s3cr3t") {
		t.Error("expected DB_PASSWORD=s3cr3t in result")
	}
}

func TestMergeEnv_SecretOverridesBase(t *testing.T) {
	base := []string{"HOME=/root", "DB_PASSWORD=plaintext"}
	secrets := map[string]string{"DB_PASSWORD": "vault_secret"}

	result := mergeEnv(base, secrets)

	if contains(result, "DB_PASSWORD=plaintext") {
		t.Error("base DB_PASSWORD should have been overridden")
	}
	if !contains(result, "DB_PASSWORD=vault_secret") {
		t.Error("expected vault DB_PASSWORD in result")
	}
}

func TestMergeEnv_EmptySecrets(t *testing.T) {
	base := []string{"HOME=/root", "PATH=/usr/bin"}
	secrets := map[string]string{}

	result := mergeEnv(base, secrets)

	if len(result) != len(base) {
		t.Fatalf("expected %d entries, got %d", len(base), len(result))
	}
}

func TestEnvKeyExtraction(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"KEY=VALUE", "KEY"},
		{"A=B=C", "A"},
		{"NOEQUALS", "NOEQUALS"},
	}

	for _, tc := range cases {
		got := envKey(tc.input)
		if got != tc.expected {
			t.Errorf("envKey(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
