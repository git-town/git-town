package config

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/sync/syncdomain"
	"github.com/spf13/cobra"
)

const configDesc = "Displays your Git Town configuration"

func RootCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: "setup",
		Args:    cobra.NoArgs,
		Short:   configDesc,
		Long:    cmdhelpers.Long(configDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfig(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&configCmd)
	configCmd.AddCommand(mainbranchConfigCmd())
	configCmd.AddCommand(offlineCmd())
	configCmd.AddCommand(perennialBranchesCmd())
	configCmd.AddCommand(pushHookCommand())
	configCmd.AddCommand(pushNewBranchesCommand())
	configCmd.AddCommand(resetConfigCommand())
	configCmd.AddCommand(setupConfigCommand())
	configCmd.AddCommand(syncFeatureStrategyCommand())
	configCmd.AddCommand(syncPerennialStrategyCommand())
	return &configCmd
}

func executeConfig(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, err := determineConfigConfig(repo.Runner)
	if err != nil {
		return err
	}
	printConfig(config)
	return nil
}

func determineConfigConfig(run *git.ProdRunner) (RootConfig, error) {
	fc := execute.FailureCollector{}
	branchTypes := run.GitTown.BranchTypes()
	deleteTrackingBranch := run.GitTown.ShipDeleteTrackingBranch
	giteaToken := run.GitTown.GiteaToken
	githubToken := run.GitTown.GitHubToken
	gitlabToken := run.GitTown.GitLabToken
	hosting := fc.Hosting(run.GitTown.HostingService())
	isOffline := run.GitTown.Offline
	lineage := run.GitTown.Lineage(run.Backend.GitTown.RemoveLocalConfigValue)
	syncPerennialStrategy := run.GitTown.SyncPerennialStrategy
	pushHook := run.GitTown.PushHook
	pushNewBranches := run.GitTown.NewBranchPush
	syncUpstream := run.GitTown.SyncUpstream
	syncFeatureStrategy := run.GitTown.SyncFeatureStrategy
	syncBeforeShip := run.GitTown.SyncBeforeShip
	return RootConfig{
		branchTypes:           branchTypes,
		deleteTrackingBranch:  deleteTrackingBranch,
		hosting:               hosting,
		giteaToken:            giteaToken,
		githubToken:           githubToken,
		gitlabToken:           gitlabToken,
		isOffline:             isOffline,
		lineage:               lineage,
		syncPerennialStrategy: syncPerennialStrategy,
		pushHook:              pushHook,
		pushNewBranches:       pushNewBranches,
		syncUpstream:          syncUpstream,
		syncFeatureStrategy:   syncFeatureStrategy,
		syncBeforeShip:        syncBeforeShip,
	}, fc.Err
}

type RootConfig struct {
	branchTypes           syncdomain.BranchTypes
	deleteTrackingBranch  configdomain.ShipDeleteTrackingBranch
	giteaToken            configdomain.GiteaToken
	githubToken           configdomain.GitHubToken
	gitlabToken           configdomain.GitLabToken
	hosting               configdomain.Hosting
	isOffline             configdomain.Offline
	lineage               configdomain.Lineage
	syncPerennialStrategy configdomain.SyncPerennialStrategy
	pushHook              configdomain.PushHook
	pushNewBranches       configdomain.NewBranchPush
	syncUpstream          configdomain.SyncUpstream
	syncFeatureStrategy   configdomain.SyncFeatureStrategy
	syncBeforeShip        configdomain.SyncBeforeShip
}

func printConfig(config RootConfig) {
	fmt.Println()
	print.Header("Branches")
	print.Entry("main branch", format.StringSetting(config.branchTypes.MainBranch.String()))
	print.Entry("perennial branches", format.StringSetting((config.branchTypes.PerennialBranches.Join(", "))))
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(config.isOffline.Bool()))
	print.Entry("run pre-push hook", format.Bool(bool(config.pushHook)))
	print.Entry("push new branches", format.Bool(config.pushNewBranches.Bool()))
	print.Entry("ship deletes the tracking branch", format.Bool(config.deleteTrackingBranch.Bool()))
	print.Entry("sync-feature strategy", config.syncFeatureStrategy.String())
	print.Entry("sync-perennial strategy", config.syncPerennialStrategy.String())
	print.Entry("sync with upstream", format.Bool(config.syncUpstream.Bool()))
	print.Entry("sync before shipping", format.Bool(config.syncBeforeShip.Bool()))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("hosting service override", format.StringSetting(config.hosting.String()))
	print.Entry("GitHub token", format.StringSetting(string(config.githubToken)))
	print.Entry("GitLab token", format.StringSetting(string(config.gitlabToken)))
	print.Entry("Gitea token", format.StringSetting(string(config.giteaToken)))
	fmt.Println()
	if !config.branchTypes.MainBranch.IsEmpty() {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.lineage))
	}
}
