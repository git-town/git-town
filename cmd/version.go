package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of the currently installed Git Town executable",
	Long:  `Displays the version of the currently installed Git Town executable`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Git Town version 2.2.0")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
