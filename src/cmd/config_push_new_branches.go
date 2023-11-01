package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/cli/format"
	"github.com/git-town/git-town/v10/src/cli/io"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/git"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/spf13/cobra"
)

const pushNewBranchesDesc = "Displays or changes whether new branches get pushed to origin"

const pushNewBranchesHelp = `
If "push-new-branches" is true, the Git Town commands hack, append, and prepend
push the new branch to the origin remote.`

func pushNewBranchesCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "If set, reads or updates the new branch push strategy for all repositories on this machine", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   "push-new-branches [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: pushNewBranchesDesc,
		Long:  long(pushNewBranchesDesc, pushNewBranchesHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigPushNewBranches(args, readGlobalFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func executeConfigPushNewBranches(args []string, global, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  false,
	})
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return setPushNewBranches(args[0], global, &repo.Runner)
	}
	return printPushNewBranches(global, &repo.Runner)
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
	io.Println(format.Bool(setting))
	return nil
}

func setPushNewBranches(text string, globalFlag bool, run *git.ProdRunner) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(messages.InputYesOrNo, text)
	}
	return run.Config.SetNewBranchPush(value, globalFlag)
}
