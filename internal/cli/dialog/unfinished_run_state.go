package dialog

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
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
	ResponseBoth     = Response("both")     // continue the old runstate and run the new program
)

// AskHowToHandleUnfinishedRunState prompts the user for how to handle the unfinished run state.
func AskHowToHandleUnfinishedRunState(command string, endBranch gitdomain.LocalBranchName, endTime time.Time, canSkip bool, input dialogcomponents.Inputs) (Response, dialogdomain.Exit, error) {
	entries := list.Entries[Response]{
		{
			Data: ResponseQuit,
			Text: messages.UnfinishedRunStateQuit,
		},
		{
			Data: ResponseContinue,
			Text: fmt.Sprintf(messages.UnfinishedRunStateContinue, command),
		},
	}
	if canSkip {
		entries = append(entries, list.Entry[Response]{
			Data: ResponseSkip,
			Text: fmt.Sprintf(messages.UnfinishedRunStateSkip, command),
		})
	}
	entries = append(entries,
		list.Entry[Response]{
			Data: ResponseUndo,
			Text: fmt.Sprintf(messages.UnfinishedRunStateUndo, command),
		},
		list.Entry[Response]{
			Data: ResponseDiscard,
			Text: messages.UnfinishedRunStateDiscard,
		},
		list.Entry[Response]{
			Data: ResponseBoth,
			Text: fmt.Sprintf(messages.UnfinishedRunStateBoth, command),
		},
	)
	selection, exit, err := dialogcomponents.RadioList(entries, 0, unfinishedRunstateTitle, fmt.Sprintf(unfinishedRunstateHelp, command, endBranch, humanize.Time(endTime)), input, "unfinished-runstate")
	fmt.Printf(messages.UnfinishedCommandHandle, dialogcomponents.FormattedSelection(string(selection), exit))
	return selection, exit, err
}
