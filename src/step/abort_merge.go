package step

// AbortMerge aborts the current merge conflict.
type AbortMerge struct {
	Empty
}

func (step *AbortMerge) Run(args RunArgs) error {
	return args.Runner.Frontend.AbortMerge()
}
