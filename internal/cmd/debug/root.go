// Package debug implements Git Town's hidden "debug" command.
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
	debugCommand.AddCommand(enterBitbucketAppPassword())
	debugCommand.AddCommand(enterBitbucketUsername())
	debugCommand.AddCommand(enterNewBranchType())
	debugCommand.AddCommand(enterDefaultBranch())
	debugCommand.AddCommand(enterDevRemote())
	debugCommand.AddCommand(enterFeatureRegex())
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
	debugCommand.AddCommand(enterSyncPrototypeStrategy())
	debugCommand.AddCommand(enterSyncUpstream())
	debugCommand.AddCommand(enterSyncTags())
	debugCommand.AddCommand(enterPushHookCmd())
	debugCommand.AddCommand(enterPushNewBranches())
	debugCommand.AddCommand(enterShipDeleteTrackingBranch())
	debugCommand.AddCommand(enterShipStrategy())
	debugCommand.AddCommand(selectCommitAuthorCmd())
	debugCommand.AddCommand(switchBranch())
	debugCommand.AddCommand(unfinishedStateCommitAuthorCmd())
	debugCommand.AddCommand(welcome())
	return debugCommand
}
