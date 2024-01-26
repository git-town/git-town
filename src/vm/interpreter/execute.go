package interpreter

import (
	"github.com/git-town/git-town/v11/src/cli/dialogs/dialogcomponents"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v11/src/undo/undoconfig"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/git-town/git-town/v11/src/vm/shared"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		nextStep := args.RunState.RunProgram.Pop()
		if nextStep == nil {
			return finished(args)
		}
		stepName := gohacks.TypeName(nextStep)
		if stepName == "SkipCurrentBranchProgram" {
			args.RunState.SkipCurrentBranchProgram()
			continue
		}
		err := nextStep.Run(shared.RunArgs{
			Connector:                       args.Connector,
			DialogTestInputs:                args.DialogTestInputs,
			Lineage:                         args.Lineage,
			PrependOpcodes:                  args.RunState.RunProgram.Prepend,
			RegisterUndoablePerennialCommit: args.RunState.RegisterUndoablePerennialCommit,
			Runner:                          args.Run,
			UpdateInitialBranchLocalSHA:     args.InitialBranchesSnapshot.Branches.UpdateLocalSHA,
		})
		if err != nil {
			return errored(nextStep, err, args)
		}
	}
}

type ExecuteArgs struct {
	*configdomain.FullConfig
	Connector               hostingdomain.Connector
	DialogTestInputs        *dialogcomponents.TestInputs
	InitialBranchesSnapshot gitdomain.BranchesStatus
	InitialConfigSnapshot   undoconfig.ConfigSnapshot
	InitialStashSnapshot    gitdomain.StashSize
	RootDir                 gitdomain.RepoRootDir
	Run                     *git.ProdRunner
	RunState                *runstate.RunState
	Verbose                 bool
}
