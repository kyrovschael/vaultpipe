package env

import "sort"

// DefaultChunkOptions returns a ChunkOptions with a batch size of 100.
func DefaultChunkOptions() ChunkOptions {
	return ChunkOptions{Size: 100}
}

// ChunkOptions controls how ChunkSnapshot splits a snapshot.
type ChunkOptions struct {
	// Size is the maximum number of entries per chunk.
	// Must be >= 1; values less than 1 are treated as 1.
	Size int

	// Sorted controls whether keys are sorted before chunking so that
	// the output is deterministic across calls.
	Sorted bool
}

// ChunkSnapshot splits snap into successive slices of at most
// opts.Size entries. Each chunk is an independent Snapshot clone.
// The original snapshot is never mutated.
func ChunkSnapshot(snap Snapshot, opts ChunkOptions) []Snapshot {
	if opts.Size < 1 {
		opts.Size = 1
	}

	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}

	if opts.Sorted {
		sort.Strings(keys)
	}

	if len(keys) == 0 {
		return nil
	}

	var chunks []Snapshot
	for i := 0; i < len(keys); i += opts.Size {
		end := i + opts.Size
		if end > len(keys) {
			end = len(keys)
		}
		chunk := make(Snapshot, end-i)
		for _, k := range keys[i:end] {
			chunk[k] = snap[k]
		}
		chunks = append(chunks, chunk)
	}
	return chunks
}
