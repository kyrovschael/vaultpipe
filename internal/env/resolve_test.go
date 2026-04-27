package env

import (
	"errors"
	"testing"
)

func TestResolveSnapshot_FirstSourceWins(t *testing.T) {
	primary := StaticSource(Snapshot{data: map[string]string{"FOO": "from_primary", "BAR": "bar_primary"}})
	fallback := StaticSource(Snapshot{data: map[string]string{"FOO": "from_fallback", "BAZ": "baz_fallback"}})

	got, err := ResolveSnapshot([]string{"FOO", "BAR", "BAZ"}, DefaultResolveOptions(), primary, fallback)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.data["FOO"] != "from_primary" {
		t.Errorf("FOO: want %q, got %q", "from_primary", got.data["FOO"])
	}
	if got.data["BAR"] != "bar_primary" {
		t.Errorf("BAR: want %q, got %q", "bar_primary", got.data["BAR"])
	}
	if got.data["BAZ"] != "baz_fallback" {
		t.Errorf("BAZ: want %q, got %q", "baz_fallback", got.data["BAZ"])
	}
}

func TestResolveSnapshot_MissingKeyLenient(t *testing.T) {
	src := StaticSource(Snapshot{data: map[string]string{"FOO": "bar"}})
	got, err := ResolveSnapshot([]string{"FOO", "MISSING"}, DefaultResolveOptions(), src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got.data["MISSING"]; ok {
		t.Error("MISSING key should be absent in lenient mode")
	}
}

func TestResolveSnapshot_MissingKeyStrict(t *testing.T) {
	src := StaticSource(Snapshot{data: map[string]string{"FOO": "bar"}})
	opts := ResolveOptions{Strict: true}
	_, err := ResolveSnapshot([]string{"FOO", "MISSING"}, opts, src)
	if err == nil {
		t.Fatal("expected error in strict mode for missing key")
	}
}

func TestResolveSnapshot_NoSources_Strict(t *testing.T) {
	_, err := ResolveSnapshot([]string{"FOO"}, ResolveOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error with no sources in strict mode")
	}
}

func TestResolveSnapshot_NoSources_Lenient(t *testing.T) {
	got, err := ResolveSnapshot([]string{"FOO"}, DefaultResolveOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.data) != 0 {
		t.Errorf("expected empty snapshot, got %v", got.data)
	}
}

func TestResolveSnapshot_SourceError(t *testing.T) {
	errSrc := errorSource{err: errors.New("boom")}
	_, err := ResolveSnapshot([]string{"FOO"}, DefaultResolveOptions(), errSrc)
	if err == nil {
		t.Fatal("expected error from failing source")
	}
}

func TestResolveSnapshot_SkipsEmptyValue(t *testing.T) {
	primary := StaticSource(Snapshot{data: map[string]string{"FOO": ""}})
	secondary := StaticSource(Snapshot{data: map[string]string{"FOO": "filled"}})

	got, err := ResolveSnapshot([]string{"FOO"}, DefaultResolveOptions(), primary, secondary)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.data["FOO"] != "filled" {
		t.Errorf("want %q, got %q", "filled", got.data["FOO"])
	}
}

// errorSource is a Source that always returns an error.
type errorSource struct{ err error }

func (e errorSource) Snapshot() (Snapshot, error) { return Snapshot{}, e.err }
