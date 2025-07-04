package config

import (
	"maps"
	"slices"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/spf13/cobra"
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
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeRemoveConfig(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRemoveConfig(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	if err = repo.UnvalidatedConfig.NormalConfig.GitIO.RemoveLocalGitConfiguration(repo.ConfigSnapshot.Local); err != nil {
		return err
	}
	aliasNames := slices.Collect(maps.Keys(repo.UnvalidatedConfig.NormalConfig.Aliases))
	slices.Sort(aliasNames)
	for _, aliasName := range aliasNames {
		if strings.HasPrefix(repo.UnvalidatedConfig.NormalConfig.Aliases[aliasName], "town ") {
			if err = repo.Git.RemoveGitAlias(repo.Frontend, aliasName); err != nil {
				return err
			}
		}
	}
	return nil
}
