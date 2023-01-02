package cmd

import (
	"fmt"
	"strconv"

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
		cli.PrintLabelAndValue("Perennial branches", cli.PrintablePerennialBranches(prodRepo.Config.PerennialBranches.List()))
		if prodRepo.Config.MainBranch() != "" {
			cli.PrintLabelAndValue("Branch Ancestry", cli.PrintableBranchAncestry(prodRepo.Config.Ancestry))
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

var newBranchPushFlagCommand = &cobra.Command{
	Use:   "new-branch-push-flag [(true | false)]",
	Short: "Displays or sets your new branch push flag",
	Long: `Displays or sets your new branch push flag

If "new-branch-push-flag" is true, Git Town pushes branches created with
hack / append / prepend on creation. Defaults to false.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printNewBranchPushFlag(prodRepo)
		} else {
			value, err := strconv.ParseBool(args[0])
			if err != nil {
				cli.Exit(fmt.Errorf(`invalid argument: %q. Please provide either "true" or "false"`, args[0]))
			}
			err = setNewBranchPushFlag(value, prodRepo)
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

func printNewBranchPushFlag(repo *git.ProdRepo) {
	if globalFlag {
		cli.Println(strconv.FormatBool(repo.Config.ShouldNewBranchPushGlobal()))
	} else {
		cli.Println(cli.PrintableNewBranchPushFlag(prodRepo.Config.ShouldNewBranchPush()))
	}
}

func setNewBranchPushFlag(value bool, repo *git.ProdRepo) error {
	return repo.Config.SetNewBranchPush(value, globalFlag)
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

var offlineCommand = &cobra.Command{
	Use:   "offline [(true | false)]",
	Short: "Displays or sets offline mode",
	Long: `Displays or sets offline mode

Git Town avoids network operations in offline mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cli.Println(cli.PrintableOfflineFlag(prodRepo.Config.Offline.Enabled()))
		} else {
			value, err := strconv.ParseBool(args[0])
			if err != nil {
				cli.Exit(fmt.Errorf(`invalid argument: %q. Please provide either "true" or "false".\n`, args[0]))
			}
			err = prodRepo.Config.Offline.Enable(value)
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
}

var perennialBranchesCommand = &cobra.Command{
	Use:   "perennial-branches",
	Short: "Displays your perennial branches",
	Long: `Displays your perennial branches

Perennial branches are long-lived branches.
They cannot be shipped.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Println(cli.PrintablePerennialBranches(prodRepo.Config.PerennialBranches.List()))
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

var updatePrennialBranchesCommand = &cobra.Command{
	Use:   "update",
	Short: "Prompts to update your perennial branches",
	Long:  `Prompts to update your perennial branches`,
	Run: func(cmd *cobra.Command, args []string) {
		err := userinput.ConfigurePerennialBranches(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

var pullBranchStrategyCommand = &cobra.Command{
	Use:   "pull-branch-strategy [(rebase | merge)]",
	Short: "Displays or sets your pull branch strategy",
	Long: `Displays or sets your pull branch strategy

The pull branch strategy specifies what strategy to use
when merging remote tracking branches into local branches
for the main branch and perennial branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cli.Println(prodRepo.Config.PullBranchStrategy())
		} else {
			err := prodRepo.Config.SetPullBranchStrategy(args[0])
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] != "rebase" && args[0] != "merge" {
			return fmt.Errorf("invalid value: %q", args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
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

var syncStrategyCommand = &cobra.Command{
	Use:   "sync-strategy [(merge | rebase)]",
	Short: "Displays or sets your sync strategy",
	Long: `Displays or sets your sync strategy

The sync strategy specifies what strategy to use
when merging remote tracking branches into local feature branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cli.Println(prodRepo.Config.SyncStrategy())
		} else {
			err := prodRepo.Config.SetSyncStrategy(args[0])
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] != "merge" && args[0] != "rebase" {
			return fmt.Errorf("invalid value: %q", args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func init() {
	configCommand.AddCommand(mainBranchConfigCommand)
	newBranchPushFlagCommand.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets your global new branch push flag")
	configCommand.AddCommand(newBranchPushFlagCommand)
	configCommand.AddCommand(offlineCommand)
	perennialBranchesCommand.AddCommand(updatePrennialBranchesCommand)
	configCommand.AddCommand(perennialBranchesCommand)
	configCommand.AddCommand(pullBranchStrategyCommand)
	configCommand.AddCommand(resetConfigCommand)
	configCommand.AddCommand(setupConfigCommand)
	configCommand.AddCommand(syncStrategyCommand)
	RootCmd.AddCommand(configCommand)
}
