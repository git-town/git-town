package opcode

// DiscardOpenChanges resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChanges struct {
	undeclaredOpcodeMethods
}

func (step *DiscardOpenChanges) Run(args RunArgs) error {
	return args.Runner.Frontend.DiscardOpenChanges()
}
