package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

const (
	undoCreateRemoteBranchHelp = `
Delete remote branch %q?

`
)

// GitHubToken lets the user enter the GitHub API token.
func UndoCreateRemoteBranch(branch gitdomain.RemoteBranchName, input components.TestInput) (result bool, aborted bool, err error) {
	fmt.Println("2222222222222222222222222222222222", input)
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
	selection, aborted, err := components.RadioList(entries, 0, confirmUndoTitle, helpText, input)
	fmt.Printf("delete branch %q: %s\n", branch, components.FormattedSelection(selection.String(), aborted))
	return selection.Bool(), aborted, err
}
