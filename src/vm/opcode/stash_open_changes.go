package opcode

type StashOpenChanges struct {
	BaseOpcode
}

func (step *StashOpenChanges) Run(args RunArgs) error {
	return args.Runner.Frontend.Stash()
}
