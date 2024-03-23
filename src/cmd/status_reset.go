package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/cli/flags"
	"github.com/git-town/git-town/v13/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v13/src/execute"
	"github.com/git-town/git-town/v13/src/messages"
	"github.com/git-town/git-town/v13/src/vm/statefile"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

// TODO: extract this and the "status" command into a "status" subfolder,
//
//	similar to the "config" subfolder.
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

func executeStatusReset(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  false,
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
