package cmd

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

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

func init() {
	newBranchPushFlagCommand.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets your global new branch push flag")
	RootCmd.AddCommand(newBranchPushFlagCommand)
}
