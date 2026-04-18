// Package rotate implements secret rotation detection for vaultpipe.
//
// A Watcher polls a SecretFetcher at a configurable interval and calls a
// ChangeHandler whenever the set of secrets differs from the previous poll.
// Detection is based on a SHA-256 signature of the key-value pairs, so any
// added, removed, or updated secret triggers a notification.
//
// Typical usage:
//
//	w := rotate.NewWatcher(fetchFn, reloadFn, 30*time.Second)
//	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package rotate
