package config

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
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
	configCmd.AddCommand(removeConfigCommand())
	configCmd.AddCommand(SetupCommand())
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
	printConfig(&repo.Runner.FullConfig)
	return nil
}

func printConfig(config *configdomain.FullConfig) {
	fmt.Println()
	print.Header("Branches")
	print.Entry("main branch", format.StringSetting(config.MainBranch.String()))
	print.Entry("perennial branches", format.StringSetting((config.PerennialBranches.Join(", "))))
	fmt.Println()
	print.Header("Configuration")
	print.Entry("offline", format.Bool(config.Offline.Bool()))
	print.Entry("run pre-push hook", format.Bool(bool(config.PushHook)))
	print.Entry("push new branches", format.Bool(config.ShouldNewBranchPush()))
	print.Entry("ship deletes the tracking branch", format.Bool(config.ShipDeleteTrackingBranch.Bool()))
	print.Entry("sync-feature strategy", config.SyncFeatureStrategy.String())
	print.Entry("sync-perennial strategy", config.SyncPerennialStrategy.String())
	print.Entry("sync with upstream", format.Bool(config.SyncUpstream.Bool()))
	print.Entry("sync before shipping", format.Bool(config.SyncBeforeShip.Bool()))
	fmt.Println()
	print.Header("Hosting")
	print.Entry("hosting platform override", format.StringSetting(config.HostingPlatform.String()))
	print.Entry("GitHub token", format.StringSetting(string(config.GitHubToken)))
	print.Entry("GitLab token", format.StringSetting(string(config.GitLabToken)))
	print.Entry("Gitea token", format.StringSetting(string(config.GiteaToken)))
	fmt.Println()
	if !config.MainBranch.IsEmpty() {
		print.LabelAndValue("Branch Lineage", format.BranchLineage(config.Lineage))
	}
}
