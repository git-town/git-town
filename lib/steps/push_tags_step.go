package steps

import (
  "github.com/Originate/gt/lib/script"
)


type PushTagsStep struct {}


func (step PushTagsStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step PushTagsStep) CreateContinueStep() Step {
  return NoOpStep{}
}


func (step PushTagsStep) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step PushTagsStep) Run() error {
  return script.RunCommand([]string{"git", "push", "--tags"})
}
