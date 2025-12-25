package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
)

const commitsToBeamTitle = `Select the commits to beam into branch %q`

// CommitsToBeam lets the user select commits to beam to the target branch.
func CommitsToBeam(commits []gitdomain.Commit, targetBranch gitdomain.LocalBranchName, git git.Commands, querier subshelldomain.Querier, inputs dialogcomponents.Inputs) (gitdomain.Commits, dialogdomain.Exit, error) {
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
			Text: fmt.Sprintf("%s %s", shortSHA, commit.Message.String()),
		}
	}
	selection, exit, err := dialogcomponents.CheckList(entries, []int{}, fmt.Sprintf(commitsToBeamTitle, targetBranch), "", inputs, "commits-to-beam")
	fmt.Printf(messages.CommitsSelected, len(selection))
	return selection, exit, err
}
