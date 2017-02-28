package step

import (
  "github.com/Originate/gt/cmd/script"
)


type FetchUpstreamStep struct {}


func (step FetchUpstreamStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step FetchUpstreamStep) CreateContinueStep() Step {
  return NoOpStep{}
}


func (step FetchUpstreamStep) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step FetchUpstreamStep) Run() error {
  return script.RunCommand([]string{"git", "fetch", "upstream"})
}
