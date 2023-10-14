package opcode

// PullCurrentBranch updates the branch with the given name with commits from its remote.
type PullCurrentBranch struct {
	Empty
}

func (step *PullCurrentBranch) Run(args RunArgs) error {
	return args.Runner.Frontend.Pull()
}
