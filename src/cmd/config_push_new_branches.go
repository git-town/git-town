package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func pushNewBranchesCommand() *cobra.Command {
	globalFlag := false
	debug := false
	cmd := cobra.Command{
		Use:   "push-new-branches [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Displays or changes whether new branches get pushed to origin",
		Long: `Displays or changes whether new branches get pushed to origin.

If "push-new-branches" is true, the Git Town commands hack, append, and prepend
push the new branch to the origin remote.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigurePushNewBranches(debug, globalFlag, args)
		},
	}
	debugFlag(&cmd, &debug)
	cmd.Flags().BoolVar(&globalFlag, "global", false, "Displays or sets your global new branch push flag")
	return &cmd
}

func runConfigurePushNewBranches(debug, global bool, args []string) error {
	repo := Repo(debug, false)
	if err := ensure(&repo, hasGitVersion, isRepository); err != nil {
		return err
	}
	if len(args) > 0 {
		return setPushNewBranches(args[0], global, &repo)
	}
	return printPushNewBranches(global, &repo)
}

func printPushNewBranches(globalFlag bool, repo *git.PublicRepo) error {
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

func setPushNewBranches(text string, globalFlag bool, repo *git.PublicRepo) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	return repo.Config.SetNewBranchPush(value, globalFlag)
}
