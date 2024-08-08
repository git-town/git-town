package debug

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	debugCommand := &cobra.Command{
		Use:    "debug",
		Short:  "Displays dialogs to help debug them.",
		Hidden: true,
	}
	debugCommand.AddCommand(enterAliases())
	debugCommand.AddCommand(enterCreatePrototypeBranches())
	debugCommand.AddCommand(enterHostingPlatform())
	debugCommand.AddCommand(enterGiteaToken())
	debugCommand.AddCommand(enterGitHubToken())
	debugCommand.AddCommand(enterGitLabToken())
	debugCommand.AddCommand(enterMainBranchCmd())
	debugCommand.AddCommand(enterParentCmd())
	debugCommand.AddCommand(enterOriginHostname())
	debugCommand.AddCommand(enterPerennialBranches())
	debugCommand.AddCommand(enterPerennialRegex())
	debugCommand.AddCommand(enterSyncFeatureStrategy())
	debugCommand.AddCommand(enterSyncPerennialStrategy())
	debugCommand.AddCommand(enterSyncUpstream())
	debugCommand.AddCommand(enterSyncTags())
	debugCommand.AddCommand(enterPushHookCmd())
	debugCommand.AddCommand(enterPushNewBranches())
	debugCommand.AddCommand(enterShipDeleteTrackingBranch())
	debugCommand.AddCommand(selectCommitAuthorCmd())
	debugCommand.AddCommand(switchBranch())
	debugCommand.AddCommand(unfinishedStateCommitAuthorCmd())
	debugCommand.AddCommand(welcome())
	return debugCommand
}
