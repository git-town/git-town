package steps

import (
	"github.com/git-town/git-town/v9/src/gohacks"
	"github.com/git-town/git-town/v9/src/step"
)

// ListBuilder allows populating StepList instances
// in concise ways out of fallible operations.
//
// This is based on ideas outlined in https://go.dev/blog/errors-are-values.
type ListBuilder struct {
	StepList                 List `exhaustruct:"optional"`
	gohacks.FailureCollector `exhaustruct:"optional"`
}

func (slb *ListBuilder) Add(step step.Step) {
	slb.StepList.Append(step)
}

func (slb *ListBuilder) AddE(step step.Step, err error) {
	if !slb.Check(err) {
		slb.Add(step)
	}
}

func (slb *ListBuilder) Wrap(options WrapOptions) {
	slb.Check(slb.StepList.Wrap(options))
}

func (slb *ListBuilder) Result() (List, error) {
	return slb.StepList, slb.Err
}
