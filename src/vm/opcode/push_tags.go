package opcode

// PushTags pushes newly created Git tags to origin.
type PushTags struct {
	Empty
}

func (step *PushTags) Run(args RunArgs) error {
	return args.Runner.Frontend.PushTags()
}
