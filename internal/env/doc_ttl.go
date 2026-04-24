// Package env – TTL support
//
// TTLSnapshot stamps each key in a Snapshot with an absolute expiry
// time. Callers can poll TTLEntry.Expired to detect stale values and
// trigger a refresh from Vault.
//
// Example:
//
//	entries, err := env.TTLSnapshot(snap, env.TTLOptions{
//		TTL:  30 * time.Second,
//		Keys: []string{"DB_PASSWORD", "API_KEY"},
//	})
//	if err != nil { ... }
//	for k, e := range entries {
//		if e.Expired(time.Now()) {
//			// re-fetch secret for k
//		}
//	}
package env
