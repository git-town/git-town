package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

const (
	undoForcePushTitle = `confirm undo`
	undoForcePushHelp  = `
Undo changes to remote branch %q?

Existing commit: %q
Commit to be pushed: %q

`
)

type BoolEntry bool

func (self BoolEntry) String() string {
	return fmt.Sprintf("%t", self)
}

func (self BoolEntry) Bool() bool {
	return bool(self)
}

// GitHubToken lets the user enter the GitHub API token.
func ForcePushBranch(branch gitdomain.RemoteBranchName, existingSHA, SHAToPush gitdomain.SHA, inputs components.TestInput) (result bool, aborted bool, err error) {
	entries := list.Entries[BoolEntry]{
		{
			Data:    false,
			Enabled: true,
			Text:    "No",
		},
		{
			Data:    true,
			Enabled: true,
			Text:    "Yes",
		},
	}
	helpText := fmt.Sprintf(undoForcePushHelp, branch, existingSHA, SHAToPush)
	entry, aborted, err := components.RadioList(entries, 0, undoForcePushTitle, helpText, inputs)
	return entry.Bool(), aborted, err
}
