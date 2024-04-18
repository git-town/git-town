package config

import (
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/git-town/git-town/v14/src/vm/statefile"
)

// Finished is called when a Git Town command that only changes configuration has finished successfully.
func Finished(args FinishedArgs) error {
	// TODO: extract the code to load a config snapshot into a reusable function
	//       since it exists in multiple places
	configGitAccess := gitconfig.Access{Runner: args.Runner.Backend.Runner}
	globalSnapshot, _, err := configGitAccess.LoadGlobal(false)
	if err != nil {
		return err
	}
	localSnapshot, _, err := configGitAccess.LoadLocal(false)
	if err != nil {
		return err
	}
	configSnapshot := undoconfig.ConfigSnapshot{
		Global: globalSnapshot,
		Local:  localSnapshot,
	}
	runState := runstate.RunState{
		AbortProgram:             program.Program{},
		BeginBranchesSnapshot:    gitdomain.EmptyBranchesSnapshot(),
		BeginConfigSnapshot:      args.BeginConfigSnapshot,
		BeginStashSize:           0,
		Command:                  args.Command,
		DryRun:                   false,
		EndBranchesSnapshot:      gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:        configSnapshot,
		EndStashSize:             0,
		FinalUndoProgram:         program.Program{},
		RunProgram:               program.Program{},
		UndoablePerennialCommits: gitdomain.SHAs{},
		UnfinishedDetails:        nil,
	}
	print.Footer(args.Verbose, args.Runner.CommandsCounter.Count(), args.Runner.FinalMessages.Result())
	return statefile.Save(&runState, args.RootDir)
}

type FinishedArgs struct {
	BeginConfigSnapshot undoconfig.ConfigSnapshot
	Command             string
	EndConfigSnapshot   undoconfig.ConfigSnapshot
	RootDir             gitdomain.RepoRootDir
	Runner              *git.ProdRunner
	Verbose             bool
}
