package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type ResetToShaStep struct {
  Hard bool
  Sha string
}

func (step ResetToShaStep) CreateAbortStep() Step {
  return nil
}

func (step ResetToShaStep) CreateContinueStep() Step {
  return nil
}

func (step ResetToShaStep) CreateUndoStep() Step {
  return nil
}

func (step ResetToShaStep) Run() error {
  if step.Sha == git.GetCurrentSha() {
    cmd := []string{"git", "reset"}
    if step.Hard {
      cmd = append(cmd, "--hard")
    }
    cmd = append(cmd, step.Sha)
    return script.RunCommand(cmd)
  }
  return nil
}
