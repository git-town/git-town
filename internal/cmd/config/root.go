// Package config implements Git Town's "config" command.
package config

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/spf13/cobra"
)

const configDesc = "Display your Git Town configuration"

func RootCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: cmdhelpers.GroupIDSetup,
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
	print.Header("Branches")
	print.Entry("contribution branches", format.BranchNames(config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypeContributionBranch)))
	print.Entry("contribution regex", format.OptionalStringerSetting(config.NormalConfig.ContributionRegex))
	print.Entry("feature regex", format.OptionalStringerSetting(config.NormalConfig.FeatureRegex))
	print.Entry("main branch", format.OptionalStringerSetting(config.UnvalidatedConfig.MainBranch))
	print.Entry("observed branches", format.BranchNames(config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypeObservedBranch)))
	print.Entry("observed regex", format.OptionalStringerSetting(config.NormalConfig.ObservedRegex))
	print.Entry("parked branches", format.BranchNames(config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypeParkedBranch)))
	print.Entry("perennial branches", format.StringsSetting(config.NormalConfig.PerennialBranches.Join(", ")))
	print.Entry("perennial regex", format.OptionalStringerSetting(config.NormalConfig.PerennialRegex))
	print.Entry("prototype branches", format.BranchNames(config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypePrototypeBranch)))
	print.Entry("unknown branch type", config.NormalConfig.UnknownBranchType.String())
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(config.NormalConfig.Offline.IsOffline()))
	fmt.Println()
	print.Header("Create")
	print.Entry("new branch type", format.OptionalStringerSetting(config.NormalConfig.NewBranchType))
	print.Entry("share new branches", config.NormalConfig.ShareNewBranches.String())
	fmt.Println()
	print.Header("Hosting")
	print.Entry("development remote", config.NormalConfig.DevRemote.String())
	print.Entry("forge type", format.OptionalStringerSetting(config.NormalConfig.ForgeType))
	print.Entry("origin hostname", format.OptionalStringerSetting(config.NormalConfig.HostingOriginHostname))
	print.Entry("Bitbucket username", format.OptionalStringerSetting(config.NormalConfig.BitbucketUsername))
	print.Entry("Bitbucket app password", format.OptionalStringerSetting(config.NormalConfig.BitbucketAppPassword))
	print.Entry("Codeberg token", format.OptionalStringerSetting(config.NormalConfig.CodebergToken))
	print.Entry("Gitea token", format.OptionalStringerSetting(config.NormalConfig.GiteaToken))
	print.Entry("GitHub connector type", format.OptionalStringerSetting(config.NormalConfig.GitHubConnectorType))
	print.Entry("GitHub token", format.OptionalStringerSetting(config.NormalConfig.GitHubToken))
	print.Entry("GitLab connector type", format.OptionalStringerSetting(config.NormalConfig.GitLabConnectorType))
	print.Entry("GitLab token", format.OptionalStringerSetting(config.NormalConfig.GitLabToken))
	fmt.Println()
	print.Header("Ship")
	print.Entry("delete tracking branch", format.Bool(config.NormalConfig.ShipDeleteTrackingBranch.IsTrue()))
	print.Entry("ship strategy", config.NormalConfig.ShipStrategy.String())
	fmt.Println()
	print.Header("Sync")
	print.Entry("run pre-push hook", format.Bool(bool(config.NormalConfig.PushHook)))
	print.Entry("feature sync strategy", config.NormalConfig.SyncFeatureStrategy.String())
	print.Entry("perennial sync strategy", config.NormalConfig.SyncPerennialStrategy.String())
	print.Entry("prototype sync strategy", config.NormalConfig.SyncPrototypeStrategy.String())
	print.Entry("sync tags", format.Bool(config.NormalConfig.SyncTags.IsTrue()))
	print.Entry("sync with upstream", format.Bool(config.NormalConfig.SyncUpstream.IsTrue()))
	fmt.Println()
	if config.NormalConfig.Lineage.Len() > 0 {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.NormalConfig.Lineage))
	}
}
