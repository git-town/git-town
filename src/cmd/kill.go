package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/gohacks/slice"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/opcode"
	"github.com/git-town/git-town/v10/src/vm/program"
	"github.com/git-town/git-town/v10/src/vm/runstate"
	"github.com/spf13/cobra"
)

const killDesc = "Removes an obsolete feature branch"

const killHelp = `
Deletes the current or provided branch from the local and origin repositories.
Does not delete perennial branches nor the main branch.`

func killCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "kill [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: killDesc,
		Long:  long(killDesc, killHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeKill(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeKill(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineKillConfig(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	steps, finalUndoProgram := killProgram(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:             "kill",
		RunProgram:          steps,
		InitialActiveBranch: initialBranchesSnapshot.Active,
		FinalUndoProgram:    finalUndoProgram,
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              config.noPushHook,
	})
}

type killConfig struct {
	hasOpenChanges bool
	initialBranch  domain.LocalBranchName
	isOffline      bool
	lineage        config.Lineage
	mainBranch     domain.LocalBranchName
	noPushHook     bool
	previousBranch domain.LocalBranchName
	targetBranch   domain.BranchInfo
}

func determineKillConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (*killConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	targetBranchName := domain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	targetBranch := branches.All.FindByLocalName(targetBranchName)
	if targetBranch == nil {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	if targetBranch.IsLocal() {
		branches.Types, lineage, err = execute.EnsureKnownBranchAncestry(targetBranchName, execute.EnsureKnownBranchAncestryArgs{
			AllBranches:   branches.All,
			BranchTypes:   branches.Types,
			DefaultBranch: mainBranch,
			Lineage:       lineage,
			MainBranch:    mainBranch,
			Runner:        &repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
	}
	if !branches.Types.IsFeatureBranch(targetBranchName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.KillOnlyFeatureBranches)
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	return &killConfig{
		hasOpenChanges: repoStatus.OpenChanges,
		initialBranch:  branches.Initial,
		isOffline:      repo.IsOffline,
		lineage:        lineage,
		mainBranch:     mainBranch,
		noPushHook:     !pushHook,
		previousBranch: previousBranch,
		targetBranch:   *targetBranch,
	}, branchesSnapshot, stashSnapshot, false, nil
}

func (self killConfig) isOnline() bool {
	return !self.isOffline
}

func (self killConfig) targetBranchParent() domain.LocalBranchName {
	return self.lineage.Parent(self.targetBranch.LocalName)
}

func killProgram(config *killConfig) (runProgram, finalUndoProgram program.Program) {
	prog := program.Program{}
	killFeatureBranch(&prog, &finalUndoProgram, *config)
	wrap(&prog, wrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch.LocalName && config.targetBranch.LocalName == config.previousBranch && config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return prog, finalUndoProgram
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog *program.Program, finalUndoProgram *program.Program, config killConfig) {
	if config.targetBranch.HasTrackingBranch() && config.isOnline() {
		prog.Add(&opcode.DeleteTrackingBranch{Branch: config.targetBranch.RemoteName})
	}
	if config.initialBranch == config.targetBranch.LocalName {
		if config.hasOpenChanges {
			prog.Add(&opcode.CommitOpenChanges{})
			// update the registered initial SHA for this branch so that undo restores the just committed changes
			prog.Add(&opcode.UpdateInitialBranchLocalSHA{Branch: config.initialBranch})
			// when undoing, manually undo the just committed changes so that they are uncommitted again
			finalUndoProgram.Add(&opcode.Checkout{Branch: config.targetBranch.LocalName})
			finalUndoProgram.Add(&opcode.UndoLastCommit{})
		}
		prog.Add(&opcode.Checkout{Branch: config.targetBranchParent()})
	}
	prog.Add(&opcode.DeleteLocalBranch{Branch: config.targetBranch.LocalName, Force: false})
	removeBranchFromLineage(removeBranchFromLineageArgs{
		branch:  config.targetBranch.LocalName,
		lineage: config.lineage,
		program: prog,
		parent:  config.targetBranchParent(),
	})
}

func removeBranchFromLineage(args removeBranchFromLineageArgs) {
	childBranches := args.lineage.Children(args.branch)
	for _, child := range childBranches {
		args.program.Add(&opcode.ChangeParent{Branch: child, Parent: args.parent})
	}
	args.program.Add(&opcode.DeleteParentBranch{Branch: args.branch})
}

type removeBranchFromLineageArgs struct {
	branch  domain.LocalBranchName
	lineage config.Lineage
	program *program.Program
	parent  domain.LocalBranchName
}
