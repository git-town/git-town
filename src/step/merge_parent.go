package step

// MergeParent merges the current parent of the current branch into the current branch.
type MergeParent struct {
	Empty
}

func (step *MergeParent) CreateAbortSteps() []Step {
	return []Step{
		&AbortMerge{},
	}
}

func (step *MergeParent) CreateContinueSteps() []Step {
	return []Step{
		&ContinueMerge{},
	}
}

func (step *MergeParent) Run(args RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	parent := args.Lineage.Parent(currentBranch)
	if parent.IsEmpty() {
		return nil
	}
	return args.Runner.Frontend.MergeBranchNoEdit(parent.BranchName())
}
