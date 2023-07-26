package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
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
	fc := failure.Collector{}
	pushNewBranches := fc.Bool(repo.ProdRunner.Config.ShouldNewBranchPush())
	pushHook := fc.Bool(repo.ProdRunner.Config.PushHook())
	isOffline := fc.Bool(repo.ProdRunner.Config.IsOffline())
	deleteOrigin := fc.Bool(repo.ProdRunner.Config.ShouldShipDeleteOriginBranch())
	pullBranchStrategy := fc.PullBranchStrategy(repo.ProdRunner.Config.PullBranchStrategy())
	shouldSyncUpstream := fc.Bool(repo.ProdRunner.Config.ShouldSyncUpstream())
	syncStrategy := fc.SyncStrategy(repo.ProdRunner.Config.SyncStrategy())
	hostingService := fc.HostingService(repo.ProdRunner.Config.HostingService())
	if fc.Err != nil {
		return fc.Err
	}
	fmt.Println()
	cli.PrintHeader("Branches")
	cli.PrintEntry("main branch", cli.StringSetting(repo.ProdRunner.Config.MainBranch()))
	cli.PrintEntry("perennial branches", cli.StringSetting(strings.Join(repo.ProdRunner.Config.PerennialBranches(), ", ")))
	fmt.Println()
	cli.PrintHeader("Configuration")
	cli.PrintEntry("offline", cli.BoolSetting(isOffline))
	cli.PrintEntry("pull branch strategy", string(pullBranchStrategy))
	cli.PrintEntry("run pre-push hook", cli.BoolSetting(pushHook))
	cli.PrintEntry("push new branches", cli.BoolSetting(pushNewBranches))
	cli.PrintEntry("ship removes the remote branch", cli.BoolSetting(deleteOrigin))
	cli.PrintEntry("sync strategy", string(syncStrategy))
	cli.PrintEntry("sync with upstream", cli.BoolSetting(shouldSyncUpstream))
	fmt.Println()
	cli.PrintHeader("Hosting")
	cli.PrintEntry("hosting service override", cli.StringSetting(string(hostingService)))
	cli.PrintEntry("GitHub token", cli.StringSetting(repo.ProdRunner.Config.GitHubToken()))
	cli.PrintEntry("GitLab token", cli.StringSetting(repo.ProdRunner.Config.GitLabToken()))
	cli.PrintEntry("Gitea token", cli.StringSetting(repo.ProdRunner.Config.GiteaToken()))
	fmt.Println()
	if repo.ProdRunner.Config.MainBranch() != "" {
		cli.PrintLabelAndValue("Branch Lineage", cli.PrintableBranchLineage(repo.ProdRunner.Config.Lineage()))
	}
	return nil
}
