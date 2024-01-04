package perennialbranches

import (
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/spf13/cobra"
)

const displaySummary = "Displays the perennial branches"

const displayHelp = `
Perennial branches are long-lived branches.
They cannot be shipped.`

func RootCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "perennial-branches",
		Args:  cobra.NoArgs,
		Short: displaySummary,
		Long:  cmdhelpers.Long(displaySummary, displayHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDisplay(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	cmd.AddCommand(addCmd())
	cmd.AddCommand(removeCmd())
	cmd.AddCommand(updateCmd())
	return &cmd
}

func executeDisplay(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	io.Println(format.StringSetting(repo.Runner.Config.PerennialBranches.Join("\n")))
	return nil
}
