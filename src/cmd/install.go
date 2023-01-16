package cmd

import (
	"github.com/spf13/cobra"
)

func installCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Commands to set up Git Town on your computer",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(aliasCommand())
	cmd.AddCommand(completionsCmd())
	return cmd
}
