package config

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.CliConfig{
				DryRun:  false,
				Verbose: verbose,
			}
			return executeGetParent(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeGetParent(args []string, cliConfig cliconfig.CliConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    false,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	var childBranch gitdomain.LocalBranchName
	if len(args) == 0 {
		childBranch, err = repo.Git.CurrentBranch(repo.Backend)
		if err != nil {
			return err
		}
	} else {
		childBranch = gitdomain.NewLocalBranchName(args[0])
	}
	parentOpt := repo.UnvalidatedConfig.NormalConfig.Lineage.Parent(childBranch)
	if parent, hasParent := parentOpt.Get(); hasParent {
		fmt.Print(parent)
	}
	print.Footer(cliConfig.Verbose, repo.CommandsCounter.Immutable(), repo.FinalMessages.Result())
	return nil
}
