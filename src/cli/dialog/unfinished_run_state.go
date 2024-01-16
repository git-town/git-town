package dialog

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/messages"
)

const unfinishedRunstateHelp = `
You have an unfinished %q command
that ended on the %q branch
%s. Please choose how to proceed.


`

type Response string

func (self Response) String() string { return string(self) }

const (
	// ResponseContinue stands for the user choosing to continue the unfinished run state.
	ResponseContinue = Response("continue")
	// ResponseDiscard stands for the user choosing to discard the unfinished run state.
	ResponseDiscard = Response("discard")
	// ResponseQuit stands for the user choosing to quit the program.
	ResponseQuit = Response("quit")
	// ResponseSkip stands for the user choosing to continue the unfinished run state by skipping the current branch.
	ResponseSkip = Response("skip")
	// ResponseUndo stands for the user choosing to undo the unfinished run state.
	ResponseUndo = Response("undo")
)

// AskHowToHandleUnfinishedRunState prompts the user for how to handle the unfinished run state.
func AskHowToHandleUnfinishedRunState(command string, endBranch gitdomain.LocalBranchName, endTime time.Time, canSkip bool, dialogTestInput TestInput) (Response, bool, error) {
	formattedOptions := map[Response]string{
		ResponseContinue: fmt.Sprintf(messages.UnfinishedRunStateContinue, command),
		ResponseDiscard:  messages.UnfinishedRunStateDiscard,
		ResponseQuit:     messages.UnfinishedRunStateQuit,
		ResponseSkip:     fmt.Sprintf(messages.UnfinishedRunStateSkip, command),
		ResponseUndo:     fmt.Sprintf(messages.UnfinishedRunStateUndo, command),
	}
	options := []string{
		formattedOptions[ResponseQuit],
		formattedOptions[ResponseContinue],
	}
	if canSkip {
		options = append(options, formattedOptions[ResponseSkip])
	}
	options = append(options, formattedOptions[ResponseUndo], formattedOptions[ResponseDiscard])
	selection, aborted, err := radioList(radioListArgs{
		entries:      options,
		defaultEntry: "",
		help:         fmt.Sprintf(unfinishedRunstateHelp, command, endBranch, humanize.Time(endTime)),
		testInput:    dialogTestInput,
	})
	fmt.Printf("Handle unfinished command: %s\n", formattedSelection(selection, aborted))
	for responseType, formattedResponseType := range formattedOptions {
		if formattedResponseType == selection {
			return responseType, aborted, err
		}
	}
	return ResponseUndo, aborted, fmt.Errorf(messages.DialogUnexpectedResponse, selection)
}
