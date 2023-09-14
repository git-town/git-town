package steps

// PushTagsStep pushes newly created Git tags to origin.
type PushTagsStep struct {
	EmptyStep
}

func (step *PushTagsStep) Run(args RunArgs) error {
	return args.Run.Frontend.PushTags()
}
