package undo

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/shared"
	"github.com/git-town/git-town/v12/src/vm/statefile"
)

// undoes the persisted runstate
func Execute(args ExecuteArgs) error {
	if args.RunState.DryRun {
		return nil
	}
	program := createProgram(args)

	// execute the undo program
	for _, opcode := range program {
		err := opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         args.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			Runner:                          args.Runner,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			fmt.Println(components.Red().Styled("NOTICE: " + err.Error()))
		}
	}

	err := statefile.Delete(args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateDeleteProblem, err)
	}
	print.Footer(args.Verbose, args.Runner.CommandsCounter.Count(), args.Runner.FinalMessages.Result())
	return nil
}

type ExecuteArgs struct {
	FullConfig       *configdomain.FullConfig
	HasOpenChanges   bool
	InitialStashSize gitdomain.StashSize
	Lineage          configdomain.Lineage
	RootDir          gitdomain.RepoRootDir
	RunState         runstate.RunState
	Runner           *git.ProdRunner
	Verbose          bool
}
