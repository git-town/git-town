package steps

type UndoLastCommitStep struct {
	EmptyStep
}

func (step *UndoLastCommitStep) Run(args RunArgs) error {
	return args.Runner.Frontend.UndoLastCommit()
}
