package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type ResetToShaStep struct {
  sha string
  hard bool
}

func (step ResetToShaStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step ResetToShaStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step ResetToShaStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step ResetToShaStep) Run() error {
  if step.sha == git.GetCurrentSha() {
    cmd := []string{"git", "reset"}
    if step.hard {
      cmd = append(cmd, "--hard")
    }
    cmd = append(cmd, step.sha)
    return script.RunCommand(cmd)
  }
  return nil
}
