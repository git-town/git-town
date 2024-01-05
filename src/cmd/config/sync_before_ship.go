package config

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/spf13/cobra"
)

const syncBeforeShipDesc = `Displays or changes whether "git ship" syncs the branch it ships`

const syncBeforeShipHelp = `
If "sync-before-ship" is enabled, the "git ship" command
executes "git sync" before shipping a branch.
This allows you to deal with breakage from resolving merge conflicts
on the feature branch instead of the main branch.
The downside is that this will trigger another CI run,
which might prevent the ship until it is finished.`

func syncBeforeShipCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "If set, reads or updates the sync-before-ship strategy for all repositories on this machine", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   "sync-before-ship [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: syncBeforeShipDesc,
		Long:  cmdhelpers.Long(syncBeforeShipDesc, syncBeforeShipHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSyncBeforeShipanches(args, readGlobalFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func executeSyncBeforeShipanches(args []string, global, verbose bool) error {
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
		return setSyncBeforeShip(args[0], global, repo.Runner)
	}
	return printSyncBeforeShip(global, repo.Runner)
}

func printSyncBeforeShip(globalFlag bool, run *git.ProdRunner) error {
	var setting *configdomain.SyncBeforeShip
	if globalFlag {
		setting = run.Config.GlobalGitConfig.SyncBeforeShip
		if setting == nil {
			defaults := configdomain.DefaultConfig()
			setting = &defaults.SyncBeforeShip
		}
	} else {
		setting = &run.Config.SyncBeforeShip
	}
	io.Println(format.Bool(setting.Bool()))
	return nil
}

func setSyncBeforeShip(text string, globalFlag bool, run *git.ProdRunner) error {
	boolValue, err := gohacks.ParseBool(text)
	if err != nil {
		return fmt.Errorf(messages.InputYesOrNo, text)
	}
	return run.Config.SetSyncBeforeShip(configdomain.SyncBeforeShip(boolValue), globalFlag)
}
