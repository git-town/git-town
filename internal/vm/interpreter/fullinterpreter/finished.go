package fullinterpreter

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runlog"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// finished is called when executing all steps has successfully finished.
func finished(args finishedArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	if err = runlog.Write(runlog.EventEnd, endBranchesSnapshot.Branches, Some(args.RunState.Command), args.RootDir); err != nil {
		return err
	}
	globalSnapshot, err := gitconfig.LoadSnapshot(args.Backend, Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedNo)
	if err != nil {
		return err
	}
	localSnapshot, err := gitconfig.LoadSnapshot(args.Backend, Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)
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
	if err = runstate.Save(args.RunState, args.RootDir); err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	args.Inputs.VerifyAllUsed()
	return nil
}

type finishedArgs struct {
	Backend         subshelldomain.RunnerQuerier
	CommandsCounter Mutable[gohacks.Counter]
	FinalMessages   stringslice.Collector
	Git             git.Commands
	Inputs          dialogcomponents.Inputs
	RootDir         gitdomain.RepoRootDir
	RunState        runstate.RunState
	Verbose         configdomain.Verbose
}
