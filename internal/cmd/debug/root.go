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
	debugCommand.AddCommand(enterCommitsToBeam())
	debugCommand.AddCommand(enterDevRemote())
	debugCommand.AddCommand(enterFeatureRegex())
	debugCommand.AddCommand(enterForgeType())
	debugCommand.AddCommand(enterGiteaToken())
	debugCommand.AddCommand(enterGitHubToken())
	debugCommand.AddCommand(enterGitLabToken())
	debugCommand.AddCommand(enterMainBranchCmd())
	debugCommand.AddCommand(enterNewBranchType())
	debugCommand.AddCommand(enterOriginHostname())
	debugCommand.AddCommand(enterPerennialBranches())
	debugCommand.AddCommand(enterPerennialRegex())
	debugCommand.AddCommand(enterSyncFeatureStrategy())
	debugCommand.AddCommand(enterSyncPerennialStrategy())
	debugCommand.AddCommand(enterSyncPrototypeStrategy())
	debugCommand.AddCommand(enterSyncUpstream())
	debugCommand.AddCommand(enterSyncTags())
	debugCommand.AddCommand(enterTokenScope())
	debugCommand.AddCommand(enterPushHookCmd())
	debugCommand.AddCommand(enterShareNewBranches())
	debugCommand.AddCommand(enterShipDeleteTrackingBranch())
	debugCommand.AddCommand(enterShipStrategy())
	debugCommand.AddCommand(enterUnknownBranch())
	debugCommand.AddCommand(selectCommitAuthorCmd())
	debugCommand.AddCommand(switchBranch())
	debugCommand.AddCommand(unfinishedStateCommitAuthorCmd())
	debugCommand.AddCommand(welcome())
	return debugCommand
}
