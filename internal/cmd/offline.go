package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	offlineDesc = "Display or set offline mode"
	offlineHelp = `
Git Town avoids network operations in offline mode.`
)

func offlineCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "offline [(yes | no)]",
		Args:    cobra.MaximumNArgs(1),
		GroupID: cmdhelpers.GroupIDSetup,
		Short:   offlineDesc,
		Long:    cmdhelpers.Long(offlineDesc, offlineHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				DryRun:  None[configdomain.DryRun](),
				Verbose: verbose,
			})
			return executeOffline(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeOffline(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  false,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	switch len(args) {
	case 0:
		displayOfflineStatus(repo.UnvalidatedConfig)
	case 1:
		err = setOfflineStatus(args[0], repo.Frontend)
		if err != nil {
			return err
		}
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "offline",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       []gitdomain.BranchName{},
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
	})
}

func displayOfflineStatus(config config.UnvalidatedConfig) {
	fmt.Println(format.Bool(config.NormalConfig.Offline.IsOffline()))
}

func setOfflineStatus(text string, runner subshelldomain.Runner) error {
	value, err := gohacks.ParseBool(text, "offline status")
	if err != nil {
		return fmt.Errorf(messages.ValueInvalid, configdomain.KeyOffline, text)
	}
	return gitconfig.SetOffline(runner, configdomain.Offline(value))
	// in the future, we could remove the offline setting here
}
