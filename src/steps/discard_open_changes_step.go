package steps

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	EmptyStep
}

func (step *DiscardOpenChangesStep) Run(args RunArgs) error {
	return args.Runner.Frontend.DiscardOpenChanges()
}
