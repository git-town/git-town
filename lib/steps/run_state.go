package steps

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has ben done so far.
type RunState struct {
	AbortStep    Step
	Command      string
	RunStepList  StepList
	UndoStepList StepList
}

// CreateAbortRunState returns a new runstate
// to be run to aborting and undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateAbortRunState() (result RunState) {
	result.Command = runState.Command
	result.RunStepList.Append(runState.AbortStep)
	result.RunStepList.AppendList(runState.UndoStepList)
	return
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (runState *RunState) CreateSkipRunState() (result RunState) {
	result.Command = runState.Command
	result.RunStepList.Append(runState.AbortStep)
	for _, step := range runState.UndoStepList.List {
		if getTypeName(step) == "CheckoutBranchStep" {
			break
		}
		result.RunStepList.Append(step)
	}
	skipping := true
	for _, step := range runState.RunStepList.List {
		if getTypeName(step) == "CheckoutBranchStep" {
			skipping = false
		}
		if !skipping {
			result.RunStepList.Append(step)
		}
	}
	return
}

// CreateUndoRunState returns a new runstate
// to be run when undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateUndoRunState() (result RunState) {
	result.Command = runState.Command
	result.RunStepList.AppendList(runState.UndoStepList)
	return
}

// SkipCurrentBranchSteps removes the steps for the current branch
// from this run state.
func (runState *RunState) SkipCurrentBranchSteps() {
	for {
		step := runState.RunStepList.Peek()
		if getTypeName(step) != "CheckoutBranchStep" {
			runState.RunStepList.Pop()
		} else {
			break
		}
	}
}
