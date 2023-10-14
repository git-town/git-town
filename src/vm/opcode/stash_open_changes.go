package opcode

type StashOpenChanges struct {
	Empty
}

func (step *StashOpenChanges) Run(args RunArgs) error {
	return args.Runner.Frontend.Stash()
}
