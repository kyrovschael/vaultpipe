package vault

import (
	"fmt"
	"strings"
)

// GetSecrets reads a KV secret at the given path and returns a map of
// key→value strings suitable for injecting into a process environment.
// Supports both KV v1 and KV v2 mounts.
func (c *Client) GetSecrets(path string) (map[string]string, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at %q", path)
	}

	data := secret.Data

	// KV v2 wraps values under a nested "data" key.
	if nested, ok := data["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			data = nestedMap
		}
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}

// EnvKey converts a Vault secret key to a canonical environment variable name.
// e.g. "db-password" → "DB_PASSWORD"
func EnvKey(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}

// SecretsToEnv converts a secrets map to a slice of "KEY=value" strings.
func SecretsToEnv(secrets map[string]string) []string {
	env := make([]string, 0, len(secrets))
	for k, v := range secrets {
		env = append(env, EnvKey(k)+"="+v)
	}
	return env
}

// MergeEnv merges Vault-derived env pairs into an existing environment slice.
// Vault values take precedence: any existing entry with the same key is replaced.
func MergeEnv(base []string, overrides []string) []string {
	keys := make(map[string]struct{}, len(overrides))
	for _, pair := range overrides {
		if idx := strings.IndexByte(pair, '='); idx > 0 {
			keys[pair[:idx]] = struct{}{}
		}
	}
	result := make([]string, 0, len(base)+len(overrides))
	for _, pair := range base {
		if idx := strings.IndexByte(pair, '='); idx > 0 {
			if _, overridden := keys[pair[:idx]]; overridden {
				continue
			}
		}
		result = append(result, pair)
	}
	return append(result, overrides...)
}

// FilterSecrets returns a new map containing only the keys present in the
// allowlist. This is useful when a Vault path contains more keys than a
// process needs, reducing the surface area of the injected environment.
func FilterSecrets(secrets map[string]string, allow []string) map[string]string {
	allowSet := make(map[string]struct{}, len(allow))
	for _, k := range allow {
		allowSet[k] = struct{}{}
	}
	filtered := make(map[string]string, len(allow))
	for k, v := range secrets {
		if _, ok := allowSet[k]; ok {
			filtered[k] = v
		}
	}
	return filtered
}
