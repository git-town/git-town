package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Git Town 4.1.2")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
