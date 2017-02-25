package step

type NoOpStep int

func (step NoOpStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step NoOpStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step NoOpStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step NoOpStep) Run() error {
  return nil
}
