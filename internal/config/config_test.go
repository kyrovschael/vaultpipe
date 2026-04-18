package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultpipe.toml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	p := writeTemp(t, `
[vault]
address  = "https://vault.example.com"
role_id  = "my-role"
secret_id = "my-secret"

[[secret]]
path = "secret/data/myapp"
env_prefix = "APP"
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("address mismatch: %q", cfg.Vault.Address)
	}
	if len(cfg.Secrets) != 1 || cfg.Secrets[0].Path != "secret/data/myapp" {
		t.Errorf("secrets mismatch: %+v", cfg.Secrets)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	p := writeTemp(t, `
[vault]
role_id = "r"
secret_id = "s"

[[secret]]
path = "secret/data/x"
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestLoad_SecretIDFile(t *testing.T) {
	dir := t.TempDir()
	sidFile := filepath.Join(dir, "secret_id")
	if err := os.WriteFile(sidFile, []byte("file-secret"), 0600); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(dir, "vaultpipe.toml")
	content := "[vault]\naddress=\"http://v\"\nrole_id=\"r\"\nsecret_id_file=\"" + sidFile + "\"\n\n[[secret]]\npath=\"secret/data/y\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.SecretID != "file-secret" {
		t.Errorf("expected secret_id from file, got %q", cfg.Vault.SecretID)
	}
}

func TestLoad_NoSecrets(t *testing.T) {
	p := writeTemp(t, `
[vault]
address = "http://v"
role_id = "r"
secret_id = "s"
`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error when no secret blocks defined")
	}
}
