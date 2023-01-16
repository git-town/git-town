package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Displays your Git Town configuration",
		Run: func(cmd *cobra.Command, args []string) {
			pushNewBranches, err := prodRepo.Config.ShouldNewBranchPush()
			if err != nil {
				cli.Exit(err)
			}
			pushHook, err := prodRepo.Config.PushHook()
			if err != nil {
				cli.Exit(err)
			}
			isOffline, err := prodRepo.Config.IsOffline()
			if err != nil {
				cli.Exit(err)
			}
			deleteOrigin, err := prodRepo.Config.ShouldShipDeleteOriginBranch()
			if err != nil {
				cli.Exit(err)
			}
			shouldSyncUpstream, err := prodRepo.Config.ShouldSyncUpstream()
			if err != nil {
				cli.Exit(err)
			}
			fmt.Println()
			cli.PrintHeader("Branches")
			cli.PrintEntry("main branch", cli.StringSetting(prodRepo.Config.MainBranch()))
			cli.PrintEntry("perennial branches", cli.StringSetting(strings.Join(prodRepo.Config.PerennialBranches(), ", ")))
			fmt.Println()
			cli.PrintHeader("Configuration")
			cli.PrintEntry("offline", cli.BoolSetting(isOffline))
			cli.PrintEntry("pull branch strategy", prodRepo.Config.PullBranchStrategy())
			cli.PrintEntry("run pre-push hook", cli.BoolSetting(pushHook))
			cli.PrintEntry("push new branches", cli.BoolSetting(pushNewBranches))
			cli.PrintEntry("ship removes the remote branch", cli.BoolSetting(deleteOrigin))
			cli.PrintEntry("sync strategy", prodRepo.Config.SyncStrategy())
			cli.PrintEntry("sync with upstream", cli.BoolSetting(shouldSyncUpstream))
			fmt.Println()
			cli.PrintHeader("Hosting")
			cli.PrintEntry("hosting service override", cli.StringSetting(prodRepo.Config.HostingService()))
			cli.PrintEntry("GitHub token", cli.StringSetting(prodRepo.Config.GitHubToken()))
			cli.PrintEntry("GitLab token", cli.StringSetting(prodRepo.Config.GitLabToken()))
			cli.PrintEntry("Gitea token", cli.StringSetting(prodRepo.Config.GiteaToken()))
			fmt.Println()
			if prodRepo.Config.MainBranch() != "" {
				cli.PrintLabelAndValue("Branch Ancestry", cli.PrintableBranchAncestry(&prodRepo.Config))
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ValidateIsRepository(prodRepo)
		},
	}
	configCmd.AddCommand(mainbranchConfigCmd())
	configCmd.AddCommand(offlineCmd())
	configCmd.AddCommand(perennialBranchesCmd())
	configCmd.AddCommand(pullBranchStrategyCommand())
	configCmd.AddCommand(pushNewBranchesCommand())
	configCmd.AddCommand(pushHookCommand())
	configCmd.AddCommand(resetConfigCommand())
	configCmd.AddCommand(setupConfigCommand())
	configCmd.AddCommand(syncStrategyCommand())

	return configCmd
}

// SYNC-STRATEGY SUBCOMMAND
