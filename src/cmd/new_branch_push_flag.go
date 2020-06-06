package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/src/command"
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
				fmt.Println("Please provide either true or false")
				os.Exit(1)
			}
			setNewBranchPushFlag(value)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return validateBooleanArgument(args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

func printNewBranchPushFlag() {
	if globalFlag {
		command.Println(strconv.FormatBool(git.Config().ShouldNewBranchPushGlobal()))
	} else {
		command.Println(git.GetPrintableNewBranchPushFlag())
	}
}

func setNewBranchPushFlag(value bool) {
	git.Config().SetNewBranchPush(value, globalFlag)
}

func init() {
	newBranchPushFlagCommand.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets your global new branch push flag")
	RootCmd.AddCommand(newBranchPushFlagCommand)
}
