package steps

// PullBranchStep updates the branch with the given name with commits from its remote.
// TODO: rename to PullCurrentBranchStep and remove the "Branch" field.
type PullBranchStep struct {
	Branch string
	EmptyStep
}

func (step *PullBranchStep) Run(args RunArgs) error {
	return args.Run.Frontend.Pull()
}
