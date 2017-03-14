package steps

import (
  "github.com/Originate/git-town/lib/git"
  "github.com/Originate/git-town/lib/script"
)


type ResetToShaStep struct {
  Hard bool
  Sha string
}


func (step ResetToShaStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step ResetToShaStep) CreateContinueStep() Step {
  return NoOpStep{}
}


func (step ResetToShaStep) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step ResetToShaStep) Run() error {
  if step.Sha == git.GetCurrentSha() {
    return nil
  }
  cmd := []string{"git", "reset"}
  if step.Hard {
    cmd = append(cmd, "--hard")
  }
  cmd = append(cmd, step.Sha)
  return script.RunCommand(cmd)
}
