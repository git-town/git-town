package cmd

import (
	"os"
	"slices"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/sync"
	"github.com/git-town/git-town/v12/src/vm/interpreter"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/spf13/cobra"
)

const appendDesc = "Creates a new feature branch as a child of the current branch"

const appendHelp = `
Syncs the current branch, forks a new feature branch with the given name off the current branch, makes the new branch a child of the current branch, pushes the new feature branch to the origin repository (if and only if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func appendCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "append <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    cmdhelpers.Long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAppend(args[0], readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeAppend(arg string, dryRun, verbose bool) error {
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineAppendConfig(gitdomain.NewLocalBranchName(arg), repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "append",
		DryRun:              dryRun,
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          appendProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                &runState,
		Verbose:                 verbose,
	})
}

type appendConfig struct {
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

func determineAppendConfig(targetBranch gitdomain.LocalBranchName, repo *execute.OpenRepoResult, dryRun, verbose bool) (*appendConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	fc := execute.FailureCollector{}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	repoStatus := fc.RepoStatus(repo.Runner.Backend.RepoStatus())
	if fc.Err != nil {
		return nil, branchesSnapshot, stashSize, false, fc.Err
	}
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	err = execute.EnsureKnownBranchAncestry(branchesSnapshot.Active, execute.EnsureKnownBranchAncestryArgs{
		Config:           &repo.Runner.FullConfig,
		AllBranches:      branchesSnapshot.Branches,
		DefaultBranch:    repo.Runner.MainBranch,
		DialogTestInputs: &dialogTestInputs,
		Runner:           repo.Runner,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	branchNamesToSync := repo.Runner.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	branchesToSync := fc.BranchInfos(branchesSnapshot.Branches.Select(branchNamesToSync))
	initialAndAncestors := repo.Runner.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	slices.Reverse(initialAndAncestors)
	return &appendConfig{
		FullConfig:                &repo.Runner.FullConfig,
		allBranches:               branchesSnapshot.Branches,
		branchesToSync:            branchesToSync,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             branchesSnapshot.Active,
		newBranchParentCandidates: initialAndAncestors,
		parentBranch:              branchesSnapshot.Active,
		previousBranch:            previousBranch,
		remotes:                   remotes,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSize, false, fc.Err
}

func appendProgram(config *appendConfig) program.Program {
	prog := program.Program{}
	for _, branch := range config.branchesToSync {
		sync.BranchProgram(branch, sync.BranchProgramArgs{
			Config:      config.FullConfig,
			BranchInfos: config.allBranches,
			Program:     &prog,
			Remotes:     config.remotes,
			PushBranch:  true,
		})
	}
	prog.Add(&opcodes.CreateBranchExistingParent{
		Ancestors: config.newBranchParentCandidates,
		Branch:    config.targetBranch,
	})
	prog.Add(&opcodes.SetExistingParent{
		Branch:    config.targetBranch,
		Ancestors: config.newBranchParentCandidates,
	})
	prog.Add(&opcodes.Checkout{Branch: config.targetBranch})
	if config.remotes.HasOrigin() && config.ShouldPushNewBranches() && config.IsOnline() {
		prog.Add(&opcodes.CreateTrackingBranch{Branch: config.targetBranch})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.initialBranch, config.previousBranch},
	})
	return prog
}
