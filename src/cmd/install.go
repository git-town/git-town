package cmd

import (
	"github.com/spf13/cobra"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Set up Git Town on your computer",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println()
	// },
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func init() {
	RootCmd.AddCommand(installCommand)
}
