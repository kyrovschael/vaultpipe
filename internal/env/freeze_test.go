package env

import (
	"sync"
	"testing"
)

func TestFreezeSnapshot_Get(t *testing.T) {
	snap := Snapshot{"KEY": "value"}
	f := FreezeSnapshot(snap)
	v, ok := f.Get("KEY")
	if !ok || v != "value" {
		t.Fatalf("expected value=\"value\" ok=true, got %q %v", v, ok)
	}
	_, ok = f.Get("MISSING")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestFreezeSnapshot_DoesNotMutateOriginal(t *testing.T) {
	snap := Snapshot{"A": "1"}
	f := FreezeSnapshot(snap)
	snap["A"] = "mutated"
	v, _ := f.Get("A")
	if v != "1" {
		t.Fatalf("frozen snapshot was mutated; got %q", v)
	}
}

func TestFreezeSnapshot_SnapshotReturnsClone(t *testing.T) {
	f := FreezeSnapshot(Snapshot{"X": "y"})
	clone := f.Snapshot()
	clone["X"] = "mutated"
	v, _ := f.Get("X")
	if v != "y" {
		t.Fatalf("clone mutation affected frozen snapshot; got %q", v)
	}
}

func TestFreezeSnapshot_Len(t *testing.T) {
	f := FreezeSnapshot(Snapshot{"A": "1", "B": "2"})
	if f.Len() != 2 {
		t.Fatalf("expected len 2, got %d", f.Len())
	}
}

func TestFreezeSnapshot_AssertKey(t *testing.T) {
	f := FreezeSnapshot(Snapshot{"PRESENT": "yes"})
	if err := f.AssertKey("PRESENT"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.AssertKey("ABSENT"); err == nil {
		t.Fatal("expected error for absent key")
	}
}

func TestFreezeSnapshot_ConcurrentReads(t *testing.T) {
	f := FreezeSnapshot(Snapshot{"SAFE": "yes"})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, ok := f.Get("SAFE")
			if !ok || v != "yes" {
				t.Errorf("unexpected result: %q %v", v, ok)
			}
		}()
	}
	wg.Wait()
}
