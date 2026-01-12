package fullinterpreter

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runlog"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// exitToShell is called when Git Town should exit to the shell without an error
func exitToShell(args ExecuteArgs) error {
	endBranchesSnapshot, err := args.Git.BranchesSnapshot(args.Backend)
	if err != nil {
		return err
	}
	if err = runlog.Write(runlog.EventEnd, endBranchesSnapshot.Branches, Some(args.RunState.Command), args.RunlogPath); err != nil {
		return err
	}
	args.RunState.EndBranchesSnapshot = Some(endBranchesSnapshot)
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
	if err = args.RunState.MarkAsUnfinished(args.Git, args.Backend, true); err != nil {
		return err
	}
	if err = runstate.Save(args.RunState, args.RunstatePath); err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	args.FinalMessages.Add(`Run "git town continue" to go to the next branch.`)
	print.Footer(args.Config.NormalConfig.Verbose, args.CommandsCounter.Immutable(), args.FinalMessages.Result())
	args.Inputs.VerifyAllUsed()
	return nil
}
