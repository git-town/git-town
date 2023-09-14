package steps

// AbortRebaseStep represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebaseStep struct {
	EmptyStep
}

func (step *AbortRebaseStep) Run(args RunArgs) error {
	return args.Run.Frontend.AbortRebase()
}
