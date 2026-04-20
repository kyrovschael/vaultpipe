package env

import (
	"testing"
)

func TestJoinSnapshot_JoinsSharedKeys(t *testing.T) {
	a := Snapshot{{Key: "PATH", Value: "/usr/bin"}, {Key: "HOME", Value: "/root"}}
	b := Snapshot{{Key: "PATH", Value: "/opt/bin"}}

	result := JoinSnapshot([]Snapshot{a, b}, DefaultJoinOptions())

	got := map[string]string{}
	for _, e := range result {
		got[e.Key] = e.Value
	}

	if got["PATH"] != "/usr/bin,/opt/bin" {
		t.Errorf("PATH: got %q, want %q", got["PATH"], "/usr/bin,/opt/bin")
	}
	if got["HOME"] != "/root" {
		t.Errorf("HOME: got %q, want %q", got["HOME"], "/root")
	}
}

func TestJoinSnapshot_CustomSeparator(t *testing.T) {
	a := Snapshot{{Key: "PLUGINS", Value: "a"}}
	b := Snapshot{{Key: "PLUGINS", Value: "b"}}

	opts := DefaultJoinOptions()
	opts.Separator = ":"
	result := JoinSnapshot([]Snapshot{a, b}, opts)

	if len(result) != 1 || result[0].Value != "a:b" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestJoinSnapshot_RestrictedKeys(t *testing.T) {
	a := Snapshot{{Key: "A", Value: "1"}, {Key: "B", Value: "x"}}
	b := Snapshot{{Key: "A", Value: "2"}, {Key: "B", Value: "y"}}

	opts := DefaultJoinOptions()
	opts.Keys = []string{"A"} // only join A; B should take last value
	result := JoinSnapshot([]Snapshot{a, b}, opts)

	got := map[string]string{}
	for _, e := range result {
		got[e.Key] = e.Value
	}

	if got["A"] != "1,2" {
		t.Errorf("A: got %q, want %q", got["A"], "1,2")
	}
	if got["B"] != "y" {
		t.Errorf("B: got %q, want %q", got["B"], "y")
	}
}

func TestJoinSnapshot_SingleSnapshot(t *testing.T) {
	s := Snapshot{{Key: "FOO", Value: "bar"}}
	result := JoinSnapshot([]Snapshot{s}, DefaultJoinOptions())
	if len(result) != 1 || result[0].Value != "bar" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestJoinSnapshot_DoesNotMutateOriginals(t *testing.T) {
	a := Snapshot{{Key: "X", Value: "1"}}
	b := Snapshot{{Key: "X", Value: "2"}}
	origA := a[0].Value

	_ = JoinSnapshot([]Snapshot{a, b}, DefaultJoinOptions())

	if a[0].Value != origA {
		t.Errorf("original snapshot mutated")
	}
}

func TestJoinSnapshot_EmptyInput(t *testing.T) {
	result := JoinSnapshot(nil, DefaultJoinOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
