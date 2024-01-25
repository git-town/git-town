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
	ResponseContinue = Response("continue") // continue the unfinished run state
	ResponseDiscard  = Response("discard")  // discard the unfinished run state
	ResponseQuit     = Response("quit")     // quit the program
	ResponseSkip     = Response("skip")     // continue the unfinished run state by skipping the current branch
	ResponseUndo     = Response("undo")     // undo the unfinished run state
)

type unfinishedRunstateDialogEntry struct {
	response Response
	text     string
}

func (self unfinishedRunstateDialogEntry) String() string {
	return self.text
}

// AskHowToHandleUnfinishedRunState prompts the user for how to handle the unfinished run state.
func AskHowToHandleUnfinishedRunState(command string, endBranch gitdomain.LocalBranchName, endTime time.Time, canSkip bool, dialogTestInput TestInput) (Response, bool, error) {
	options := []unfinishedRunstateDialogEntry{
		{response: ResponseQuit, text: messages.UnfinishedRunStateQuit},
		{response: ResponseContinue, text: fmt.Sprintf(messages.UnfinishedRunStateContinue, command)},
	}
	if canSkip {
		options = append(options, unfinishedRunstateDialogEntry{response: ResponseSkip, text: fmt.Sprintf(messages.UnfinishedRunStateSkip, command)})
	}
	options = append(options,
		unfinishedRunstateDialogEntry{response: ResponseUndo, text: fmt.Sprintf(messages.UnfinishedRunStateUndo, command)},
		unfinishedRunstateDialogEntry{response: ResponseDiscard, text: messages.UnfinishedRunStateDiscard},
	)
	selection, aborted, err := radioList(options, 0, fmt.Sprintf(unfinishedRunstateHelp, command, endBranch, humanize.Time(endTime)), dialogTestInput)
	fmt.Printf("Handle unfinished command: %s\n", formattedSelection(selection.response.String(), aborted))
	return selection.response, aborted, err
}
