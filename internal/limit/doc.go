// Package limit provides a concurrency limiter for outbound Vault requests.
//
// Use New to create a Limiter from a Config, then call Acquire before each
// Vault API call and invoke the returned release function when done:
//
//	release, err := limiter.Acquire(ctx)
//	if err != nil {
//	    return err
//	}
//	defer release()
//
// WithTimeout wraps a context with the configured per-request deadline so
// individual fetches cannot stall indefinitely.
package limit
