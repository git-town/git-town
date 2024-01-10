package debug

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	debugCommand := &cobra.Command{
		Use:    "debug",
		Short:  "Displays dialogs to help develop them.",
		Hidden: true,
	}
	debugCommand.AddCommand(enterMainBranchCmd())
	debugCommand.AddCommand(enterParentCmd())
	debugCommand.AddCommand(enterPerennialBranchesCmd())
	return debugCommand
}
