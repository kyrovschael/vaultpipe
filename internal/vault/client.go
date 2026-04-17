package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with vaultpipe-specific helpers.
type Client struct {
	api *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
	RoleID  string
	SecretID string
}

// NewClient creates a new Vault client from the provided config.
// Falls back to environment variables (VAULT_ADDR, VAULT_TOKEN) when fields are empty.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr != "" {
		vcfg.Address = addr
	}

	raw, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("vault: create client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	if token != "" {
		raw.SetToken(token)
	} else if cfg.RoleID != "" && cfg.SecretID != "" {
		if err := loginAppRole(raw, cfg.RoleID, cfg.SecretID); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("vault: no authentication method provided")
	}

	return &Client{api: raw}, nil
}

func loginAppRole(c *vaultapi.Client, roleID, secretID string) error {
	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}
	secret, err := c.Logical().Write("auth/approle/login", data)
	if err != nil {
		return fmt.Errorf("vault: approle login: %w", err)
	}
	if secret == nil || secret.Auth == nil {
		return fmt.Errorf("vault: approle login returned no auth info")
	}
	c.SetToken(secret.Auth.ClientToken)
	return nil
}
