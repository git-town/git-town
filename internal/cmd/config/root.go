// Package config implements Git Town's "config" command.
package config

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/format"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/spf13/cobra"
)

const configDesc = "Display your Git Town configuration"

func RootCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: "setup",
		Args:    cobra.NoArgs,
		Short:   configDesc,
		Long:    cmdhelpers.Long(configDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeDisplayConfig(verbose)
		},
	}
	addVerboseFlag(&configCmd)
	configCmd.AddCommand(getParentCommand())
	configCmd.AddCommand(removeConfigCommand())
	configCmd.AddCommand(SetupCommand())
	return &configCmd
}

func executeDisplayConfig(verbose configdomain.Verbose) error {
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
	printConfig(repo.UnvalidatedConfig)
	return nil
}

func printConfig(config config.UnvalidatedConfig) {
	fmt.Println()
	// TODO: organize these entries exactly like the config file is organized
	print.Header("Branches")
	print.Entry("main branch", format.OptionalStringerSetting(config.UnvalidatedConfig.MainBranch))
	print.Entry("perennial branches", format.StringsSetting((config.NormalConfig.PerennialBranches.Join(", "))))
	print.Entry("perennial regex", format.OptionalStringerSetting(config.NormalConfig.PerennialRegex))
	print.Entry("default branch type", config.NormalConfig.DefaultBranchType.String())
	print.Entry("feature regex", format.OptionalStringerSetting(config.NormalConfig.FeatureRegex))
	print.Entry("parked branches", format.StringsSetting((config.NormalConfig.ParkedBranches.Join(", "))))
	print.Entry("contribution branches", format.StringsSetting((config.NormalConfig.ContributionBranches.Join(", "))))
	print.Entry("contribution regex", format.OptionalStringerSetting((config.NormalConfig.ContributionRegex)))
	print.Entry("observed branches", format.StringsSetting((config.NormalConfig.ObservedBranches.Join(", "))))
	print.Entry("observed regex", format.OptionalStringerSetting((config.NormalConfig.ObservedRegex)))
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(config.NormalConfig.Offline.IsTrue()))
	print.Entry("run pre-push hook", format.Bool(bool(config.NormalConfig.PushHook)))
	print.Entry("push new branches", format.Bool(config.NormalConfig.ShouldPushNewBranches()))
	print.Entry("ship strategy", config.NormalConfig.ShipStrategy.String())
	print.Entry("ship deletes the tracking branch", format.Bool(config.NormalConfig.ShipDeleteTrackingBranch.IsTrue()))
	print.Entry("sync-feature strategy", config.NormalConfig.SyncFeatureStrategy.String())
	print.Entry("sync-perennial strategy", config.NormalConfig.SyncPerennialStrategy.String())
	print.Entry("sync-prototype strategy", config.NormalConfig.SyncPrototypeStrategy.String())
	print.Entry("sync with upstream", format.Bool(config.NormalConfig.SyncUpstream.IsTrue()))
	print.Entry("sync tags", format.Bool(config.NormalConfig.SyncTags.IsTrue()))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("hosting platform override", format.OptionalStringerSetting(config.NormalConfig.HostingPlatform))
	print.Entry("GitHub token", format.OptionalStringerSetting(config.NormalConfig.GitHubToken))
	print.Entry("GitLab token", format.OptionalStringerSetting(config.NormalConfig.GitLabToken))
	print.Entry("Gitea token", format.OptionalStringerSetting(config.NormalConfig.GiteaToken))
	fmt.Println()
	if config.NormalConfig.Lineage.Len() > 0 {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.NormalConfig.Lineage))
	}
}
