package runstate

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/steps"
)

// StepListBuilder allows populating StepList instances
// in concise ways out of fallible operations.
//
// This is based on ideas outlined in https://go.dev/blog/errors-are-values.
type StepListBuilder struct {
	StepList     StepList `exhaustruct:"optional"`
	ErrorChecker `exhaustruct:"optional"`
}

func (slb *StepListBuilder) Add(step steps.Step) {
	slb.StepList.Append(step)
}

func (slb *StepListBuilder) AddE(step steps.Step, err error) {
	if !slb.Check(err) {
		slb.Add(step)
	}
}

func (slb *StepListBuilder) Wrap(options WrapOptions, repo *git.BackendCommands, mainBranch string) {
	slb.Check(slb.StepList.Wrap(options, repo, mainBranch))
}

func (slb *StepListBuilder) Result() (StepList, error) {
	return slb.StepList, slb.Err
}
