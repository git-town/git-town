package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/sync/syncprograms"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
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
		Long:  cmdhelpers.Long(killDesc, killHelp),
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
		Run:                     repo.Runner,
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
	branchToKill   gitdomain.BranchInfo
	branchWhenDone gitdomain.LocalBranchName
	hasOpenChanges bool
	initialBranch  gitdomain.LocalBranchName
	isOnline       configdomain.Online
	lineage        configdomain.Lineage
	mainBranch     gitdomain.LocalBranchName
	noPushHook     configdomain.NoPushHook
	previousBranch gitdomain.LocalBranchName
}

func determineKillConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (*killConfig, undodomain.BranchesSnapshot, undodomain.StashSnapshot, bool, error) {
	lineage := repo.Runner.GitTown.Lineage(repo.Runner.Backend.GitTown.RemoveLocalConfigValue)
	pushHook := repo.Runner.GitTown.PushHook
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
	mainBranch := repo.Runner.GitTown.MainBranch
	branchNameToKill := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	branchToKill := branches.All.FindByLocalName(branchNameToKill)
	if branchToKill == nil {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToKill)
	}
	if branchToKill.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return nil, branchesSnapshot, stashSnapshot, exit, fmt.Errorf(messages.KillBranchOtherWorktree, branchNameToKill)
	}
	if branchToKill.IsLocal() {
		branches.Types, lineage, err = execute.EnsureKnownBranchAncestry(branchToKill.LocalName, execute.EnsureKnownBranchAncestryArgs{
			AllBranches:   branches.All,
			BranchTypes:   branches.Types,
			DefaultBranch: mainBranch,
			Lineage:       lineage,
			MainBranch:    mainBranch,
			Runner:        repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
	}
	if !branches.Types.IsFeatureBranch(branchToKill.LocalName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.KillOnlyFeatureBranches)
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	var branchWhenDone gitdomain.LocalBranchName
	if branchNameToKill == branches.Initial {
		branchWhenDone = previousBranch
	} else {
		branchWhenDone = branches.Initial
	}
	return &killConfig{
		branchToKill:   *branchToKill,
		branchWhenDone: branchWhenDone,
		hasOpenChanges: repoStatus.OpenChanges,
		initialBranch:  branches.Initial,
		isOnline:       repo.IsOffline.ToOnline(),
		lineage:        lineage,
		mainBranch:     mainBranch,
		noPushHook:     pushHook.Negate(),
		previousBranch: previousBranch,
	}, branchesSnapshot, stashSnapshot, false, nil
}

func (self killConfig) branchToKillParent() gitdomain.LocalBranchName {
	return self.lineage.Parent(self.branchToKill.LocalName)
}

func killProgram(config *killConfig) (runProgram, finalUndoProgram program.Program) {
	prog := program.Program{}
	killFeatureBranch(&prog, &finalUndoProgram, *config)
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		RunInGitRoot:             true,
		StashOpenChanges:         config.initialBranch != config.branchToKill.LocalName && config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch, config.initialBranch},
	})
	return prog, finalUndoProgram
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(prog *program.Program, finalUndoProgram *program.Program, config killConfig) {
	if config.branchToKill.HasTrackingBranch() && config.isOnline.Bool() {
		prog.Add(&opcode.DeleteTrackingBranch{Branch: config.branchToKill.RemoteName})
	}
	if config.initialBranch == config.branchToKill.LocalName {
		if config.hasOpenChanges {
			prog.Add(&opcode.CommitOpenChanges{})
			// update the registered initial SHA for this branch so that undo restores the just committed changes
			prog.Add(&opcode.UpdateInitialBranchLocalSHA{Branch: config.initialBranch})
			// when undoing, manually undo the just committed changes so that they are uncommitted again
			finalUndoProgram.Add(&opcode.Checkout{Branch: config.branchToKill.LocalName})
			finalUndoProgram.Add(&opcode.UndoLastCommit{})
		}
		prog.Add(&opcode.Checkout{Branch: config.branchWhenDone})
	}
	prog.Add(&opcode.DeleteLocalBranch{Branch: config.branchToKill.LocalName, Force: false})
	syncprograms.RemoveBranchFromLineage(syncprograms.RemoveBranchFromLineageArgs{
		Branch:  config.branchToKill.LocalName,
		Lineage: config.lineage,
		Program: prog,
		Parent:  config.branchToKillParent(),
	})
}
