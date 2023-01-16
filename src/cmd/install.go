package cmd

import (
	"github.com/spf13/cobra"
)

func installCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Commands to set up Git Town on your computer",
		Args:  cobra.NoArgs,
	}
	installCmd.AddCommand(aliasCommand())
	installCmd.AddCommand(completionsCmd())
	return installCmd
}
