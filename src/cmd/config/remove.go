package config

import (
	"slices"
	"strings"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/execute"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
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
	err = repo.UnvalidatedConfig.GitConfig.RemoveLocalGitConfiguration(repo.UnvalidatedConfig.Config.Value.Lineage)
	if err != nil {
		return err
	}
	aliasNames := maps.Keys(repo.UnvalidatedConfig.Config.Value.Aliases)
	slices.Sort(aliasNames)
	for _, aliasName := range aliasNames {
		if strings.HasPrefix(repo.UnvalidatedConfig.Config.Value.Aliases[aliasName], "town ") {
			err = repo.Git.RemoveGitAlias(repo.Frontend, aliasName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
