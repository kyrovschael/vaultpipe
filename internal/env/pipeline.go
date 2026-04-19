package env

// StepFn is a function that transforms a Snapshot, returning the result or an
// error that halts the pipeline.
type StepFn func(Snapshot) (Snapshot, error)

// Pipeline executes a sequence of StepFns against an initial Snapshot,
// threading the output of each step into the next.
type Pipeline struct {
	steps []StepFn
}

// NewPipeline creates an empty Pipeline.
func NewPipeline(steps ...StepFn) *Pipeline {
	p := &Pipeline{}
	p.steps = append(p.steps, steps...)
	return p
}

// Add appends one or more steps to the pipeline.
func (p *Pipeline) Add(steps ...StepFn) *Pipeline {
	p.steps = append(p.steps, steps...)
	return p
}

// Run executes all steps in order starting from src.
func (p *Pipeline) Run(src Snapshot) (Snapshot, error) {
	current := src
	for _, step := range p.steps {
		var err error
		current, err = step(current)
		if err != nil {
			return Snapshot{}, err
		}
	}
	return current, nil
}

// Lift wraps a pure Snapshot→Snapshot function into a StepFn.
func Lift(fn func(Snapshot) Snapshot) StepFn {
	return func(s Snapshot) (Snapshot, error) {
		return fn(s), nil
	}
}
