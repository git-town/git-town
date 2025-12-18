package status

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/cache"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state"
	"github.com/git-town/git-town/v22/internal/subshell"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

func resetRunstateCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetDesc,
		Long:  cmdhelpers.Long(statusResetDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeStatusReset(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeStatusReset(cliConfig configdomain.PartialConfig) error {
	commandsCounter := NewMutable(new(gohacks.Counter))
	backendRunner := subshell.BackendRunner{
		Dir:             None[string](),
		CommandsCounter: commandsCounter,
		Verbose:         cliConfig.Verbose.GetOr(false),
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.WithPrevious[gitdomain.LocalBranchName]{},
		RemotesCache:       &cache.Cache[gitdomain.Remotes]{},
	}
	rootDir, hasRootDir := gitCommands.RootDirectory(backendRunner).Get()
	if !hasRootDir {
		return errors.New(messages.RepoOutside)
	}
	runstateExisted, err := state.Delete(rootDir, state.FileTypeRunstate)
	if err != nil {
		return err
	}
	if runstateExisted {
		fmt.Println(messages.RunstateDeleted)
	} else {
		fmt.Println(messages.RunstateDoesntExist)
	}
	runlogExisted, err := state.Delete(rootDir, state.FileTypeRunlog)
	if err != nil {
		return err
	}
	if runlogExisted {
		fmt.Println(messages.RunLogDeleted)
	} else {
		fmt.Println(messages.RunLogDoesntExist)
	}
	return nil
}
