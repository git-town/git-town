package status

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v19/internal/cmd/cmdhelpers"
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
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetDesc,
		Long:  cmdhelpers.Long(statusResetDesc),
		RunE: func(_ *cobra.Command, _ []string) error {
			return executeStatusReset()
		},
	}
	return &cmd
}

func executeStatusReset() error {
	commandsCounter := NewMutable(new(gohacks.Counter))
	backendRunner := subshell.BackendRunner{
		Dir:             None[string](),
		CommandsCounter: commandsCounter,
		Verbose:         false,
	}
	gitCommands := git.Commands{
		CurrentBranchCache: &cache.WithPrevious[gitdomain.LocalBranchName]{},
		RemotesCache:       &cache.Cache[gitdomain.Remotes]{},
	}
	rootDir, hasRootDir := gitCommands.RootDirectory(backendRunner).Get()
	if !hasRootDir {
		return errors.New(messages.RepoOutside)
	}
	err := statefile.Delete(rootDir)
	if err != nil {
		return err
	}
	fmt.Println(messages.RunstateDeleted)
	return nil
}
