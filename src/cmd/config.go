package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/cli/format"
	"github.com/git-town/git-town/v10/src/cli/print"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/git"
	"github.com/git-town/git-town/v10/src/gohacks"
	"github.com/spf13/cobra"
)

const configDesc = "Displays your Git Town configuration"

func configCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: "setup",
		Args:    cobra.NoArgs,
		Short:   configDesc,
		Long:    long(configDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfig(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&configCmd)
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
	config, err := determineConfigConfig(&repo.Runner)
	if err != nil {
		return err
	}
	printConfig(config)
	return nil
}

func determineConfigConfig(run *git.ProdRunner) (ConfigConfig, error) {
	fc := gohacks.FailureCollector{}
	branchTypes := run.Config.BranchTypes()
	deleteOrigin := fc.Bool(run.Config.ShouldShipDeleteOriginBranch())
	giteaToken := run.Config.GiteaToken()
	githubToken := run.Config.GitHubToken()
	gitlabToken := run.Config.GitLabToken()
	hosting := fc.Hosting(run.Config.HostingService())
	isOffline := fc.Bool(run.Config.IsOffline())
	lineage := run.Config.Lineage(run.Backend.Config.RemoveLocalConfigValue)
	pullBranchStrategy := fc.PullBranchStrategy(run.Config.PullBranchStrategy())
	pushHook := fc.Bool(run.Config.PushHook())
	pushNewBranches := fc.Bool(run.Config.ShouldNewBranchPush())
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	return ConfigConfig{
		branchTypes:        branchTypes,
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
	branchTypes        domain.BranchTypes
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
	print.Header("Branches")
	print.Entry("main branch", format.StringSetting(config.branchTypes.MainBranch.String()))
	print.Entry("perennial branches", format.StringSetting((config.branchTypes.PerennialBranches.Join(", "))))
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(config.isOffline))
	print.Entry("pull branch strategy", config.pullBranchStrategy.String())
	print.Entry("run pre-push hook", format.Bool(config.pushHook))
	print.Entry("push new branches", format.Bool(config.pushNewBranches))
	print.Entry("ship removes the remote branch", format.Bool(config.deleteOrigin))
	print.Entry("sync strategy", config.syncStrategy.String())
	print.Entry("sync with upstream", format.Bool(config.shouldSyncUpstream))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("hosting service override", format.StringSetting(config.hosting.String()))
	print.Entry("GitHub token", format.StringSetting(config.githubToken))
	print.Entry("GitLab token", format.StringSetting(config.gitlabToken))
	print.Entry("Gitea token", format.StringSetting(config.giteaToken))
	fmt.Println()
	if !config.branchTypes.MainBranch.IsEmpty() {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.lineage))
	}
}
