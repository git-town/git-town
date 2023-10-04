package step

// DiscardOpenChanges resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChanges struct {
	Empty
}

func (step *DiscardOpenChanges) Run(args RunArgs) error {
	return args.Runner.Frontend.DiscardOpenChanges()
}
