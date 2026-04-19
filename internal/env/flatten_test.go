package env

import (
	"testing"
)

func TestFlattenMap_Simple(t *testing.T) {
	src := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	}
	got := FlattenMap(src, DefaultFlattenOptions())
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
}

func TestFlattenMap_DeepNesting(t *testing.T) {
	src := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	got := FlattenMap(src, DefaultFlattenOptions())
	if got["A_B_C"] != "deep" {
		t.Errorf("expected A_B_C=deep, got %q", got["A_B_C"])
	}
}

func TestFlattenMap_SkipsNonString(t *testing.T) {
	src := map[string]any{
		"count": 42,
		"name":  "vault",
	}
	got := FlattenMap(src, DefaultFlattenOptions())
	if _, ok := got["COUNT"]; ok {
		t.Error("expected non-string value to be skipped")
	}
	if got["NAME"] != "vault" {
		t.Errorf("expected NAME=vault, got %q", got["NAME"])
	}
}

func TestFlattenMap_CustomSeparator(t *testing.T) {
	opts := FlattenOptions{Separator: ".", UpperKeys: false}
	src := map[string]any{
		"db": map[string]any{"host": "h"},
	}
	got := FlattenMap(src, opts)
	if got["db.host"] != "h" {
		t.Errorf("expected db.host=h, got %v", got)
	}
}

func TestFlattenMap_Empty(t *testing.T) {
	got := FlattenMap(map[string]any{}, DefaultFlattenOptions())
	if len(got) != 0 {
		t.Errorf("expected empty snapshot, got %v", got)
	}
}
