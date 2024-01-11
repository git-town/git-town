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
	debugCommand.AddCommand(selectCommitAuthorCmd())
	debugCommand.AddCommand(unfinishedStateCommitAuthorCmd())
	return debugCommand
}
