package config

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cli/format"
	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
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
	printConfig(repo.UnvalidatedConfig.Config.Get())
	return nil
}

func printConfig(config configdomain.UnvalidatedConfig) {
	fmt.Println()
	print.Header("Branches")
	print.Entry("main branch", format.StringSetting(config.MainBranch.String()))
	print.Entry("perennial branches", format.StringsSetting((config.PerennialBranches.Join(", "))))
	print.Entry("perennial regex", format.StringSetting(config.PerennialRegex.String()))
	print.Entry("parked branches", format.StringsSetting((config.ParkedBranches.Join(", "))))
	print.Entry("contribution branches", format.StringsSetting((config.ContributionBranches.Join(", "))))
	print.Entry("observed branches", format.StringsSetting((config.ObservedBranches.Join(", "))))
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(config.Offline.IsTrue()))
	print.Entry("run pre-push hook", format.Bool(bool(config.PushHook)))
	print.Entry("push new branches", format.Bool(config.ShouldPushNewBranches()))
	print.Entry("ship deletes the tracking branch", format.Bool(config.ShipDeleteTrackingBranch.IsTrue()))
	print.Entry("sync-feature strategy", config.SyncFeatureStrategy.String())
	print.Entry("sync-perennial strategy", config.SyncPerennialStrategy.String())
	print.Entry("sync with upstream", format.Bool(config.SyncUpstream.IsTrue()))
	print.Entry("sync tags", format.Bool(config.SyncTags.IsTrue()))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("hosting platform override", format.StringSetting(config.HostingPlatform.String()))
	print.Entry("GitHub token", format.OptionalStringerSetting(config.GitHubToken))
	print.Entry("GitLab token", format.OptionalStringerSetting(config.GitLabToken))
	print.Entry("Gitea token", format.OptionalStringerSetting(config.GiteaToken))
	fmt.Println()
	if config.Lineage.Len() > 0 {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.Lineage))
	}
}
