package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
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
			printNewBranchPushFlag()
		} else {
			value, err := strconv.ParseBool(args[0])
			if err != nil {
				fmt.Printf(`Error: invalid argument: %q. Please provide either "true" or "false".\n`, args[0])
				os.Exit(1)
			}
			setNewBranchPushFlag(value)
		}
	},
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func printNewBranchPushFlag() {
	if globalFlag {
		cli.Println(strconv.FormatBool(git.Config().ShouldNewBranchPushGlobal()))
	} else {
		cli.Println(cli.PrintableNewBranchPushFlag(prodRepo.ShouldNewBranchPush()))
	}
}

func setNewBranchPushFlag(value bool) {
	git.Config().SetNewBranchPush(value, globalFlag)
}

func init() {
	newBranchPushFlagCommand.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets your global new branch push flag")
	RootCmd.AddCommand(newBranchPushFlagCommand)
}
