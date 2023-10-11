package step

import "github.com/git-town/git-town/v9/src/domain"

// Merge merges the branch with the given name into the current branch.
type Merge struct {
	Branch domain.BranchName
	Empty
}

func (step *Merge) CreateAbortSteps() []Step {
	return []Step{
		&AbortMerge{},
	}
}

func (step *Merge) CreateContinueSteps() []Step {
	return []Step{
		&ContinueMerge{},
	}
}

func (step *Merge) Run(args RunArgs) error {
	return args.Runner.Frontend.MergeBranchNoEdit(step.Branch)
}
