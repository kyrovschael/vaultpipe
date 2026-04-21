package env

import (
	"strings"
	"testing"
)

func TestPartitionSnapshot_BasicSplit(t *testing.T) {
	src := Snapshot{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_ENV": "prod"}
	matched, rest := PartitionSnapshot(src, DefaultPartitionOptions(), func(k, _ string) bool {
		return strings.HasPrefix(k, "DB_")
	})

	if len(matched) != 2 {
		t.Fatalf("want 2 matched, got %d", len(matched))
	}
	if len(rest) != 1 {
		t.Fatalf("want 1 rest, got %d", len(rest))
	}
	if _, ok := matched["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in matched")
	}
	if _, ok := rest["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in rest")
	}
}

func TestPartitionSnapshot_AllMatch(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2"}
	matched, rest := PartitionSnapshot(src, DefaultPartitionOptions(), func(_, _ string) bool { return true })
	if len(matched) != 2 {
		t.Fatalf("want 2, got %d", len(matched))
	}
	if len(rest) != 0 {
		t.Fatalf("want 0 rest, got %d", len(rest))
	}
}

func TestPartitionSnapshot_NoneMatch(t *testing.T) {
	src := Snapshot{"A": "1", "B": "2"}
	matched, rest := PartitionSnapshot(src, DefaultPartitionOptions(), func(_, _ string) bool { return false })
	if len(matched) != 0 {
		t.Fatalf("want 0 matched, got %d", len(matched))
	}
	if len(rest) != 2 {
		t.Fatalf("want 2 rest, got %d", len(rest))
	}
}

func TestPartitionSnapshot_DoesNotMutateSource(t *testing.T) {
	src := Snapshot{"X": "original"}
	matched, _ := PartitionSnapshot(src, DefaultPartitionOptions(), func(_, _ string) bool { return true })
	matched["X"] = "mutated"
	if src["X"] != "original" {
		t.Error("source snapshot was mutated")
	}
}

func TestPartitionSnapshot_CaseFold(t *testing.T) {
	src := Snapshot{"db_host": "localhost", "APP_ENV": "prod"}
	opts := DefaultPartitionOptions()
	opts.CaseFold = true

	matched, rest := PartitionSnapshot(src, opts, func(k, _ string) bool {
		return strings.HasPrefix(k, "DB_")
	})

	if len(matched) != 1 {
		t.Fatalf("want 1 matched with case fold, got %d", len(matched))
	}
	if _, ok := matched["db_host"]; !ok {
		t.Error("expected db_host in matched (original case preserved)")
	}
	if len(rest) != 1 {
		t.Fatalf("want 1 rest, got %d", len(rest))
	}
}

func TestPartitionSnapshot_EmptySnapshot(t *testing.T) {
	matched, rest := PartitionSnapshot(Snapshot{}, DefaultPartitionOptions(), func(_, _ string) bool { return true })
	if len(matched) != 0 || len(rest) != 0 {
		t.Error("expected both partitions empty for empty source")
	}
}
