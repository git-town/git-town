package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v19/internal/git"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/messages"
)

const commitsToBeamTitle = `Select the commits to beam into branch %q`

// lets the user select commits to beam to the target branch
func CommitsToBeam(commits []gitdomain.Commit, targetBranch gitdomain.LocalBranchName, git git.Commands, querier gitdomain.Querier, inputs components.TestInput) (gitdomain.Commits, bool, error) {
	if len(commits) == 0 {
		return gitdomain.Commits{}, false, nil
	}
	entries := make(list.Entries[gitdomain.Commit], len(commits))
	for c, commit := range commits {
		shortSHA, err := git.ShortenSHA(querier, commit.SHA)
		if err != nil {
			return gitdomain.Commits{}, false, err
		}
		entries[c] = list.Entry[gitdomain.Commit]{
			Data: commit,
			Text: fmt.Sprintf("%s (%s)", commit.Message.String(), shortSHA),
		}
	}
	selection, aborted, err := components.CheckList(entries, []int{}, fmt.Sprintf(commitsToBeamTitle, targetBranch), "", inputs)
	fmt.Printf(messages.CommitsSelected, len(selection))
	return selection, aborted, err
}
