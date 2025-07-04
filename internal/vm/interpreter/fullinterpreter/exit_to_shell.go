package fullinterpreter

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runlog"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// exitToShell is called when Git Town should exit to the shell without an error
func exitToShell(args ExecuteArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	if err = runlog.Write(runlog.EventEnd, endBranchesSnapshot.Branches, Some(args.RunState.Command), args.RootDir); err != nil {
		return err
	}
	args.RunState.EndBranchesSnapshot = Some(endBranchesSnapshot)
	configGitAccess := gitconfig.IO{Shell: args.Backend}
	globalSnapshot, err := configGitAccess.Load(Some(configdomain.ConfigScopeGlobal), configdomain.UpdateOutdatedNo)
	if err != nil {
		return err
	}
	localSnapshot, err := configGitAccess.Load(Some(configdomain.ConfigScopeLocal), configdomain.UpdateOutdatedNo)
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
	if err = args.RunState.MarkAsUnfinished(args.Git, args.Backend, true); err != nil {
		return err
	}
	if err = runstate.Save(args.RunState, args.RootDir); err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	args.FinalMessages.Add(`Run "git town continue" to go to the next branch.`)
	print.Footer(args.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	return nil
}
