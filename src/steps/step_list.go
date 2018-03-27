package steps

import (
	"encoding/json"
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
)

// StepList is a list of steps
// with convenience functions for adding and removing steps.
type StepList struct {
	List []Step
}

// Append adds the given step to the end of this StepList.
func (stepList *StepList) Append(step Step) {
	stepList.List = append(stepList.List, step)
}

// AppendList adds all elements of the given StepList to the end of this StepList.
func (stepList *StepList) AppendList(otherList StepList) {
	stepList.List = append(stepList.List, otherList.List...)
}

// IsEmpty returns whether or not this StepList has any elements
func (stepList *StepList) isEmpty() bool {
	return len(stepList.List) == 0
}

// Peek returns the first element of this StepList.
func (stepList *StepList) Peek() (result Step) {
	if stepList.isEmpty() {
		return nil
	}
	return stepList.List[0]
}

// Pop removes and returns the first element of this StepList.
func (stepList *StepList) Pop() (result Step) {
	if stepList.isEmpty() {
		return nil
	}
	result = stepList.List[0]
	stepList.List = stepList.List[1:]
	return
}

// Prepend adds the given step to the beginning of this StepList.
func (stepList *StepList) Prepend(step Step) {
	stepList.List = append([]Step{step}, stepList.List...)
}

// PrependList adds all elements of the given StepList to the start of this StepList.
func (stepList *StepList) PrependList(otherList StepList) {
	stepList.List = append(otherList.List, stepList.List...)
}

// WrapOptions represents the options given to Wrap.
type WrapOptions struct {
	RunInGitRoot     bool
	StashOpenChanges bool
}

// Wrap wraps the list with steps that
// change to the Git root directory or stash away open changes.
func (stepList *StepList) Wrap(options WrapOptions) {
	stepList.Append(&PreserveCheckoutHistoryStep{
		InitialBranch:                     git.GetCurrentBranchName(),
		InitialPreviouslyCheckedOutBranch: git.GetPreviouslyCheckedOutBranch(),
	})

	if options.StashOpenChanges && git.HasOpenChanges() {
		stepList.Prepend(&StashOpenChangesStep{})
		stepList.Append(&RestoreOpenChangesStep{})
	}

	initialDirectory, err := os.Getwd()
	exit.If(err)
	gitRootDirectory := git.GetRootDirectory()

	if options.RunInGitRoot && initialDirectory != gitRootDirectory {
		stepList.Prepend(&ChangeDirectoryStep{Directory: gitRootDirectory})
		stepList.Append(&ChangeDirectoryStep{Directory: initialDirectory})
	}
}

// MarshalJSON marshals the step list to JSON
func (stepList *StepList) MarshalJSON() (b []byte, e error) {
	jsonSteps := make([]*JSONStep, len(stepList.List))
	for i, step := range stepList.List {
		jsonSteps[i] = &JSONStep{Step: step}
	}
	return json.Marshal(jsonSteps)
}

// UnmarshalJSON unmarshals the step list from JSON
func (stepList *StepList) UnmarshalJSON(b []byte) error {
	var jsonSteps []JSONStep
	err := json.Unmarshal(b, &jsonSteps)
	if err != nil {
		return err
	}
	if len(jsonSteps) > 0 {
		stepList.List = make([]Step, len(jsonSteps))
		for i, jsonStep := range jsonSteps {
			stepList.List[i] = jsonStep.Step
		}
	}
	return nil
}
