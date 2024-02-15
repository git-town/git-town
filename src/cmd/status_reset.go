package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/statefile"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeStatusReset(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeStatusReset(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
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
