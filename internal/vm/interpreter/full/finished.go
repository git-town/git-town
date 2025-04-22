package interpreter

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/print"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/config/gitconfig"
	"github.com/git-town/git-town/v19/internal/git"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/gohacks"
	"github.com/git-town/git-town/v19/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/undo/undoconfig"
	"github.com/git-town/git-town/v19/internal/vm/runstate"
	"github.com/git-town/git-town/v19/internal/vm/statefile"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// finished is called when executing all steps has successfully finished.
func finished(args finishedArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	configGitAccess := gitconfig.Access{Runner: args.Backend}
	globalSnapshot, _, err := configGitAccess.LoadGlobal(false)
	if err != nil {
		return err
	}
	localSnapshot, _, err := configGitAccess.LoadLocal(false)
	if err != nil {
		return err
	}
	args.RunState.EndConfigSnapshot = Some(undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	})
	endStashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return err
	}
	args.RunState.EndStashSize = Some(endStashSize)
	args.RunState.MarkAsFinished(endBranchesSnapshot)
	err = statefile.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	return nil
}

type finishedArgs struct {
	Backend         gitdomain.RunnerQuerier
	CommandsCounter Mutable[gohacks.Counter]
	FinalMessages   stringslice.Collector
	Git             git.Commands
	RootDir         gitdomain.RepoRootDir
	RunState        runstate.RunState
	Verbose         configdomain.Verbose
}
