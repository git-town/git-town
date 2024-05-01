package cmd

import (
	"os"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const appendDesc = "Create a new feature branch as a child of the current branch"

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
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineAppendData(gitdomain.NewLocalBranchName(arg), repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "append",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            appendProgram(*data),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        &data.dialogTestInputs,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     data.runner,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type appendData struct {
	allBranches               gitdomain.BranchInfos
	branchesToSync            gitdomain.BranchInfos
	config                    configdomain.FullConfig
	dialogTestInputs          components.TestInputs
	dryRun                    bool
	hasOpenChanges            bool
	initialBranch             gitdomain.LocalBranchName
	newBranchParentCandidates gitdomain.LocalBranchNames
	parentBranch              gitdomain.LocalBranchName
	previousBranch            gitdomain.LocalBranchName
	remotes                   gitdomain.Remotes
	runner                    *git.ProdRunner
	targetBranch              gitdomain.LocalBranchName
}

func determineAppendData(targetBranch gitdomain.LocalBranchName, repo *execute.OpenRepoResult, dryRun, verbose bool) (*appendData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	fc := execute.FailureCollector{}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	runner := git.ProdRunner{
		Backend:         repo.Backend,
		CommandsCounter: repo.CommandsCounter,
		Config:          repo.Config,
		FinalMessages:   repo.FinalMessages,
		Frontend:        repo.Frontend,
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 !repoStatus.OpenChanges,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		Runner:                &runner,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	remotes := fc.Remotes(repo.Backend.Remotes())
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	err = execute.EnsureKnownBranchesAncestry(execute.EnsureKnownBranchesAncestryArgs{
		BranchesToVerify: gitdomain.LocalBranchNames{branchesSnapshot.Active},
		Config:           repo.Config,
		DefaultChoice:    repo.Config.Config.MainBranch,
		DialogTestInputs: &dialogTestInputs,
		LocalBranches:    branchesSnapshot.Branches.LocalBranches(),
		MainBranch:       repo.Config.Config.MainBranch,
		Runner:           &runner,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	branchNamesToSync := repo.Config.Config.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	branchesToSync := fc.BranchInfos(branchesSnapshot.Branches.Select(branchNamesToSync...))
	initialAndAncestors := repo.Config.Config.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	slices.Reverse(initialAndAncestors)
	return &appendData{
		allBranches:               branchesSnapshot.Branches,
		branchesToSync:            branchesToSync,
		config:                    repo.Config.Config,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             branchesSnapshot.Active,
		newBranchParentCandidates: initialAndAncestors,
		parentBranch:              branchesSnapshot.Active,
		previousBranch:            previousBranch,
		remotes:                   remotes,
		runner:                    &runner,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSize, false, fc.Err
}

func appendProgram(config appendData) program.Program {
	prog := program.Program{}
	if !config.hasOpenChanges {
		for _, branch := range config.branchesToSync {
			sync.BranchProgram(branch, sync.BranchProgramArgs{
				BranchInfos:   config.allBranches,
				Config:        config.config,
				InitialBranch: config.initialBranch,
				Program:       &prog,
				Remotes:       config.remotes,
				PushBranch:    true,
			})
		}
	}
	prog.Add(&opcodes.CreateAndCheckoutBranchExistingParent{
		Ancestors: config.newBranchParentCandidates,
		Branch:    config.targetBranch,
	})
	if config.remotes.HasOrigin() && config.config.ShouldPushNewBranches() && config.config.IsOnline() {
		prog.Add(&opcodes.CreateTrackingBranch{Branch: config.targetBranch})
	}
	prog.Add(&opcodes.SetExistingParent{
		Branch:    config.targetBranch,
		Ancestors: config.newBranchParentCandidates,
	})
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.initialBranch, config.previousBranch},
	})
	return prog
}
