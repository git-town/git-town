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
	addRedactFlag, readRedactFlag := flags.Redact()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: cmdhelpers.GroupIDConfig,
		Args:    cobra.NoArgs,
		Short:   configDesc,
		Long:    cmdhelpers.Long(configDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			displayTypes, errDisplayTypes := readDisplayTypesFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			redact, errRedact := readRedactFlag(cmd)
			if err := cmp.Or(errDisplayTypes, errRedact, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      displayTypes,
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeDisplayConfig(cliConfig, redact)
		},
	}
	addDisplayTypesFlag(&configCmd)
	addVerboseFlag(&configCmd)
	addRedactFlag(&configCmd)
	configCmd.AddCommand(getParentCommand())
	configCmd.AddCommand(removeConfigCommand())
	return &configCmd
}

func executeDisplayConfig(cliConfig configdomain.PartialConfig, redact configdomain.Redact) error {
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
	printConfig(repo.UnvalidatedConfig, redact)
	return nil
}

func printConfig(config config.UnvalidatedConfig, redact configdomain.Redact) {
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
	print.Entry("git user email", formatToken(config.NormalConfig.GitUserEmail, redact))
	fmt.Println()
	print.Header("Create")
	print.Entry("branch prefix", format.OptionalStringerSetting(config.NormalConfig.BranchPrefix))
	print.Entry("new branch type", format.OptionalStringerSetting(config.NormalConfig.NewBranchType))
	print.Entry("share new branches", config.NormalConfig.ShareNewBranches.String())
	print.Entry("stash uncommitted changes", format.Bool(config.NormalConfig.Stash.ShouldStash()))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("browser", format.OptionalStringerSetting(config.NormalConfig.Browser))
	print.Entry("development remote", config.NormalConfig.DevRemote.String())
	print.Entry("forge type", format.OptionalStringerSetting(config.NormalConfig.ForgeType))
	print.Entry("origin hostname", format.OptionalStringerSetting(config.NormalConfig.HostingOriginHostname))
	print.Entry("Bitbucket username", format.OptionalStringerSetting(config.NormalConfig.BitbucketUsername))
	print.Entry("Bitbucket app password", formatToken(config.NormalConfig.BitbucketAppPassword, redact))
	print.Entry("Forgejo token", formatToken(config.NormalConfig.ForgejoToken, redact))
	print.Entry("Gitea token", formatToken(config.NormalConfig.GiteaToken, redact))
	print.Entry("GitHub connector", format.OptionalStringerSetting(config.NormalConfig.GithubConnectorType))
	print.Entry("GitHub token", formatToken(config.NormalConfig.GithubToken, redact))
	print.Entry("GitLab connector", format.OptionalStringerSetting(config.NormalConfig.GitlabConnectorType))
	print.Entry("GitLab token", formatToken(config.NormalConfig.GitlabToken, redact))
	fmt.Println()
	print.Header("Propose")
	print.Entry("breadcrumb", format.StringsSetting(config.NormalConfig.ProposalBreadcrumb.String()))
	print.Entry("breadcrumb direction", format.StringsSetting(config.NormalConfig.ProposalBreadcrumbDirection.String()))
	fmt.Println()
	print.Header("Ship")
	print.Entry("delete tracking branch", format.Bool(config.NormalConfig.ShipDeleteTrackingBranch.ShouldDeleteTrackingBranch()))
	print.Entry("ignore uncommitted changes", format.Bool(config.NormalConfig.IgnoreUncommitted.AllowUncommitted()))
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

// formatToken returns a formatted token value. If redact is true and the token is set, it returns "(configured)".
func formatToken[T fmt.Stringer](token Option[T], redact configdomain.Redact) string {
	if redact.ShouldRedact() && token.IsSome() {
		return "(configured)"
	}
	return format.OptionalStringerSetting(token)
}
