package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/userinput"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		cli.PrintLabelAndValue("Main branch", cli.PrintableMainBranch(prodRepo.Config.MainBranch()))
		cli.PrintLabelAndValue("Perennial branches", cli.PrintablePerennialBranches(prodRepo.Config.PerennialBranches()))
		mainBranch := prodRepo.Config.MainBranch()
		if mainBranch != "" {
			cli.PrintLabelAndValue("Branch Ancestry", cli.PrintableBranchAncestry(&prodRepo.Config))
		}
		cli.PrintLabelAndValue("Pull branch strategy", prodRepo.Config.PullBranchStrategy())
		cli.PrintLabelAndValue("New Branch Push Flag", cli.PrintableNewBranchPushFlag(prodRepo.Config.ShouldNewBranchPush()))
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

var resetConfigCommand = &cobra.Command{
	Use:   "reset",
	Short: "Resets your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := prodRepo.Config.RemoveLocalGitConfiguration()
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

var setupConfigCommand = &cobra.Command{
	Use:   "setup",
	Short: "Prompts to setup your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := userinput.ConfigureMainBranch(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		err = userinput.ConfigurePerennialBranches(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func init() {
	configCommand.AddCommand(resetConfigCommand)
	configCommand.AddCommand(setupConfigCommand)
	RootCmd.AddCommand(configCommand)
}
