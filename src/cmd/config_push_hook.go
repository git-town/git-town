package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/spf13/cobra"
)

const pushHookDesc = "Configures whether Git Town should run Git's pre-push hook."

const pushHookHelp = `
Enabled by default. When disabled, Git Town prevents Git's pre-push hook from running.`

func pushHookCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "If set, reads or updates the push hook flag for all repos on this machine", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   "push-hook [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: pushHookDesc,
		Long:  long(pushHookDesc, pushHookHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigPushHook(args, readGlobalFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func executeConfigPushHook(args []string, global, verbose bool) error {
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
		return setPushHook(args[0], global, repo.Runner)
	}
	return printPushHook(global, repo.Runner)
}

func printPushHook(globalFlag bool, run *git.ProdRunner) error {
	var setting configdomain.PushHook
	var err error
	if globalFlag {
		setting, err = run.GitTown.PushHookGlobal()
	} else {
		setting, err = run.GitTown.PushHook()
	}
	if err != nil {
		return err
	}
	io.Println(format.Bool(bool(setting)))
	return nil
}

func setPushHook(text string, global bool, run *git.ProdRunner) error {
	valueBool, err := gohacks.ParseBool(text)
	if err != nil {
		return fmt.Errorf(messages.InputYesOrNo, text)
	}
	value := configdomain.PushHook(valueBool)
	if global {
		return run.GitTown.SetPushHookGlobally(value)
	}
	return run.GitTown.SetPushHookLocally(value)
}
