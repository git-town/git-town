package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/spf13/cobra"
)

func enterUnknownBranch() *cobra.Command {
	return &cobra.Command{
		Use: "unknown-branch-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			repo, err := execute.OpenRepo(execute.OpenRepoArgs{
				CliConfig: cliconfig.CliConfig{
					DryRun:  false,
					Verbose: false,
				},
				PrintBranchNames: false,
				PrintCommands:    true,
				ValidateGitRepo:  true,
				ValidateIsOnline: false,
			})
			if err != nil {
				return err
			}
			dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
			_, _, err = dialog.UnknownBranchType(repo.UnvalidatedConfig, dialogTestInputs)
			return err
		},
	}
}
