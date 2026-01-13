package status

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/cache"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runlog"
	"github.com/git-town/git-town/v22/internal/state/runstate"
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
	configDirUser, err := configdomain.SystemUserConfigDir()
	if err != nil {
		return err
	}
	configDirRepo := configDirUser.RepoConfigDir(rootDir)

	// delete runstate
	runstatePath := runstate.NewRunstatePath(configDirRepo)
	err = os.Remove(runstatePath.String())
	switch {
	case os.IsNotExist(err):
		fmt.Println(messages.RunstateDoesntExist)
	case err != nil:
		return err
	default:
		fmt.Println(messages.RunstateDeleted)
	}

	// delete runlog
	runlogPath := runlog.NewRunlogPath(configDirRepo)
	err = os.Remove(runlogPath.String())
	switch {
	case os.IsNotExist(err):
		fmt.Println(messages.RunLogDoesntExist)
	case err != nil:
		return err
	default:
		fmt.Println(messages.RunLogDeleted)
	}

	return nil
}
