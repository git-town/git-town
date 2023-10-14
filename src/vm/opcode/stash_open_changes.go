package opcode

type StashOpenChanges struct {
	undeclaredOpcodeMethods
}

func (step *StashOpenChanges) Run(args RunArgs) error {
	return args.Runner.Frontend.Stash()
}
