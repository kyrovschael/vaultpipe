package env

import (
	"strings"
	"testing"
)

func TestTransformSnapshot_UpperKeys(t *testing.T) {
	s := Snapshot{"foo": "bar", "baz": "qux"}
	out := TransformSnapshot(s, UpperKeys())
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestTransformSnapshot_LowerKeys(t *testing.T) {
	s := Snapshot{"FOO": "1", "BAR": "2"}
	out := TransformSnapshot(s, LowerKeys())
	if out["foo"] != "1" || out["bar"] != "2" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestTransformSnapshot_TrimValues(t *testing.T) {
	s := Snapshot{"KEY": "  hello  ", "OTHER": "\tworld\n"}
	out := TransformSnapshot(s, TrimValues())
	if out["KEY"] != "hello" || out["OTHER"] != "world" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestTransformSnapshot_DoesNotMutateOriginal(t *testing.T) {
	s := Snapshot{"key": "value"}
	_ = TransformSnapshot(s, UpperKeys())
	if _, ok := s["key"]; !ok {
		t.Fatal("original snapshot was mutated")
	}
}

func TestTransformSnapshot_BothFns(t *testing.T) {
	s := Snapshot{"foo": "  bar  "}
	opts := TransformOptions{
		KeyFn:   strings.ToUpper,
		ValueFn: strings.TrimSpace,
	}
	out := TransformSnapshot(s, opts)
	if out["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %v", out)
	}
}

func TestTransformSnapshot_NilFns(t *testing.T) {
	s := Snapshot{"k": "v"}
	out := TransformSnapshot(s, DefaultTransformOptions())
	if out["k"] != "v" {
		t.Fatalf("expected k=v, got %v", out)
	}
}

func TestTransformSnapshot_KeyCollision(t *testing.T) {
	s := Snapshot{"foo": "a", "FOO": "b"}
	out := TransformSnapshot(s, UpperKeys())
	if len(out) != 1 {
		t.Fatalf("expected collision to produce 1 key, got %d", len(out))
	}
	v := out["FOO"]
	if v != "a" && v != "b" {
		t.Fatalf("unexpected value after collision: %q", v)
	}
}
