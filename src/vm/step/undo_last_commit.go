package step

type UndoLastCommit struct {
	Empty
}

func (step *UndoLastCommit) Run(args RunArgs) error {
	return args.Runner.Frontend.UndoLastCommit()
}
