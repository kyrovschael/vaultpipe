// Package env provides the Pipeline type for composing ordered Snapshot
// transformations.
//
// A Pipeline chains StepFn values — each receiving the output of the previous
// step — and halts immediately if any step returns an error.
//
// Example:
//
//	p := env.NewPipeline(
//		env.Lift(env.UpperKeys),
//		env.Lift(env.TrimValues),
//	)
//	out, err := p.Run(src)
package env
