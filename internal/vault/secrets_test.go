package vault

import (
	"sort"
	"testing"
)

func TestEnvKey(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"db-password", "DB_PASSWORD"},
		{"API_KEY", "API_KEY"},
		{"my-secret-value", "MY_SECRET_VALUE"},
		{"simple", "SIMPLE"},
	}
	for _, tc := range cases {
		got := EnvKey(tc.input)
		if got != tc.want {
			t.Errorf("EnvKey(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestSecretsToEnv(t *testing.T) {
	secrets := map[string]string{
		"db-host":     "localhost",
		"db-password": "s3cr3t",
	}

	env := SecretsToEnv(secrets)
	sort.Strings(env)

	want := []string{
		"DB_HOST=localhost",
		"DB_PASSWORD=s3cr3t",
	}

	if len(env) != len(want) {
		t.Fatalf("SecretsToEnv: got %d entries, want %d", len(env), len(want))
	}
	for i := range want {
		if env[i] != want[i] {
			t.Errorf("entry %d: got %q, want %q", i, env[i], want[i])
		}
	}
}
