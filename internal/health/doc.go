// Package health exposes utilities for checking the liveness and readiness
// of a HashiCorp Vault server before vaultpipe injects secrets into a
// child process.
//
// Usage:
//
//	checker := health.NewChecker(vaultClient)
//	status := health.WaitUntilHealthy(ctx, checker, health.DefaultRetryConfig())
//	if !status.Healthy {
//		log.Fatalf("vault not healthy: %v", status.Error)
//	}
package health
