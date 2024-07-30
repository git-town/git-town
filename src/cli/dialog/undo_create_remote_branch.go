package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

const (
	undoCreateRemoteBranchTitle = `confirm undo`
	undoCreateRemoteBranchHelp  = `
Delete remote branch %q?

`
)

// GitHubToken lets the user enter the GitHub API token.
func UndoCreateRemoteBranch(branch gitdomain.RemoteBranchName, inputs components.TestInput) (result bool, aborted bool, err error) {
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
	helpText := fmt.Sprintf(undoCreateRemoteBranchHelp, branch)
	entry, aborted, err := components.RadioList(entries, 0, undoCreateRemoteBranchTitle, helpText, inputs)
	return entry.Bool(), aborted, err
}
