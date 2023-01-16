package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pushNewBranchesCommand() *cobra.Command {
	globalFlag := false
	cmd := cobra.Command{
		Use:   "push-new-branches [--global] [(yes | no)]",
		Short: "Displays or changes whether new branches get pushed to origin",
		Long: `Displays or changes whether new branches get pushed to origin.

If "push-new-branches" is true, the Git Town commands hack, append, and prepend
push the new branch to the origin remote.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				err := printPushNewBranches(globalFlag, prodRepo)
				if err != nil {
					cli.Exit(err)
				}
			} else {
				err := setPushNewBranches(args[0], globalFlag, prodRepo)
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
	cmd.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets your global new branch push flag")
	return &cmd
}

func printPushNewBranches(globalFlag bool, repo *git.ProdRepo) error {
	var setting bool
	var err error
	if globalFlag {
		setting, err = repo.Config.ShouldNewBranchPushGlobal()
	} else {
		setting, err = repo.Config.ShouldNewBranchPush()
	}
	if err != nil {
		return err
	}
	cli.Println(cli.FormatBool(setting))
	return nil
}

func setPushNewBranches(text string, globalFlag bool, repo *git.ProdRepo) error {
	value, err := cli.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	return repo.Config.SetNewBranchPush(value, globalFlag)
}
