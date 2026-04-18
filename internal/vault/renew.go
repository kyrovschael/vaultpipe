package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/api"
)

// RenewToken attempts to renew the Vault token and schedules background renewal.
// It returns a cancel function to stop the renewal loop.
func RenewToken(ctx context.Context, client *api.Client) (context.CancelFunc, error) {
	secret, err := client.Auth().Token().RenewSelf(0)
	if err != nil {
		return nil, err
	}

	ttl, err := secret.TokenTTL()
	if err != nil || ttl == 0 {
		// Non-renewable token; no-op cancel
		return func() {}, nil
	}

	renewCtx, cancel := context.WithCancel(ctx)
	go renewLoop(renewCtx, client, ttl)
	return cancel, nil
}

func renewLoop(ctx context.Context, client *api.Client, ttl time.Duration) {
	// Renew at 2/3 of TTL to avoid expiry
	interval := ttl * 2 / 3
	if interval < 5*time.Second {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s// Token failed; exit loop
			tnewTTL, err := secret.TokenTTL()
			if err != nil || newTTL == 0 {
				return
			}
			newInterval := newTTL * 2 / 3
			if newInterval < 5*time.Second {
				newInterval = 5 * time.Second
			}
			ticker.Reset(newInterval)
		}
	}
}
