package env

import (
	"errors"
	"sort"
)

// ScoreFunc computes a numeric score for a single environment entry.
// Higher scores sort earlier in the result.
type ScoreFunc func(key, value string) float64

// ScoreOptions controls the behaviour of ScoreSnapshot.
type ScoreOptions struct {
	// Ascending reverses the default sort order so that lower scores sort first.
	Ascending bool

	// Keys restricts scoring to the given keys. When empty all keys are scored.
	Keys []string
}

// DefaultScoreOptions returns a ScoreOptions with sensible defaults.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{}
}

type scoredEntry struct {
	Entry
	score float64
}

// ScoreSnapshot returns a new snapshot whose entries are ordered by the score
// returned by fn. Entries not in opts.Keys (when non-empty) are appended after
// scored entries, preserving their original relative order.
func ScoreSnapshot(s Snapshot, opts ScoreOptions, fn ScoreFunc) (Snapshot, error) {
	if fn == nil {
		return Snapshot{}, errors.New("env: ScoreSnapshot: fn must not be nil")
	}

	allowed := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowed[k] = true
	}

	var scored []scoredEntry
	var rest []Entry

	for _, e := range s.Entries {
		if len(allowed) == 0 || allowed[e.Key] {
			scored = append(scored, scoredEntry{Entry: e, score: fn(e.Key, e.Value)})
		} else {
			rest = append(rest, e)
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if opts.Ascending {
			return scored[i].score < scored[j].score
		}
		return scored[i].score > scored[j].score
	})

	out := make([]Entry, 0, len(s.Entries))
	for _, se := range scored {
		out = append(out, se.Entry)
	}
	out = append(out, rest...)

	return Snapshot{Entries: out}, nil
}
