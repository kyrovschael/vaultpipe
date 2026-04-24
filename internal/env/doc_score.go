// Package env provides environment variable manipulation utilities.
//
// # ScoreSnapshot
//
// ScoreSnapshot assigns a numeric score to each entry in a snapshot using a
// caller-supplied scoring function. Entries are returned sorted by score in
// descending order (highest score first).
//
// The original snapshot is never mutated.
//
//	scored, err := env.ScoreSnapshot(snap, env.DefaultScoreOptions(), func(k, v string) float64 {
//		if strings.HasPrefix(k, "VAULT_") {
//			return 10.0
//		}
//		return 1.0
//	})
package env
