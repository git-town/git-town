package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	debug := false
	configCmd := cobra.Command{
		Use:     "config",
		GroupID: "setup",
		Args:    cobra.NoArgs,
		Short:   "Displays your Git Town configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfig(debug)
		},
	}
	debugFlag(&configCmd, &debug)
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
	repo := Repo(debug, false)
	err := ensure(&repo, isRepository)
	if err != nil {
		return err
	}
	ec := runstate.ErrorChecker{}
	pushNewBranches := ec.Bool(repo.Config.ShouldNewBranchPush())
	pushHook := ec.Bool(repo.Config.PushHook())
	isOffline := ec.Bool(repo.Config.IsOffline())
	deleteOrigin := ec.Bool(repo.Config.ShouldShipDeleteOriginBranch())
	pullBranchStrategy := ec.PullBranchStrategy(repo.Config.PullBranchStrategy())
	shouldSyncUpstream := ec.Bool(repo.Config.ShouldSyncUpstream())
	syncStrategy := ec.SyncStrategy(repo.Config.SyncStrategy())
	hostingService := ec.HostingService(repo.Config.HostingService())
	if ec.Err != nil {
		return ec.Err
	}
	fmt.Println()
	cli.PrintHeader("Branches")
	cli.PrintEntry("main branch", cli.StringSetting(repo.Config.MainBranch()))
	cli.PrintEntry("perennial branches", cli.StringSetting(strings.Join(repo.Config.PerennialBranches(), ", ")))
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
	cli.PrintEntry("GitHub token", cli.StringSetting(repo.Config.GitHubToken()))
	cli.PrintEntry("GitLab token", cli.StringSetting(repo.Config.GitLabToken()))
	cli.PrintEntry("Gitea token", cli.StringSetting(repo.Config.GiteaToken()))
	fmt.Println()
	if repo.Config.MainBranch() != "" {
		cli.PrintLabelAndValue("Branch Ancestry", cli.PrintableBranchAncestry(repo.Config))
	}
	return nil
}
