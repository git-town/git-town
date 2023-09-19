package runvm

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/undo"
)

// Execute runs the commands in the given runstate.
func Execute(args ExecuteArgs) error {
	for {
		step := args.RunState.RunStepList.Pop()
		if step == nil {
			return finished(args)
		}
		stepName := runstate.TypeName(step)
		if stepName == "SkipCurrentBranchSteps" {
			args.RunState.SkipCurrentBranchSteps()
			continue
		}
		if stepName == "PushBranchAfterCurrentBranchSteps" {
			err := args.RunState.AddPushBranchStepAfterCurrentBranchSteps(&args.Run.Backend)
			if err != nil {
				return err
			}
			continue
		}
		err := step.Run(steps.RunArgs{
			Runner:    args.Run,
			Connector: args.Connector,
			Lineage:   args.Lineage,
		})
		if err != nil {
			return errored(step, err, args)
		}
	}
}

// finished is called when executing all steps has successfully finished.
func finished(args ExecuteArgs) error {
	args.RunState.MarkAsFinished()
	undoSteps, err := createUndoList(createUndoListArgs{
		Run:                     args.Run,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
	})
	if err != nil {
		return err
	}
	args.RunState.UndoStepList = undoSteps
	if args.RunState.IsAbort || args.RunState.IsUndo {
		err := persistence.Delete(args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateDeleteProblem, err)
		}
	} else {
		err := persistence.Save(args.RunState, args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateSaveProblem, err)
		}
	}
	fmt.Println()
	args.Run.Stats.PrintAnalysis()
	return nil
}

// errored is called when the given step has resulted in the given error.
func errored(step steps.Step, runErr error, args ExecuteArgs) error {
	args.RunState.AbortStepList.Append(step.CreateAbortSteps()...)
	undoSteps, err := createUndoList(createUndoListArgs{
		Run:                     args.Run,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
	})
	if err != nil {
		return err
	}
	args.RunState.UndoStepList = undoSteps
	if step.ShouldAutomaticallyAbortOnError() {
		return autoAbort(step, runErr, args)
	}
	args.RunState.RunStepList.Prepend(step.CreateContinueSteps()...)
	err = args.RunState.MarkAsUnfinished(&args.Run.Backend)
	if err != nil {
		return err
	}
	currentBranch, err := args.Run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	rebasing, err := args.Run.Backend.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if args.RunState.Command == "sync" && !(rebasing && args.Run.Config.IsMainBranch(currentBranch)) {
		args.RunState.UnfinishedDetails.CanSkip = true
	}
	err = persistence.Save(args.RunState, args.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateSaveProblem, err)
	}
	message := runErr.Error() + messages.AbortContinueGuidance
	if args.RunState.UnfinishedDetails.CanSkip {
		message += messages.ContinueSkipGuidance
	}
	message += "\n"
	return fmt.Errorf(message)
}

// autoAbort is called when a step that produced an error triggers an auto-abort.
func autoAbort(step steps.Step, runErr error, args ExecuteArgs) error {
	cli.PrintError(fmt.Errorf(messages.RunAutoAborting, runErr.Error()))
	abortRunState := args.RunState.CreateAbortRunState()
	err := Execute(ExecuteArgs{
		RunState:                &abortRunState,
		Run:                     args.Run,
		Connector:               args.Connector,
		RootDir:                 args.RootDir,
		Lineage:                 args.Lineage,
		InitialBranchesSnapshot: args.InitialBranchesSnapshot,
		InitialConfigSnapshot:   args.InitialConfigSnapshot,
	})
	if err != nil {
		return fmt.Errorf(messages.RunstateAbortStepProblem, err)
	}
	return step.CreateAutomaticAbortError()
}

type ExecuteArgs struct {
	RunState                *runstate.RunState
	Run                     *git.ProdRunner
	Connector               hosting.Connector
	RootDir                 string
	InitialBranchesSnapshot undo.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
	Lineage                 config.Lineage
}

func createUndoList(args createUndoListArgs) (runstate.StepList, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return runstate.StepList{}, errors.New(messages.DirCurrentProblem)
	}
	finalConfigSnapshot := undo.ConfigSnapshot{
		Cwd:       currentDirectory,
		GitConfig: config.LoadGitConfig(args.Run.Backend),
	}
	configDiff := finalConfigSnapshot.Diff(args.InitialConfigSnapshot)
	undoConfigSteps := configDiff.UndoSteps()
	allBranches, active, err := args.Run.Backend.BranchInfos()
	if err != nil {
		return runstate.StepList{}, err
	}
	finalBranchesSnapshot := undo.BranchesSnapshot{
		Branches: allBranches,
		Active:   active,
	}
	bba := args.InitialBranchesSnapshot.Changes(finalBranchesSnapshot)
	branchesDiff := bba.Diff()
	undoBranchesSteps := branchesDiff.Steps(args.Run.Config.Lineage(), args.Run.Config.BranchTypes(), args.InitialBranchesSnapshot.Active, finalBranchesSnapshot.Active)
	undoConfigSteps.AppendList(undoBranchesSteps)
	return undoConfigSteps, nil
}

type createUndoListArgs struct {
	Run                     *git.ProdRunner
	InitialBranchesSnapshot undo.BranchesSnapshot
	InitialConfigSnapshot   undo.ConfigSnapshot
}
