package opcode

// PushTags pushes newly created Git tags to origin.
type PushTags struct {
	undeclaredOpcodeMethods
}

func (step *PushTags) Run(args RunArgs) error {
	return args.Runner.Frontend.PushTags()
}
