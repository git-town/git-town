package steps

// AbortMergeStep aborts the current merge conflict.
type AbortMergeStep struct {
	EmptyStep
}

func (step *AbortMergeStep) Run(args RunArgs) error {
	return args.Run.Frontend.AbortMerge()
}
