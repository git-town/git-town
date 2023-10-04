package steps

import (
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/step"
)

// StepListBuilder allows populating StepList instances
// in concise ways out of fallible operations.
//
// This is based on ideas outlined in https://go.dev/blog/errors-are-values.
type StepListBuilder struct {
	StepList                 List `exhaustruct:"optional"`
	gohacks.FailureCollector `exhaustruct:"optional"`
}

func (slb *StepListBuilder) Add(step step.Step) {
	slb.StepList.Append(step)
}

func (slb *StepListBuilder) AddE(step step.Step, err error) {
	if !slb.Check(err) {
		slb.Add(step)
	}
}

func (slb *StepListBuilder) Wrap(options WrapOptions) {
	slb.Check(slb.StepList.Wrap(options))
}

func (slb *StepListBuilder) Result() (List, error) {
	return slb.StepList, slb.Err
}
