// Package health provides Vault connectivity and auth health checks.
package health

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Status holds the result of a health check.
type Status struct {
	Healthy   bool
	Sealed    bool
	Standby   bool
	Version   string
	Latency   time.Duration
	Error     error
}

// Checker performs health checks against a Vault server.
type Checker struct {
	client *vaultapi.Client
}

// NewChecker creates a Checker wrapping the given Vault client.
func NewChecker(client *vaultapi.Client) *Checker {
	return &Checker{client: client}
}

// Check queries the Vault health endpoint and returns a Status.
func (c *Checker) Check(ctx context.Context) Status {
	start := time.Now()

	sys := c.client.Sys()
	health, err := sys.HealthWithContext(ctx)
	latency := time.Since(start)

	if err != nil {
		return Status{
			Healthy: false,
			Latency: latency,
			Error:   fmt.Errorf("vault health check failed: %w", err),
		}
	}

	return Status{
		Healthy: health.Initialized && !health.Sealed,
		Sealed:  health.Sealed,
		Standby: health.Standby,
		Version: health.Version,
		Latency: latency,
	}
}
