package opcode

type UndoLastCommit struct {
	BaseOpcode
}

func (step *UndoLastCommit) Run(args RunArgs) error {
	return args.Runner.Frontend.UndoLastCommit()
}
