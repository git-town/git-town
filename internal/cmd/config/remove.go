package config

import (
	"maps"
	"slices"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
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
			cliConfig := cliconfig.CliConfig{
				DryRun:  false,
				Verbose: verbose,
			}
			return executeRemoveConfig(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRemoveConfig(cliConfig cliconfig.CliConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	if err = gitconfig.RemoveLocalGitConfiguration(repo.Backend, repo.ConfigSnapshot.Local); err != nil {
		return err
	}
	aliasNames := slices.Collect(maps.Keys(repo.UnvalidatedConfig.NormalConfig.Aliases))
	slices.Sort(aliasNames)
	for _, aliasName := range aliasNames {
		if strings.HasPrefix(repo.UnvalidatedConfig.NormalConfig.Aliases[aliasName], "town ") {
			if err = gitconfig.RemoveAlias(repo.Frontend, aliasName); err != nil {
				return err
			}
		}
	}
	return nil
}
