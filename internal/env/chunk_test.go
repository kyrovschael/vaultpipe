package env

import (
	"testing"
)

func TestChunkSnapshot_EvenSplit(t *testing.T) {
	snap := Snapshot{"A": "1", "B": "2", "C": "3", "D": "4"}
	opts := ChunkOptions{Size: 2, Sorted: true}
	chunks := ChunkSnapshot(snap, opts)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 || len(chunks[1]) != 2 {
		t.Fatalf("unexpected chunk sizes: %d, %d", len(chunks[0]), len(chunks[1]))
	}
}

func TestChunkSnapshot_UnevenSplit(t *testing.T) {
	snap := Snapshot{"A": "1", "B": "2", "C": "3"}
	opts := ChunkOptions{Size: 2, Sorted: true}
	chunks := ChunkSnapshot(snap, opts)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if len(chunks[1]) != 1 {
		t.Fatalf("last chunk should have 1 entry, got %d", len(chunks[1]))
	}
}

func TestChunkSnapshot_SizeOne(t *testing.T) {
	snap := Snapshot{"X": "1", "Y": "2", "Z": "3"}
	opts := ChunkOptions{Size: 1, Sorted: true}
	chunks := ChunkSnapshot(snap, opts)
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
}

func TestChunkSnapshot_EmptySnapshot(t *testing.T) {
	chunks := ChunkSnapshot(Snapshot{}, DefaultChunkOptions())
	if chunks != nil {
		t.Fatalf("expected nil for empty snapshot, got %v", chunks)
	}
}

func TestChunkSnapshot_DoesNotMutateOriginal(t *testing.T) {
	snap := Snapshot{"A": "original"}
	opts := ChunkOptions{Size: 10, Sorted: true}
	chunks := ChunkSnapshot(snap, opts)
	chunks[0]["A"] = "mutated"
	if snap["A"] != "original" {
		t.Fatal("original snapshot was mutated")
	}
}

func TestChunkSnapshot_InvalidSizeTreatedAsOne(t *testing.T) {
	snap := Snapshot{"A": "1", "B": "2"}
	opts := ChunkOptions{Size: 0, Sorted: true}
	chunks := ChunkSnapshot(snap, opts)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks with size=0 treated as 1, got %d", len(chunks))
	}
}

func TestChunkSnapshot_LargerThanSnapshot(t *testing.T) {
	snap := Snapshot{"A": "1", "B": "2"}
	opts := ChunkOptions{Size: 100, Sorted: true}
	chunks := ChunkSnapshot(snap, opts)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 {
		t.Fatalf("expected 2 entries in single chunk, got %d", len(chunks[0]))
	}
}
