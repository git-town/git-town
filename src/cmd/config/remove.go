package config

import (
	"slices"
	"strings"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const removeConfigDesc = "Removes the Git Town configuration"

func removeConfigCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "remove",
		Args:  cobra.NoArgs,
		Short: removeConfigDesc,
		Long:  cmdhelpers.Long(removeConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRemoveConfig(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRemoveConfig(verbose bool) error {
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
	err = repo.Runner.Config.GitConfig.RemoveLocalGitConfiguration(repo.Runner.Config.FullConfig.Lineage)
	if err != nil {
		return err
	}
	aliasNames := maps.Keys(repo.Runner.Config.FullConfig.Aliases)
	slices.Sort(aliasNames)
	for _, aliasName := range aliasNames {
		if strings.HasPrefix(repo.Runner.Config.FullConfig.Aliases[aliasName], "town ") {
			err = repo.Runner.Frontend.RemoveGitAlias(aliasName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
