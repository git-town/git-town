package steps

// PullCurrentBranchStep updates the branch with the given name with commits from its remote.
type PullCurrentBranchStep struct {
	EmptyStep
}

func (step *PullCurrentBranchStep) Run(args RunArgs) error {
	return args.Runner.Frontend.Pull()
}
