// Package config implements Git Town's "config" command.
package config

import (
	"cmp"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/format"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const configDesc = "Display your Git Town configuration"

func RootCmd() *cobra.Command {
	addDisplayTypesFlag, readDisplayTypesFlag := flags.Displaytypes()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: cmdhelpers.GroupIDSetup,
		Args:    cobra.NoArgs,
		Short:   configDesc,
		Long:    cmdhelpers.Long(configDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			displayTypes, errDisplayTypes := readDisplayTypesFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errDisplayTypes, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:  None[configdomain.AutoResolve](),
				AutoSync:     None[configdomain.AutoSync](),
				Detached:     None[configdomain.Detached](),
				DisplayTypes: displayTypes,
				DryRun:       None[configdomain.DryRun](),
				Order:        None[configdomain.Order](),
				PushBranches: None[configdomain.PushBranches](),
				Stash:        None[configdomain.Stash](),
				Verbose:      verbose,
			})
			return executeDisplayConfig(cliConfig)
		},
	}
	addDisplayTypesFlag(&configCmd)
	addVerboseFlag(&configCmd)
	configCmd.AddCommand(getParentCommand())
	configCmd.AddCommand(removeConfigCommand())
	return &configCmd
}

func executeDisplayConfig(cliConfig configdomain.PartialConfig) error {
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
	print.Entry("order", config.NormalConfig.Order.String())
	print.Entry("display types", config.NormalConfig.DisplayTypes.String())
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(config.NormalConfig.Offline.IsOffline()))
	print.Entry("git user name", format.OptionalStringerSetting(config.NormalConfig.GitUserName))
	print.Entry("git user email", format.OptionalStringerSetting(config.NormalConfig.GitUserEmail))
	fmt.Println()
	print.Header("Create")
	print.Entry("new branch type", format.OptionalStringerSetting(config.NormalConfig.NewBranchType))
	print.Entry("share new branches", config.NormalConfig.ShareNewBranches.String())
	print.Entry("stash uncommitted changes", format.Bool(config.NormalConfig.Stash.ShouldStash()))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("development remote", config.NormalConfig.DevRemote.String())
	print.Entry("forge type", format.OptionalStringerSetting(config.NormalConfig.ForgeType))
	print.Entry("origin hostname", format.OptionalStringerSetting(config.NormalConfig.HostingOriginHostname))
	print.Entry("Bitbucket username", format.OptionalStringerSetting(config.NormalConfig.BitbucketUsername))
	print.Entry("Bitbucket app password", format.OptionalStringerSetting(config.NormalConfig.BitbucketAppPassword))
	print.Entry("Forgejo token", format.OptionalStringerSetting(config.NormalConfig.ForgejoToken))
	print.Entry("Gitea token", format.OptionalStringerSetting(config.NormalConfig.GiteaToken))
	print.Entry("GitHub connector type", format.OptionalStringerSetting(config.NormalConfig.GitHubConnectorType))
	print.Entry("GitHub token", format.OptionalStringerSetting(config.NormalConfig.GitHubToken))
	print.Entry("GitLab connector type", format.OptionalStringerSetting(config.NormalConfig.GitLabConnectorType))
	print.Entry("GitLab token", format.OptionalStringerSetting(config.NormalConfig.GitLabToken))
	fmt.Println()
	print.Header("Proposals")
	print.Entry("show lineage", format.StringsSetting(config.NormalConfig.ProposalsShowLineage.String()))
	fmt.Println()
	print.Header("Ship")
	print.Entry("delete tracking branch", format.Bool(config.NormalConfig.ShipDeleteTrackingBranch.ShouldDeleteTrackingBranch()))
	print.Entry("ship strategy", config.NormalConfig.ShipStrategy.String())
	fmt.Println()
	print.Header("Sync")
	print.Entry("auto-resolve phantom conflicts", format.Bool(config.NormalConfig.AutoResolve.ShouldAutoResolve()))
	print.Entry("auto-sync", format.Bool(config.NormalConfig.AutoSync.ShouldSync()))
	print.Entry("run detached", format.Bool(config.NormalConfig.Detached.ShouldWorkDetached()))
	print.Entry("run pre-push hook", format.Bool(config.NormalConfig.PushHook.ShouldRunPushHook()))
	print.Entry("feature sync strategy", config.NormalConfig.SyncFeatureStrategy.String())
	print.Entry("perennial sync strategy", config.NormalConfig.SyncPerennialStrategy.String())
	print.Entry("prototype sync strategy", config.NormalConfig.SyncPrototypeStrategy.String())
	print.Entry("push branches", format.Bool(config.NormalConfig.PushBranches.ShouldPush()))
	print.Entry("sync tags", format.Bool(config.NormalConfig.SyncTags.ShouldSyncTags()))
	print.Entry("sync with upstream", format.Bool(config.NormalConfig.SyncUpstream.ShouldSyncUpstream()))
	print.Entry("auto-resolve phantom conflicts", format.Bool(config.NormalConfig.AutoResolve.ShouldAutoResolve()))
	fmt.Println()
	if config.NormalConfig.Lineage.Len() > 0 {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.NormalConfig.Lineage, config.NormalConfig.Order))
	}
}
