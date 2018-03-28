package steps

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
)

// RunState represents the current state of a Git Town command,
// including which operations are left to do,
// and how to undo what has ben done so far.
type RunState struct {
	AbortStepList StepList
	Command       string
	IsAbort       bool
	isUndo        bool
	RunStepList   StepList
	UndoStepList  StepList
}

// LoadPreviousRunState loads the run state from disk if it exists or creates a new run state
func LoadPreviousRunState(command string) *RunState {
	filename := getRunResultFilename(command)
	if util.DoesFileExist(filename) {
		var runState RunState
		content, err := ioutil.ReadFile(filename)
		exit.If(err)
		err = json.Unmarshal(content, &runState)
		exit.If(err)
		return &runState
	}
	return &RunState{
		Command: command,
	}
}

// DeleteRunState deletes the run state from disk
func DeleteRunState(command string) {
	exit.If(os.Remove(getRunResultFilename(command)))
}

// AddPushBranchStepAfterCurrentBranchSteps inserts a PushBranchStep
// after all the steps for the current branch
func (runState *RunState) AddPushBranchStepAfterCurrentBranchSteps() {
	popped := StepList{}
	for {
		step := runState.RunStepList.Peek()
		if !isCheckoutBranchStep(step) {
			popped.Append(runState.RunStepList.Pop())
		} else {
			runState.RunStepList.Prepend(&PushBranchStep{BranchName: git.GetCurrentBranchName()})
			runState.RunStepList.PrependList(popped)
			break
		}
	}
}

// CreateAbortRunState returns a new runstate
// to be run to aborting and undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateAbortRunState() (result RunState) {
	result.Command = runState.Command
	result.IsAbort = true
	result.RunStepList.AppendList(runState.AbortStepList)
	result.RunStepList.AppendList(runState.UndoStepList)
	return
}

// CreateSkipRunState returns a new Runstate
// that skips operations for the current branch.
func (runState *RunState) CreateSkipRunState() (result RunState) {
	result.Command = runState.Command
	result.RunStepList.AppendList(runState.AbortStepList)
	for _, step := range runState.UndoStepList.List {
		if isCheckoutBranchStep(step) {
			break
		}
		result.RunStepList.Append(step)
	}
	skipping := true
	for _, step := range runState.RunStepList.List {
		if isCheckoutBranchStep(step) {
			skipping = false
		}
		if !skipping {
			result.RunStepList.Append(step)
		}
	}
	return
}

// CreateUndoRunState returns a new runstate
// to be run when undoing the Git Town command
// represented by this runstate.
func (runState *RunState) CreateUndoRunState() (result RunState) {
	result.Command = runState.Command
	result.isUndo = true
	result.RunStepList.AppendList(runState.UndoStepList)
	return
}

// Save saves the run state to disk
func (runState *RunState) Save() {
	content, err := json.MarshalIndent(runState, "", "  ")
	exit.If(err)
	filename := getRunResultFilename(runState.Command)
	err = ioutil.WriteFile(filename, content, 0644)
	exit.If(err)
}

// SkipCurrentBranchSteps removes the steps for the current branch
// from this run state.
func (runState *RunState) SkipCurrentBranchSteps() {
	for {
		step := runState.RunStepList.Peek()
		if !isCheckoutBranchStep(step) {
			runState.RunStepList.Pop()
		} else {
			break
		}
	}
}

func isCheckoutBranchStep(step Step) bool {
	return getTypeName(step) == "*CheckoutBranchStep"
}
