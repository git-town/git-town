package step

// AbortRebase represents aborting on ongoing merge conflict.
// This step is used in the abort scripts for Git Town commands.
type AbortRebase struct {
	Empty
}

func (step *AbortRebase) Run(args RunArgs) error {
	return args.Runner.Frontend.AbortRebase()
}
