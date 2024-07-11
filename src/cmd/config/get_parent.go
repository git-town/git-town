package config

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const getParentDesc = "Displays the parent branch for the current or given branch"

func getParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "get-parent [branch]",
		Args:  cobra.MaximumNArgs(1),
		Short: getParentDesc,
		Long:  cmdhelpers.Long(getParentDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGetParent(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeGetParent(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	var branch gitdomain.LocalBranchName
	if len(args) == 0 {
		branch, err = repo.Git.CurrentBranch(repo.Backend)
		if err != nil {
			return err
		}
	} else {
		branch = gitdomain.NewLocalBranchName(args[0])
	}
	parent := repo.UnvalidatedConfig.Config.Value.Lineage.Parent(branch)
	fmt.Println(parent)
	print.Footer(verbose, repo.CommandsCounter.Get(), repo.FinalMessages.Result())
	return nil
}
