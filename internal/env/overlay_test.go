package env

import (
	"testing"
)

func TestOverlaySnapshot_OverwriteExisting(t *testing.T) {
	base := FromSlice([]string{"A=1", "B=2"})
	overlay := FromSlice([]string{"B=99", "C=3"})
	out := OverlaySnapshot(base, DefaultOverlayOptions(), overlay)
	assertVal(t, out, "A", "1")
	assertVal(t, out, "B", "99")
	assertVal(t, out, "C", "3")
}

func TestOverlaySnapshot_NoOverwrite(t *testing.T) {
	base := FromSlice([]string{"A=1", "B=2"})
	overlay := FromSlice([]string{"B=99", "C=3"})
	opts := OverlayOptions{OverwriteExisting: false}
	out := OverlaySnapshot(base, opts, overlay)
	assertVal(t, out, "B", "2") // base wins
	assertVal(t, out, "C", "3") // new key still added
}

func TestOverlaySnapshot_SkipEmpty(t *testing.T) {
	base := FromSlice([]string{"A=hello"})
	overlay := FromSlice([]string{"A="})
	opts := OverlayOptions{OverwriteExisting: true, SkipEmpty: true}
	out := OverlaySnapshot(base, opts, overlay)
	assertVal(t, out, "A", "hello") // empty overlay skipped
}

func TestOverlaySnapshot_DoesNotMutateBase(t *testing.T) {
	base := FromSlice([]string{"X=original"})
	overlay := FromSlice([]string{"X=changed"})
	_ = OverlaySnapshot(base, DefaultOverlayOptions(), overlay)
	assertVal(t, base, "X", "original")
}

func TestOverlaySnapshot_MultipleLayers(t *testing.T) {
	base := FromSlice([]string{"A=1"})
	l1 := FromSlice([]string{"A=2"})
	l2 := FromSlice([]string{"A=3"})
	out := OverlaySnapshot(base, DefaultOverlayOptions(), l1, l2)
	assertVal(t, out, "A", "3") // last layer wins
}

func assertVal(t *testing.T, s Snapshot, key, want string) {
	t.Helper()
	v, ok := s.Lookup(key)
	if !ok {
		t.Fatalf("key %q not found", key)
	}
	if v != want {
		t.Fatalf("key %q: got %q, want %q", key, v, want)
	}
}
