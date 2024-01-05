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

const syncUpstreamDesc = `Displays or changes whether "git sync" pulls updates from the "upstream" remote`

const syncUpstreamHelp = `
If "sync-upstream" is enabled, and your Git repository has an "upstream" remote, "git sync" will also pull updates from the main branch at that upstream remote.`

func syncUpstreamCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addGlobalFlag, readGlobalFlag := flags.Bool("global", "g", "If set, reads or updates the sync-upstream strategy for all repositories on this machine", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   "sync-upstream [--global] [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: syncUpstreamDesc,
		Long:  cmdhelpers.Long(syncUpstreamDesc, syncUpstreamHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSyncUpstream(args, readGlobalFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addGlobalFlag(&cmd)
	return &cmd
}

func executeSyncUpstream(args []string, global, verbose bool) error {
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
		return setSyncUpstream(args[0], global, repo.Runner)
	}
	return printSyncUpstream(global, repo.Runner)
}

func printSyncUpstream(globalFlag bool, run *git.ProdRunner) error {
	var setting *configdomain.SyncUpstream
	if globalFlag {
		setting = run.Config.GlobalGitConfig.SyncUpstream
		if setting == nil {
			defaults := configdomain.DefaultConfig()
			setting = &defaults.SyncUpstream
		}
	} else {
		setting = &run.Config.SyncUpstream
	}
	io.Println(format.Bool(setting.Bool()))
	return nil
}

func setSyncUpstream(text string, globalFlag bool, run *git.ProdRunner) error {
	boolValue, err := gohacks.ParseBool(text)
	if err != nil {
		return fmt.Errorf(messages.InputYesOrNo, text)
	}
	return run.Config.SetSyncUpstream(configdomain.SyncUpstream(boolValue), globalFlag)
}
