package steps

type SkipCurrentBranchSteps struct {}


func (step SkipCurrentBranchSteps) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step SkipCurrentBranchSteps) CreateContinueStep() Step {
  return NoOpStep{}
}


func (step SkipCurrentBranchSteps) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step SkipCurrentBranchSteps) Run() error {
  return nil
}
