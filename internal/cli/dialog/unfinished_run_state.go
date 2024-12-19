package dialog

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	unfinishedRunstateTitle = `unfinished Git Town command`
	unfinishedRunstateHelp  = `
You have an unfinished %q command
that ended on the %q branch
%s. Please choose how to proceed.


`
)

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
func AskHowToHandleUnfinishedRunState(command string, endBranch gitdomain.LocalBranchName, endTime time.Time, canSkip bool, dialogTestInput components.TestInput) (Response, bool, error) {
	entries := list.Entries[Response]{
		{
			Data:    ResponseQuit,
			Enabled: true,
			Text:    messages.UnfinishedRunStateQuit,
		},
		{
			Data:    ResponseContinue,
			Enabled: true,
			Text:    fmt.Sprintf(messages.UnfinishedRunStateContinue, command),
		},
	}
	if canSkip {
		entries = append(entries, list.Entry[Response]{
			Data:    ResponseSkip,
			Enabled: true,
			Text:    fmt.Sprintf(messages.UnfinishedRunStateSkip, command),
		})
	}
	entries = append(entries,
		list.Entry[Response]{
			Data:    ResponseUndo,
			Enabled: true,
			Text:    fmt.Sprintf(messages.UnfinishedRunStateUndo, command),
		},
		list.Entry[Response]{
			Data:    ResponseDiscard,
			Enabled: true,
			Text:    messages.UnfinishedRunStateDiscard,
		},
	)
	selection, aborted, err := components.RadioList(entries, 0, unfinishedRunstateTitle, fmt.Sprintf(unfinishedRunstateHelp, command, endBranch, humanize.Time(endTime)), dialogTestInput)
	fmt.Printf(messages.UnfinishedCommandHandle, components.FormattedSelection(string(selection), aborted))
	return selection, aborted, err
}
