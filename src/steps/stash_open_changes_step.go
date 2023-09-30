package steps

type StashOpenChangesStep struct {
	EmptyStep
}

func (step *StashOpenChangesStep) Run(args RunArgs) error {
	return args.Runner.Frontend.Stash()
}
