package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds the full vaultpipe configuration.
type Config struct {
	Vault   VaultConfig   `toml:"vault"`
	Secrets []SecretMount  `toml:"secret"`
}

// VaultConfig holds Vault connection and auth settings.
type VaultConfig struct {
	Address  string `toml:"address"`
	RoleID   string `toml:"role_id"`
	SecretID string `toml:"secret_id"`
	// SecretIDFile allows reading secret_id from a file at runtime.
	SecretIDFile string `toml:"secret_id_file"`
}

// SecretMount maps a Vault KV path to a set of env var overrides.
type SecretMount struct {
	// Path is the Vault KV path, e.g. "secret/data/myapp".
	Path string `toml:"path"`
	// EnvPrefix is optionally prepended to derived env var names.
	EnvPrefix string `toml:"env_prefix"`
}

// Load reads a TOML config file from the given path.
func Load(path string) (*Config, error) {
	if path == "" {
		path = "vaultpipe.toml"
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if _, err := toml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// Resolve secret_id from file if provided.
	if cfg.Vault.SecretIDFile != "" && cfg.Vault.SecretID == "" {
		data, err := os.ReadFile(cfg.Vault.SecretIDFile)
		if err != nil {
			return nil, fmt.Errorf("config: read secret_id_file: %w", err)
		}
		cfg.Vault.SecretID = string(data)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("config: vault.address is required")
	}
	if c.Vault.RoleID == "" {
		return fmt.Errorf("config: vault.role_id is required")
	}
	if len(c.Secrets) == 0 {
		return fmt.Errorf("config: at least one [[secret]] block is required")
	}
	for i, s := range c.Secrets {
		if s.Path == "" {
			return fmt.Errorf("config: secret[%d].path is required", i)
		}
	}
	return nil
}
