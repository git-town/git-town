package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli/flags"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/vm/persistence"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

func resetRunstateCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetDesc,
		Long:  long(statusResetDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeStatusReset(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func executeStatusReset(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	err = persistence.Delete(repo.RootDir)
	if err != nil {
		return err
	}
	fmt.Println("Runstate file deleted.")
	return nil
}
