package env

import (
	"errors"
	"testing"
)

func TestPipeline_Empty(t *testing.T) {
	src := FromSlice([]string{"A=1"})
	out, err := NewPipeline().Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Get("A") != "1" {
		t.Fatalf("expected A=1, got %q", out.Get("A"))
	}
}

func TestPipeline_StepsAppliedInOrder(t *testing.T) {
	src := FromSlice([]string{"a=hello", "b=world"})

	p := NewPipeline(
		Lift(UpperKeys),
		Lift(func(s Snapshot) Snapshot { return TrimValues(s) }),
	)
	out, err := p.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Get("A") != "hello" {
		t.Fatalf("expected A=hello, got %q", out.Get("A"))
	}
	if out.Get("B") != "world" {
		t.Fatalf("expected B=world, got %q", out.Get("B"))
	}
}

func TestPipeline_HaltsOnError(t *testing.T) {
	sentinel := errors.New("step error")
	called := false
	p := NewPipeline(
		func(s Snapshot) (Snapshot, error) { return Snapshot{}, sentinel },
		func(s Snapshot) (Snapshot, error) { called = true; return s, nil },
	)
	_, err := p.Run(FromSlice(nil))
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if called {
		t.Fatal("second step should not have been called")
	}
}

func TestPipeline_Add_Chaining(t *testing.T) {
	p := NewPipeline()
	p.Add(Lift(UpperKeys)).Add(Lift(func(s Snapshot) Snapshot { return s }))
	src := FromSlice([]string{"x=1"})
	out, err := p.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Get("X") != "1" {
		t.Fatalf("expected X=1")
	}
}

func TestLift_WrapsFunction(t *testing.T) {
	step := Lift(UpperKeys)
	src := FromSlice([]string{"foo=bar"})
	out, err := step(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Get("FOO") != "bar" {
		t.Fatalf("expected FOO=bar")
	}
}
