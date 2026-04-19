package env

import (
	"context"
	"errors"
	"testing"
)

func TestStaticSource_ReturnsClone(t *testing.T) {
	original := Snapshot{"KEY": "val"}
	src := StaticSource(original)
	got, err := src(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "val" {
		t.Fatalf("expected val, got %q", got["KEY"])
	}
	// mutation of result must not affect original
	got["KEY"] = "mutated"
	if original["KEY"] != "val" {
		t.Fatal("StaticSource did not return a clone")
	}
}

func TestSliceSource_ParsesPairs(t *testing.T) {
	src := SliceSource([]string{"A=1", "B=2", "NOEQUALS"})
	got, err := src(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "1" || got["B"] != "2" {
		t.Fatalf("unexpected snapshot: %v", got)
	}
}

func TestChainSource_LaterWins(t *testing.T) {
	a := StaticSource(Snapshot{"X": "a", "Y": "a"})
	b := StaticSource(Snapshot{"X": "b"})
	src := ChainSource(a, b)
	got, err := src(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["X"] != "b" {
		t.Fatalf("expected b to win, got %q", got["X"])
	}
	if got["Y"] != "a" {
		t.Fatalf("expected Y=a, got %q", got["Y"])
	}
}

func TestChainSource_PropagatesError(t *testing.T) {
	boom := func(_ context.Context) (Snapshot, error) {
		return nil, errors.New("boom")
	}
	src := ChainSource(StaticSource(Snapshot{"K": "v"}), boom)
	_, err := src(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMapSource_ReturnsAllKeys(t *testing.T) {
	src := MapSource(map[string]string{"FOO": "bar", "BAZ": "qux"})
	got, err := src(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
}

func TestOSSource_ReturnsNonEmpty(t *testing.T) {
	src := OSSource()
	got, err := src(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected non-empty OS snapshot")
	}
}
