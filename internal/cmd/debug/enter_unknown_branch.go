package debug

import (
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

func enterUnknownBranch() *cobra.Command {
	return &cobra.Command{
		Use: "unknown-branch-type",
		RunE: func(_ *cobra.Command, _ []string) error {
			repo, err := execute.OpenRepo(execute.OpenRepoArgs{
				CliConfig: cliconfig.New(cliconfig.NewArgs{
					DryRun:  None[configdomain.DryRun](),
					Verbose: None[configdomain.Verbose](),
				}),
				PrintBranchNames: false,
				PrintCommands:    true,
				ValidateGitRepo:  true,
				ValidateIsOnline: false,
			})
			if err != nil {
				return err
			}
			inputs := dialogcomponents.LoadInputs(os.Environ())
			_, _, err = dialog.UnknownBranchType(dialog.Args[configdomain.UnknownBranchType]{
				Global: repo.UnvalidatedConfig.GitGlobal.UnknownBranchType,
				Inputs: inputs,
				Local:  repo.UnvalidatedConfig.GitLocal.UnknownBranchType,
			})
			return err
		},
	}
}
