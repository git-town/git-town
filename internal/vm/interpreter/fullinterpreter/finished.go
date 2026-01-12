package fullinterpreter

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runlog"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type finishedArgs struct {
	Backend         subshelldomain.RunnerQuerier
	CommandsCounter Mutable[gohacks.Counter]
	ConfigDir       configdomain.ConfigDirRepo
	FinalMessages   stringslice.Collector
	Git             git.Commands
	Inputs          dialogcomponents.Inputs
	RunState        runstate.RunState
	Verbose         configdomain.Verbose
}

// finished is called when executing all steps has successfully finished.
func finished(args finishedArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	runlogPath := runlog.NewRunlogPath(args.ConfigDir)
	if err = runlog.Write(runlog.EventEnd, endBranchesSnapshot.Branches, Some(args.RunState.Command), runlogPath); err != nil {
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
	args.RunState.EndConfigSnapshot = Some(configdomain.EndConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	})
	endStashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return err
	}
	args.RunState.EndStashSize = Some(endStashSize)
	args.RunState.MarkAsFinished(endBranchesSnapshot)
	runstatePath := runstate.NewRunstatePath(args.ConfigDir)
	if err = runstate.Save(args.RunState, runstatePath); err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	print.Footer(args.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	args.Inputs.VerifyAllUsed()
	return nil
}
