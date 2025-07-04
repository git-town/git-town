package debug

import (
	"os"
	"strconv"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/pkg/asserts"
	"github.com/spf13/cobra"
)

func enterCommitsToBeam() *cobra.Command {
	return &cobra.Command{
		Use:  "commits-to-beam <number of commits>",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			amount := asserts.NoError1(strconv.ParseInt(args[0], 10, 64))
			repo := asserts.NoError1(execute.OpenRepo(execute.OpenRepoArgs{
				DryRun:           false,
				PrintBranchNames: true,
				PrintCommands:    true,
				ValidateGitRepo:  true,
				ValidateIsOnline: false,
				Verbose:          false,
			}))
			allCommits := asserts.NoError1(repo.Git.CommitsInPerennialBranch(repo.Backend))
			commits := make([]gitdomain.Commit, amount)
			for i := range amount {
				commits[i] = allCommits[i]
			}
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err := dialog.CommitsToBeam(commits, "target-branch", repo.Git, repo.Backend, dialogTestInputs.Next())
			return err
		},
	}
}
