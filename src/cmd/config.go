package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/spf13/cobra"
)

const configDesc = "Displays your Git Town configuration"

func configCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: "setup",
		Args:    cobra.NoArgs,
		Short:   configDesc,
		Long:    long(configDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfig(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&configCmd)
	configCmd.AddCommand(mainbranchConfigCmd())
	configCmd.AddCommand(offlineCmd())
	configCmd.AddCommand(perennialBranchesCmd())
	configCmd.AddCommand(pullBranchStrategyCommand())
	configCmd.AddCommand(pushNewBranchesCommand())
	configCmd.AddCommand(pushHookCommand())
	configCmd.AddCommand(resetConfigCommand())
	configCmd.AddCommand(setupConfigCommand())
	configCmd.AddCommand(syncStrategyCommand())
	return &configCmd
}

func runConfig(debug bool) error {
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 false,
		HandleUnfinishedState: false,
		OmitBranchNames:       true,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineConfigConfig(&repo.Runner)
	if err != nil {
		return err
	}
	printConfig(config)
	return nil
}

func determineConfigConfig(run *git.ProdRunner) (ConfigConfig, error) {
	fc := failure.Collector{}
	branchDurations := run.Config.BranchDurations()
	deleteOrigin := fc.Bool(run.Config.ShouldShipDeleteOriginBranch())
	giteaToken := run.Config.GiteaToken()
	githubToken := run.Config.GitHubToken()
	gitlabToken := run.Config.GitLabToken()
	hosting := fc.Hosting(run.Config.HostingService())
	isOffline := fc.Bool(run.Config.IsOffline())
	lineage := run.Config.Lineage()
	pullBranchStrategy := fc.PullBranchStrategy(run.Config.PullBranchStrategy())
	pushHook := fc.Bool(run.Config.PushHook())
	pushNewBranches := fc.Bool(run.Config.ShouldNewBranchPush())
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	return ConfigConfig{
		branchDurations:    branchDurations,
		deleteOrigin:       deleteOrigin,
		hosting:            hosting,
		giteaToken:         giteaToken,
		githubToken:        githubToken,
		gitlabToken:        gitlabToken,
		isOffline:          isOffline,
		lineage:            lineage,
		pullBranchStrategy: pullBranchStrategy,
		pushHook:           pushHook,
		pushNewBranches:    pushNewBranches,
		shouldSyncUpstream: shouldSyncUpstream,
		syncStrategy:       syncStrategy,
	}, fc.Err
}

type ConfigConfig struct {
	branchDurations    domain.BranchDurations
	deleteOrigin       bool
	giteaToken         string
	githubToken        string
	gitlabToken        string
	hosting            config.Hosting
	isOffline          bool
	lineage            config.Lineage
	pullBranchStrategy config.PullBranchStrategy
	pushHook           bool
	pushNewBranches    bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func printConfig(config ConfigConfig) {
	fmt.Println()
	cli.PrintHeader("Branches")
	cli.PrintEntry("main branch", cli.StringSetting(config.branchDurations.MainBranch.String()))
	cli.PrintEntry("perennial branches", cli.StringSetting((config.branchDurations.PerennialBranches.Join(", "))))
	fmt.Println()
	cli.PrintHeader("Configuration")
	cli.PrintEntry("offline", cli.BoolSetting(config.isOffline))
	cli.PrintEntry("pull branch strategy", config.pullBranchStrategy.String())
	cli.PrintEntry("run pre-push hook", cli.BoolSetting(config.pushHook))
	cli.PrintEntry("push new branches", cli.BoolSetting(config.pushNewBranches))
	cli.PrintEntry("ship removes the remote branch", cli.BoolSetting(config.deleteOrigin))
	cli.PrintEntry("sync strategy", config.syncStrategy.String())
	cli.PrintEntry("sync with upstream", cli.BoolSetting(config.shouldSyncUpstream))
	fmt.Println()
	cli.PrintHeader("Hosting")
	cli.PrintEntry("hosting service override", cli.StringSetting(config.hosting.String()))
	cli.PrintEntry("GitHub token", cli.StringSetting(config.githubToken))
	cli.PrintEntry("GitLab token", cli.StringSetting(config.gitlabToken))
	cli.PrintEntry("Gitea token", cli.StringSetting(config.giteaToken))
	fmt.Println()
	if !config.branchDurations.MainBranch.IsEmpty() {
		cli.PrintLabelAndValue("Branch Lineage", cli.PrintableBranchLineage(config.lineage))
	}
}
