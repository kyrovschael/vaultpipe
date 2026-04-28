package env

import (
	"testing"
)

func TestNormalizeSnapshot_UpperCasesKeys(t *testing.T) {
	src := StaticSource(Snapshot{
		"db_host": "localhost",
		"db_port": "5432",
	})
	out, err := NormalizeSnapshot(src, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key after normalisation")
	}
	if _, ok := out["DB_PORT"]; !ok {
		t.Error("expected DB_PORT key after normalisation")
	}
}

func TestNormalizeSnapshot_DropsEmptyValues(t *testing.T) {
	src := StaticSource(Snapshot{
		"KEEP": "value",
		"DROP": "",
		"TRIM": "   ",
	})
	out, err := NormalizeSnapshot(src, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DROP"]; ok {
		t.Error("expected DROP to be removed")
	}
	if _, ok := out["TRIM"]; ok {
		t.Error("expected TRIM (whitespace-only) to be removed")
	}
	if v := out["KEEP"]; v != "value" {
		t.Errorf("expected KEEP=value, got %q", v)
	}
}

func TestNormalizeSnapshot_SkipUpperKeys(t *testing.T) {
	src := StaticSource(Snapshot{
		"lower_key": "hello",
	})
	opts := DefaultNormalizeOptions()
	opts.SkipUpperKeys = true
	out, err := NormalizeSnapshot(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["lower_key"]; !ok {
		t.Error("expected lower_key to be preserved when SkipUpperKeys=true")
	}
	if _, ok := out["LOWER_KEY"]; ok {
		t.Error("did not expect LOWER_KEY when SkipUpperKeys=true")
	}
}

func TestNormalizeSnapshot_DoesNotMutateSource(t *testing.T) {
	original := Snapshot{
		"my_key": "val",
	}
	src := StaticSource(original)
	_, err := NormalizeSnapshot(src, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := original["my_key"]; !ok {
		t.Error("source snapshot was mutated")
	}
	if _, ok := original["MY_KEY"]; ok {
		t.Error("source snapshot was mutated (upper-cased key appeared)")
	}
}

func TestNormalizeSnapshot_SanitizesKeys(t *testing.T) {
	src := StaticSource(Snapshot{
		"valid-but-dashed": "ok",
	})
	out, err := NormalizeSnapshot(src, DefaultNormalizeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Sanitize replaces '-' with '_', then upper-case produces VALID_BUT_DASHED.
	if _, ok := out["VALID_BUT_DASHED"]; !ok {
		t.Errorf("expected sanitized+upper key VALID_BUT_DASHED, got keys: %v", snapshotKeys(out))
	}
}
