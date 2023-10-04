package steps

import (
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/step"
)

// Builder allows populating List instances
// in concise ways out of fallible operations.
//
// This is based on ideas outlined in https://go.dev/blog/errors-are-values.
type Builder struct {
	StepList                 List `exhaustruct:"optional"`
	gohacks.FailureCollector `exhaustruct:"optional"`
}

func (b *Builder) Add(step step.Step) {
	b.StepList.Append(step)
}

func (b *Builder) AddE(step step.Step, err error) {
	if !b.Check(err) {
		b.Add(step)
	}
}

func (b *Builder) Wrap(options WrapOptions) {
	b.Check(b.StepList.Wrap(options))
}

func (b *Builder) Result() (List, error) {
	return b.StepList, b.Err
}
