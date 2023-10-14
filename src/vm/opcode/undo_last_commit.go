package opcode

type UndoLastCommit struct {
	undeclaredOpcodeMethods
}

func (step *UndoLastCommit) Run(args RunArgs) error {
	return args.Runner.Frontend.UndoLastCommit()
}
