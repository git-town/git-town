// Package config implements Git Town's "config" command.
package config

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/format"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
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
			return executeDisplayConfig(readVerboseFlag(cmd))
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
	printConfig(*repo.NormalConfig.Config.Value, repo.UnvalidatedConfig.Config.Get())
	return nil
}

func printConfig(normalConfig configdomain.NormalConfig, unvalidatedConfig configdomain.UnvalidatedConfig) {
	fmt.Println()
	print.Header("Branches")
	print.Entry("main branch", format.StringSetting(unvalidatedConfig.MainBranch.String()))
	print.Entry("perennial branches", format.StringsSetting((normalConfig.PerennialBranches.Join(", "))))
	print.Entry("perennial regex", format.StringSetting(normalConfig.PerennialRegex.String()))
	print.Entry("default branch type", format.StringSetting(normalConfig.DefaultBranchType.String()))
	print.Entry("feature regex", format.StringSetting(normalConfig.FeatureRegex.String()))
	print.Entry("parked branches", format.StringsSetting((normalConfig.ParkedBranches.Join(", "))))
	print.Entry("contribution branches", format.StringsSetting((normalConfig.ContributionBranches.Join(", "))))
	print.Entry("contribution regex", format.StringsSetting((normalConfig.ContributionRegex.String())))
	print.Entry("observed branches", format.StringsSetting((normalConfig.ObservedBranches.Join(", "))))
	print.Entry("observed regex", format.StringsSetting((normalConfig.ObservedRegex.String())))
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(normalConfig.Offline.IsTrue()))
	print.Entry("run pre-push hook", format.Bool(bool(normalConfig.PushHook)))
	print.Entry("push new branches", format.Bool(normalConfig.ShouldPushNewBranches()))
	print.Entry("ship strategy", normalConfig.ShipStrategy.String())
	print.Entry("ship deletes the tracking branch", format.Bool(normalConfig.ShipDeleteTrackingBranch.IsTrue()))
	print.Entry("sync-feature strategy", normalConfig.SyncFeatureStrategy.String())
	print.Entry("sync-perennial strategy", normalConfig.SyncPerennialStrategy.String())
	print.Entry("sync with upstream", format.Bool(normalConfig.SyncUpstream.IsTrue()))
	print.Entry("sync tags", format.Bool(normalConfig.SyncTags.IsTrue()))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("hosting platform override", format.StringSetting(normalConfig.HostingPlatform.String()))
	print.Entry("GitHub token", format.OptionalStringerSetting(normalConfig.GitHubToken))
	print.Entry("GitLab token", format.OptionalStringerSetting(normalConfig.GitLabToken))
	print.Entry("Gitea token", format.OptionalStringerSetting(normalConfig.GiteaToken))
	fmt.Println()
	if normalConfig.Lineage.Len() > 0 {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(normalConfig.Lineage))
	}
}
