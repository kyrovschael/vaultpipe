package env

import (
	"context"
	"strings"
	"testing"
)

func TestGroupSnapshot_DefaultKeyFn(t *testing.T) {
	src := StaticSource(Snapshot{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
		"NOUNDERSCORE": "value",
	})

	groups, err := GroupSnapshot(context.Background(), src, DefaultGroupOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if groups["DB"]["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST in DB group, got %v", groups["DB"])
	}
	if groups["DB"]["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT in DB group")
	}
	if groups["APP"]["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV in APP group")
	}
	if groups["NOUNDERSCORE"]["NOUNDERSCORE"] != "value" {
		t.Errorf("expected NOUNDERSCORE in its own group")
	}
}

func TestGroupSnapshot_CustomKeyFn(t *testing.T) {
	src := StaticSource(Snapshot{
		"VAULT_ADDR":  "https://vault",
		"VAULT_TOKEN": "s.abc",
		"OTHER_KEY":   "val",
	})

	opts := GroupOptions{
		KeyFn: func(key string) string {
			if strings.HasPrefix(key, "VAULT_") {
				return "vault"
			}
			return "other"
		},
	}

	groups, err := GroupSnapshot(context.Background(), src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups["vault"]) != 2 {
		t.Errorf("expected 2 vault keys, got %d", len(groups["vault"]))
	}
	if len(groups["other"]) != 1 {
		t.Errorf("expected 1 other key, got %d", len(groups["other"]))
	}
}

func TestGroupSnapshot_DoesNotMutateGroups(t *testing.T) {
	src := StaticSource(Snapshot{"DB_HOST": "localhost"})
	groups, _ := GroupSnapshot(context.Background(), src, DefaultGroupOptions())

	groups["DB"]["DB_HOST"] = "mutated"

	groups2, _ := GroupSnapshot(context.Background(), src, DefaultGroupOptions())
	if groups2["DB"]["DB_HOST"] != "localhost" {
		t.Error("mutation of group affected source")
	}
}

func TestGroupSnapshot_EmptyKeyFallback(t *testing.T) {
	src := StaticSource(Snapshot{"NOPREFIX": "v"})
	opts := GroupOptions{KeyFn: func(key string) string { return "" }}

	groups, err := GroupSnapshot(context.Background(), src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if groups["_"]["NOPREFIX"] != "v" {
		t.Errorf("expected fallback group '_', got %v", groups)
	}
}

func TestGroupNames_Sorted(t *testing.T) {
	groups := map[string]Snapshot{
		"Z": {"Z_A": "1"},
		"A": {"A_B": "2"},
		"M": {"M_C": "3"},
	}
	names := GroupNames(groups)
	if names[0] != "A" || names[1] != "M" || names[2] != "Z" {
		t.Errorf("expected sorted names, got %v", names)
	}
}
