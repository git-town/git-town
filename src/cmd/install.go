package cmd

import (
	"github.com/spf13/cobra"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Commands to set up Git Town on your computer",
	Args:  cobra.NoArgs,
}

func init() {
	RootCmd.AddCommand(installCommand)
}
