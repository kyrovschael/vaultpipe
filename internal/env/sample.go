package env

import (
	"fmt"
	"math/rand"
	"sort"
)

// DefaultSampleOptions returns a SampleOptions with safe defaults.
func DefaultSampleOptions() SampleOptions {
	return SampleOptions{
		N:    10,
		Seed: 0,
	}
}

// SampleOptions controls the behaviour of SampleSnapshot.
type SampleOptions struct {
	// N is the maximum number of entries to return.
	// If N >= len(snapshot) all entries are returned.
	N int

	// Seed is the random seed used for selection.
	// A fixed seed produces a deterministic result.
	Seed int64
}

// SampleSnapshot returns a pseudo-random subset of at most opts.N entries
// from src. The source snapshot is never mutated. Keys are sorted before
// sampling so that the result is deterministic for a given seed.
func SampleSnapshot(src Snapshot, opts SampleOptions) (Snapshot, error) {
	if opts.N < 0 {
		return nil, fmt.Errorf("env: SampleSnapshot: N must be >= 0, got %d", opts.N)
	}

	clone := CloneSnapshot(src, DefaultCloneOptions())

	keys := make([]string, 0, len(clone))
	for k := range clone {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if opts.N == 0 || len(keys) == 0 {
		return Snapshot{}, nil
	}

	//nolint:gosec // non-cryptographic sampling is intentional
	r := rand.New(rand.NewSource(opts.Seed))
	r.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	n := opts.N
	if n > len(keys) {
		n = len(keys)
	}

	out := make(Snapshot, n)
	for _, k := range keys[:n] {
		out[k] = clone[k]
	}
	return out, nil
}
