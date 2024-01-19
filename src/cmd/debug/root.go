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
	debugCommand.AddCommand(enterMainBranchCmd())
	debugCommand.AddCommand(enterParentCmd())
	debugCommand.AddCommand(enterPerennialBranchesCmd())
	debugCommand.AddCommand(enterSyncFeatureStrategy())
	debugCommand.AddCommand(enterSyncPerennialStrategy())
	debugCommand.AddCommand(enterPushHookCmd())
	debugCommand.AddCommand(enterPushNewBranches())
	debugCommand.AddCommand(enterShipDeleteTrackingBranch())
	debugCommand.AddCommand(enterSyncBeforeShip())
	debugCommand.AddCommand(selectCommitAuthorCmd())
	debugCommand.AddCommand(unfinishedStateCommitAuthorCmd())
	return debugCommand
}
