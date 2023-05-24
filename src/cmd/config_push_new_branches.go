package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/spf13/cobra"
)

const pushNewBranchesDesc = "Displays or changes whether new branches get pushed to origin"

const pushNewBranchesHelp = `
If "push-new-branches" is true, the Git Town commands hack, append, and prepend
push the new branch to the origin remote.`

func pushNewBranchesCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "If set, reads or updates the new branch push strategy for all repositories on this machine")
	cmd := cobra.Command{
		Use:   "push-new-branches [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: pushNewBranchesDesc,
		Long:  long(pushNewBranchesDesc, pushNewBranchesHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pushNewBranches(args, readGlobalFlag(cmd), readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func pushNewBranches(args []string, global, debug bool) error {
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		OmitBranchNames:       true,
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	if len(args) > 0 {
		return setPushNewBranches(args[0], global, &run)
	}
	return printPushNewBranches(global, &run)
}

func printPushNewBranches(globalFlag bool, run *git.ProdRunner) error {
	var setting bool
	var err error
	if globalFlag {
		setting, err = run.Config.ShouldNewBranchPushGlobal()
	} else {
		setting, err = run.Config.ShouldNewBranchPush()
	}
	if err != nil {
		return err
	}
	cli.Println(cli.FormatBool(setting))
	return nil
}

func setPushNewBranches(text string, globalFlag bool, run *git.ProdRunner) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no"`, text)
	}
	return run.Config.SetNewBranchPush(value, globalFlag)
}
