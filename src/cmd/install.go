package cmd

import (
	"github.com/spf13/cobra"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Set up Git Town on your computer",
	Args:  cobra.NoArgs,
}

func init() {
	RootCmd.AddCommand(installCommand)
}
