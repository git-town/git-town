package cmd

import (
	"github.com/spf13/cobra"
)

func installCommand(rootCmd *cobra.Command) *cobra.Command {
	installCmd := cobra.Command{
		Use:   "install",
		Short: "Commands to set up Git Town on your computer",
		Args:  cobra.NoArgs,
	}
	installCmd.AddCommand(aliasCommand())
	installCmd.AddCommand(completionsCmd(rootCmd))
	return &installCmd
}
