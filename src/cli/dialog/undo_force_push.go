package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

const (
	undoForcePushTitle = `force-push to remote branch`
	undoForcePushHelp  = `
Should I force-push remote branch %q?

Existing commit: %q
Commit to be pushed: %q

`
)

type YesNoEntry struct {
	Text  string
	Value bool
}

func (self YesNoEntry) String() string {
	return fmt.Sprintf("%t", self.Value)
}

// GitHubToken lets the user enter the GitHub API token.
func ForcePushBranch(branch gitdomain.RemoteBranchName, inputs components.TestInput) (result bool, aborted bool, err error) {
	entries := list.Entries[YesNoEntry]{}
	entry, aborted, err := components.RadioList(entries, 0, undoForcePushTitle, undoForcePushHelp, inputs)
	return entry.Value, aborted, err
}
