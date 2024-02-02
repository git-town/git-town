package cmd

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/sync"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/spf13/cobra"
)

const prependDesc = "Creates a new feature branch as the parent of the current branch"

const prependHelp = `
Syncs the parent branch, cuts a new feature branch with the given name off the parent branch, makes the new branch the parent of the current branch, pushes the new feature branch to the origin repository (if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.`

func prependCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "prepend <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   prependDesc,
		Long:    cmdhelpers.Long(prependDesc, prependHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePrepend(args, readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executePrepend(args []string, dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determinePrependConfig(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "prepend",
		DryRun:              dryRun,
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          prependProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		FullConfig:              config.FullConfig,
		RunState:                &runState,
		Run:                     repo.Runner,
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type prependConfig struct {
	*configdomain.FullConfig
	allBranches               gitdomain.BranchInfos
	branchesToSync            gitdomain.BranchInfos
	dialogTestInputs          components.TestInputs
	dryRun                    bool
	hasOpenChanges            bool
	initialBranch             gitdomain.LocalBranchName
	newBranchParentCandidates gitdomain.LocalBranchNames
	parentBranch              gitdomain.LocalBranchName
	previousBranch            gitdomain.LocalBranchName
	remotes                   gitdomain.Remotes
	targetBranch              gitdomain.LocalBranchName
}

func determinePrependConfig(args []string, repo *execute.OpenRepoResult, dryRun, verbose bool) (*prependConfig, gitdomain.BranchesStatus, gitdomain.StashSize, bool, error) {
	fc := execute.FailureCollector{}
	branchesSnapshot, stashSnapshot, dialogTestInputs, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus := fc.RepoStatus(repo.Runner.Backend.RepoStatus())
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	targetBranch := gitdomain.NewLocalBranchName(args[0])
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	if !repo.Runner.IsFeatureBranch(branchesSnapshot.Active) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.SetParentNoFeatureBranch, branchesSnapshot.Active)
	}
	err = execute.EnsureKnownBranchAncestry(branchesSnapshot.Active, execute.EnsureKnownBranchAncestryArgs{
		Config:           &repo.Runner.FullConfig,
		AllBranches:      branchesSnapshot.Branches,
		DefaultBranch:    repo.Runner.MainBranch,
		DialogTestInputs: &dialogTestInputs,
		Runner:           repo.Runner,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	branchNamesToSync := repo.Runner.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	branchesToSync := fc.BranchesSyncStatus(branchesSnapshot.Branches.Select(branchNamesToSync))
	parent := repo.Runner.Lineage.Parent(branchesSnapshot.Active)
	parentAndAncestors := repo.Runner.Lineage.BranchAndAncestors(parent)
	slices.Reverse(parentAndAncestors)
	return &prependConfig{
		FullConfig:                &repo.Runner.FullConfig,
		allBranches:               branchesSnapshot.Branches,
		branchesToSync:            branchesToSync,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             branchesSnapshot.Active,
		newBranchParentCandidates: parentAndAncestors,
		parentBranch:              parent,
		previousBranch:            previousBranch,
		remotes:                   remotes,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSnapshot, false, fc.Err
}

func prependProgram(config *prependConfig) program.Program {
	prog := program.Program{}
	for _, branchToSync := range config.branchesToSync {
		sync.BranchProgram(branchToSync, sync.BranchProgramArgs{
			Config:      config.FullConfig,
			BranchInfos: config.allBranches,
			Program:     &prog,
			PushBranch:  true,
			Remotes:     config.remotes,
		})
	}
	prog.Add(&opcode.CreateBranchExistingParent{
		Ancestors: config.newBranchParentCandidates,
		Branch:    config.targetBranch,
	})
	// set the parent of the newly created branch
	prog.Add(&opcode.SetExistingParent{
		Branch:    config.targetBranch,
		Ancestors: config.newBranchParentCandidates,
	})
	// set the parent of the branch prepended to
	prog.Add(&opcode.SetParentIfBranchExists{
		Branch: config.initialBranch,
		Parent: config.targetBranch,
	})
	prog.Add(&opcode.Checkout{Branch: config.targetBranch})
	if config.remotes.HasOrigin() && config.ShouldPushNewBranches() && config.IsOnline() {
		prog.Add(&opcode.CreateTrackingBranch{Branch: config.targetBranch})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})
	return prog
}
