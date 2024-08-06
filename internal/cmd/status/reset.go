package status

import (
	"fmt"

	"github.com/git-town/git-town/v14/internal/cli/flags"
	"github.com/git-town/git-town/v14/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/execute"
	"github.com/git-town/git-town/v14/internal/messages"
	"github.com/git-town/git-town/v14/internal/vm/statefile"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

func resetRunstateCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetDesc,
		Long:  cmdhelpers.Long(statusResetDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeStatusReset(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeStatusReset(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	err = statefile.Delete(repo.RootDir)
	if err != nil {
		return err
	}
	fmt.Println(messages.RunstateDeleted)
	return nil
}
