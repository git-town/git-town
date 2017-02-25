package step

type Step interface {
  CreateAbortStep() Step
  CreateContinueStep() Step
  CreateUndoStep() Step
  Run() error
}
