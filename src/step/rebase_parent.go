package step

// RebaseParent rebases the current branch against its current parent branch.
type RebaseParent struct {
	Empty
}

func (step *RebaseParent) CreateAbortSteps() []Step {
	return []Step{
		&AbortRebase{},
	}
}

func (step *RebaseParent) CreateContinueSteps() []Step {
	return []Step{
		&ContinueRebase{},
	}
}

func (step *RebaseParent) Run(args RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	parent := args.Lineage.Parent(currentBranch)
	return args.Runner.Frontend.Rebase(parent.BranchName())
}
