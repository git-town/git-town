package steps

// RunOptions bundles the parameters for running a Git Town command.
type RunOptions struct {
	CanSkip              func() bool
	Command              string
	IsAbort              bool
	IsContinue           bool
	IsSkip               bool
	IsUndo               bool
	SkipMessageGenerator func() string
	StepListGenerator    func() StepList
}

// ShouldLoadStateFromDisk returns whether or not the state should be load inserts a PushBranchStep
// after all the steps for the current branch
func (r *RunOptions) ShouldLoadStateFromDisk() bool {
	return r.IsAbort || r.IsContinue || r.IsSkip || r.IsUndo
}
