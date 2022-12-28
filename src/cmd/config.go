package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
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
		if prodRepo.Config.MainBranch() != "" {
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

var mainBranchConfigCommand = &cobra.Command{
	Use:   "main-branch [<branch>]",
	Short: "Displays or sets your main development branch",
	Long: `Displays or sets your main development branch

The main branch is the Git branch from which new feature branches are cut.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printMainBranch()
		} else {
			err := setMainBranch(args[0], prodRepo)
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func printMainBranch() {
	cli.Println(cli.PrintableMainBranch(prodRepo.Config.MainBranch()))
}

func setMainBranch(branchName string, repo *git.ProdRepo) error {
	hasBranch, err := repo.Silent.HasLocalBranch(branchName)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf("there is no branch named %q", branchName)
	}
	return repo.Config.SetMainBranch(branchName)
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
	configCommand.AddCommand(mainBranchConfigCommand)
	configCommand.AddCommand(resetConfigCommand)
	configCommand.AddCommand(setupConfigCommand)
	RootCmd.AddCommand(configCommand)
}
