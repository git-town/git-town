package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const commitsToBeamTitle = `Select the commits to beam into branch %s`

// lets the user select commits to beam to the target branch
func CommitsToBeam(commits []gitdomain.Commit, targetBranch gitdomain.LocalBranchName, inputs components.TestInput) (gitdomain.Commits, bool, error) {
	if len(commits) == 0 {
		return gitdomain.Commits{}, false, nil
	}
	entries := make(list.Entries[gitdomain.Commit], len(commits))
	for c, commit := range commits {
		entries[c] = list.Entry[gitdomain.Commit]{
			Data: commit,
			Text: commit.Message.String(),
		}
	}
	selection, aborted, err := components.CheckList(entries, []int{}, commitsToBeamTitle, "", inputs)
	fmt.Printf(messages.CommitsSelected, len(selection))
	return selection, aborted, err
}
