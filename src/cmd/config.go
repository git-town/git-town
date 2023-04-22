package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v8/src/cli"
	"github.com/git-town/git-town/v8/src/execute"
	"github.com/git-town/git-town/v8/src/failure"
	"github.com/git-town/git-town/v8/src/flags"
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
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		OmitBranchNames:       true,
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	fc := failure.Collector{}
	pushNewBranches := fc.Bool(run.Config.ShouldNewBranchPush())
	pushHook := fc.Bool(run.Config.PushHook())
	isOffline := fc.Bool(run.Config.IsOffline())
	deleteOrigin := fc.Bool(run.Config.ShouldShipDeleteOriginBranch())
	pullBranchStrategy := fc.PullBranchStrategy(run.Config.PullBranchStrategy())
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	hostingService := fc.HostingService(run.Config.HostingService())
	if fc.Err != nil {
		return fc.Err
	}
	fmt.Println()
	cli.PrintHeader("Branches")
	cli.PrintEntry("main branch", cli.StringSetting(run.Config.MainBranch()))
	cli.PrintEntry("perennial branches", cli.StringSetting(strings.Join(run.Config.PerennialBranches(), ", ")))
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
	cli.PrintEntry("GitHub token", cli.StringSetting(run.Config.GitHubToken()))
	cli.PrintEntry("GitLab token", cli.StringSetting(run.Config.GitLabToken()))
	cli.PrintEntry("Gitea token", cli.StringSetting(run.Config.GiteaToken()))
	fmt.Println()
	if run.Config.MainBranch() != "" {
		cli.PrintLabelAndValue("Branch Ancestry", cli.PrintableBranchAncestry(run.Config))
	}
	return nil
}
