package dialog

import (
	"fmt"

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

func (self BoolEntry) String() string {
	return fmt.Sprintf("%v", self.Bool())
}

func (self BoolEntry) Bool() bool {
	return bool(self)
}

// GitHubToken lets the user enter the GitHub API token.
func UndoForcePush(branch gitdomain.RemoteBranchName, existingSHA, SHAToPush gitdomain.SHA, inputs components.TestInput) (result bool, aborted bool, err error) {
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
	helpText := fmt.Sprintf(undoForcePushHelp, branch, existingSHA, SHAToPush)
	selection, aborted, err := components.RadioList(entries, 0, confirmUndoTitle, helpText, inputs)
	fmt.Printf("undo force-push %q: %s\n", branch, components.FormattedSelection(selection.String(), aborted))
	return selection.Bool(), aborted, err
}
