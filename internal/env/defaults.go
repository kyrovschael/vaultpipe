package env

// DefaultDenyPatterns is the set of environment variable patterns that
// vaultpipe blocks by default to prevent credential leakage into child
// processes.
var DefaultDenyPatterns = []string{
	// Vault credentials
	"VAULT_TOKEN",
	"VAULT_ROLE_ID",
	"VAULT_SECRET_ID",

	// AWS credentials
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_SESSION_TOKEN",
	"AWS_SECURITY_TOKEN",

	// GCP / GKE
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_CREDENTIALS",

	// GitHub Actions
	"GITHUB_TOKEN",
	"ACTIONS_RUNTIME_TOKEN",

	// Generic secrets
	"SECRET_*",
	"PRIVATE_*",
}

// DefaultDenyList returns a DenyList pre-populated with DefaultDenyPatterns.
func DefaultDenyList() *DenyList {
	return NewDenyList(DefaultDenyPatterns)
}
