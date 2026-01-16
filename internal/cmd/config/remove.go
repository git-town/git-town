package config

import (
	"maps"
	"slices"
	"strings"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeRemoveConfig(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRemoveConfig(cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    true,
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
	slice.NaturalSort(aliasNames)
	for _, aliasName := range aliasNames {
		if strings.HasPrefix(repo.UnvalidatedConfig.NormalConfig.Aliases[aliasName], "town ") {
			if err = gitconfig.RemoveAlias(repo.Frontend, aliasName); err != nil {
				return err
			}
		}
	}
	return nil
}
