package debug

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	debugCommand := &cobra.Command{
		Use:    "debug <thing>",
		Hidden: true,
	}
	debugCommand.AddCommand(enterMainBranchCmd())
	debugCommand.AddCommand(enterParentCmd())
	return debugCommand
}
