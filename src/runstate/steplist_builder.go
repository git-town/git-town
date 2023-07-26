package runstate

import (
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/steps"
)

// StepListBuilder allows populating StepList instances
// in concise ways out of fallible operations.
//
// This is based on ideas outlined in https://go.dev/blog/errors-are-values.
type StepListBuilder struct {
	StepList          StepList `exhaustruct:"optional"`
	failure.Collector `exhaustruct:"optional"`
}

func (slb *StepListBuilder) Add(step steps.Step) {
	slb.StepList.Append(step)
}

func (slb *StepListBuilder) AddE(step steps.Step, err error) {
	if !slb.Check(err) {
		slb.Add(step)
	}
}

func (slb *StepListBuilder) Wrap(options WrapOptions, backend *git.BackendCommands) {
	slb.Check(slb.StepList.Wrap(options, backend))
}

func (slb *StepListBuilder) Result() (StepList, error) {
	return slb.StepList, slb.Err
}
