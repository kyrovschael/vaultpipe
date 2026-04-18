// Package env provides helpers for manipulating the environment variable set
// that is forwarded to child processes launched by vaultpipe.
//
// The primary concern is safety: sensitive variables such as Vault credentials
// or cloud provider keys present in the parent process environment must not
// leak into the child unless explicitly requested.
//
// Usage:
//
//	dl := env.NewDenyList([]string{"VAULT_*", "AWS_*", "GITHUB_TOKEN"})
//	clean := dl.Filter(os.Environ())
package env
