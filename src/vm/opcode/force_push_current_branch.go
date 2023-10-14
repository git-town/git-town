package opcode

// ForcePushCurrentBranch force-pushes the branch with the given name to the origin remote.
type ForcePushCurrentBranch struct {
	NoPushHook bool
	BaseOpcode
}

func (step *ForcePushCurrentBranch) Run(args RunArgs) error {
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	shouldPush, err := args.Runner.Backend.ShouldPushBranch(currentBranch, currentBranch.TrackingBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Runner.Config.DryRun {
		return nil
	}
	return args.Runner.Frontend.ForcePushBranch(step.NoPushHook)
}
