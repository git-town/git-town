package dialog

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

const (
	confirmUndoTitle  = `confirm undo action`
	undoForcePushHelp = `
Undo changes to remote branch %q?

Existing commit SHA: %s
Commit to be pushed: %s

`
)

type BoolEntry bool

func (self BoolEntry) Bool() bool {
	return bool(self)
}

func (self BoolEntry) String() string {
	return strconv.FormatBool(self.Bool())
}

// GitHubToken lets the user enter the GitHub API token.
func UndoForcePush(branch gitdomain.RemoteBranchName, existingSHA, shaToPush gitdomain.SHA, input components.TestInput) (result bool, aborted bool, err error) {
	fmt.Println("11111111111111111111111111111", input)
	entries := list.Entries[BoolEntry]{
		{
			Data:    true,
			Enabled: true,
			Text:    "Yes",
		},
		{
			Data:    false,
			Enabled: true,
			Text:    "No",
		},
	}
	helpText := fmt.Sprintf(undoForcePushHelp, branch, existingSHA, shaToPush)
	selection, aborted, err := components.RadioList(entries, 0, confirmUndoTitle, helpText, input)
	fmt.Printf("undo force-push %q: %s\n", branch, components.FormattedSelection(selection.String(), aborted))
	return selection.Bool(), aborted, err
}
