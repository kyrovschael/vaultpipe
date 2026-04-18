package audit

// Well-known event names used throughout vaultpipe.
const (
	// EventSecretsLoaded is emitted after secrets are successfully fetched.
	EventSecretsLoaded = "secrets.loaded"

	// EventSecretsLoadFailed is emitted when fetching secrets fails.
	EventSecretsLoadFailed = "secrets.load_failed"

	// EventProcessStart is emitted just before the child process is exec'd.
	EventProcessStart = "process.start"

	// EventProcessExit is emitted after the child process exits.
	EventProcessExit = "process.exit"

	// EventAuthSuccess is emitted after successful Vault authentication.
	EventAuthSuccess = "auth.success"

	// EventAuthFailed is emitted when Vault authentication fails.
	EventAuthFailed = "auth.failed"
)
