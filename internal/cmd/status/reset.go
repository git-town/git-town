package status

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/flags"
	"github.com/git-town/git-town/v19/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/git"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/gohacks"
	"github.com/git-town/git-town/v19/internal/gohacks/cache"
	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/subshell"
	"github.com/git-town/git-town/v19/internal/vm/statefile"
	. "github.com/git-town/git-town/v19/pkg/prelude"
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
			return executeStatusReset(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeStatusReset(verbose configdomain.Verbose) error {
	commandsCounter := NewMutable(new(gohacks.Counter))
	backendRunner := subshell.BackendRunner{
		Dir:             None[string](),
		CommandsCounter: commandsCounter,
		Verbose:         verbose,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.WithPrevious[gitdomain.LocalBranchName]{},
		RemotesCache:       &cache.Cache[gitdomain.Remotes]{},
	}
	rootDir, hasRootDir := gitCommands.RootDirectory(backendRunner).Get()
	if !hasRootDir {
		return errors.New(messages.RepoOutside)
	}
	existed, err := statefile.Delete(rootDir)
	if err != nil {
		return err
	}
	if existed {
		fmt.Println(messages.RunstateDeleted)
	} else {
		fmt.Println(messages.RunstateDoesntExist)
	}
	return nil
}
