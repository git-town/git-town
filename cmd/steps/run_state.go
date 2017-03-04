package steps

type RunState struct {
  AbortStep Step
  Command string
  RunStepList StepList
  SkipMessage string
  UndoStepList StepList
}

func (runState RunState) CreateAbortRunState() (result RunState) {
  result.Command = runState.Command
  result.RunStepList.Append(runState.AbortStep)
  result.RunStepList.AppendList(runState.UndoStepList)
  result.SkipMessage = runState.SkipMessage
  return result
}
