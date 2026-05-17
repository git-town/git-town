package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/browser/browserdomain"
	"github.com/git-town/git-town/v23/internal/cli/flags"
	"github.com/git-town/git-town/v23/internal/cli/format"
	"github.com/git-town/git-town/v23/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v23/internal/config"
	"github.com/git-town/git-town/v23/internal/config/cliconfig"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/config/gitconfig"
	"github.com/git-town/git-town/v23/internal/execute"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	"github.com/git-town/git-town/v23/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v23/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v23/pkg/prelude"
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
		GroupID: cmdhelpers.GroupIDConfig,
		Short:   offlineDesc,
		Long:    cmdhelpers.Long(offlineDesc, offlineHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				BrowserEnabled:    None[browserdomain.BrowserEnabled](),
				BrowserExecutable: None[browserdomain.BrowserExecutable](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Interactive:       None[configdomain.Interactive](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
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
		IgnoreUnknown:    false,
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
		err = setOfflineStatus(stringss.TrimSpace(args[0]), repo.Frontend)
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
		ConfigDir:             repo.ConfigDir,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		TouchedBranches:       []gitdomain.BranchName{},
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
	})
}

func displayOfflineStatus(config config.UnvalidatedConfig) {
	fmt.Println(format.Bool(config.NormalConfig.Offline.IsOffline()))
}

func setOfflineStatus(text stringss.TrimmedString, runner subshelldomain.Runner) error {
	value, err := gohacks.ParseBool[configdomain.Offline](text, "offline status")
	if err != nil {
		return fmt.Errorf(messages.ValueInvalid, configdomain.KeyOffline, text)
	}
	return gitconfig.SetOffline(runner, value)
	// in the future, we could remove the offline setting here
}
