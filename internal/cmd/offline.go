package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cli/format"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/execute"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks"
	"github.com/git-town/git-town/v17/internal/messages"
	configInterpreter "github.com/git-town/git-town/v17/internal/vm/interpreter/config"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/spf13/cobra"
)

const offlineDesc = "Display or set offline mode"

const offlineHelp = `
Git Town avoids network operations in offline mode.`

func offlineCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "offline [(yes | no)]",
		Args:    cobra.MaximumNArgs(1),
		GroupID: "setup",
		Short:   offlineDesc,
		Long:    cmdhelpers.Long(offlineDesc, offlineHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeOffline(args, verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeOffline(args []string, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  false,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	switch len(args) {
	case 0:
		displayOfflineStatus(repo.UnvalidatedConfig)
	case 1:
		err = setOfflineStatus(args[0], repo.UnvalidatedConfig)
		if err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "offline",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       []gitdomain.BranchName{},
		Verbose:               verbose,
	})
}

func displayOfflineStatus(config config.UnvalidatedConfig) {
	fmt.Println(format.Bool(config.NormalConfig.Offline.IsTrue()))
}

func setOfflineStatus(text string, config config.UnvalidatedConfig) error {
	value, err := gohacks.ParseBool(text, "offline status")
	if err != nil {
		return fmt.Errorf(messages.ValueInvalid, configdomain.KeyOffline, text)
	}
	if offline, has := value.Get(); has {
		return config.NormalConfig.SetOffline(configdomain.Offline(offline))
	}
	// in the future, we could remove the offline setting here
	return nil
}
