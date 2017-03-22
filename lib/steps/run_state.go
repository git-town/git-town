package steps

type RunState struct {
	AbortStep    Step
	Command      string
	RunStepList  StepList
	UndoStepList StepList
}

func (runState *RunState) CreateAbortRunState() (result RunState) {
	result.Command = runState.Command
	result.RunStepList.Append(runState.AbortStep)
	result.RunStepList.AppendList(runState.UndoStepList)
	return
}

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

func (runState *RunState) CreateUndoRunState() (result RunState) {
	result.Command = runState.Command
	result.RunStepList.AppendList(runState.UndoStepList)
	return
}

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
