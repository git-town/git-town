package opcode

// AbortMerge aborts the current merge conflict.
type AbortMerge struct {
	undeclaredOpcodeMethods
}

func (step *AbortMerge) Run(args RunArgs) error {
	return args.Runner.Frontend.AbortMerge()
}
