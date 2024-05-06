package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const prependDesc = "Create a new feature branch as the parent of the current branch"

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
	data, initialBranchesSnapshot, initialStashSize, exit, err := determinePrependData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "prepend",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            prependProgram(data),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        &data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type prependData struct {
	allBranches               gitdomain.BranchInfos
	branchesToSync            gitdomain.BranchInfos
	config                    config.Config
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

func emptyPrependData() prependData {
	return prependData{} //exhaustruct:ignore
}

func determinePrependData(args []string, repo execute.OpenRepoResult, dryRun, verbose bool) (prependData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return emptyPrependData(), gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	fc := execute.FailureCollector{}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 !repoStatus.OpenChanges,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return emptyPrependData(), branchesSnapshot, stashSize, exit, err
	}
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	remotes := fc.Remotes(repo.Backend.Remotes())
	targetBranch := gitdomain.NewLocalBranchName(args[0])
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return emptyPrependData(), branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return emptyPrependData(), branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	repo.Config, exit, err = validate.Config(validate.ConfigArgs{
		Backend:            &repo.Backend,
		BranchesToValidate: gitdomain.LocalBranchNames{branchesSnapshot.Active},
		FinalMessages:      repo.FinalMessages,
		LocalBranches:      branchesSnapshot.Branches.LocalBranches().Names(),
		TestInputs:         &dialogTestInputs,
		Unvalidated:        repo.Config,
	})
	if err != nil || exit {
		return emptyPrependData(), branchesSnapshot, stashSize, exit, err
	}
	branchNamesToSync := repo.Config.Config.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	branchesToSync := fc.BranchInfos(branchesSnapshot.Branches.Select(branchNamesToSync...))
	parent, hasParent := repo.Config.Config.Lineage.Parent(branchesSnapshot.Active).Get()
	if !hasParent {
		return emptyPrependData(), branchesSnapshot, stashSize, false, fmt.Errorf(messages.SetParentNoFeatureBranch, branchesSnapshot.Active)
	}
	parentAndAncestors := repo.Config.Config.Lineage.BranchAndAncestors(parent)
	slices.Reverse(parentAndAncestors)
	return prependData{
		allBranches:               branchesSnapshot.Branches,
		branchesToSync:            branchesToSync,
		config:                    repo.Config,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             branchesSnapshot.Active,
		newBranchParentCandidates: parentAndAncestors,
		parentBranch:              parent,
		previousBranch:            previousBranch,
		remotes:                   remotes,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSize, false, fc.Err
}

func prependProgram(data prependData) program.Program {
	prog := program.Program{}
	for _, branchToSync := range data.branchesToSync {
		sync.BranchProgram(branchToSync, sync.BranchProgramArgs{
			BranchInfos:   data.allBranches,
			Config:        data.config.Config,
			InitialBranch: data.initialBranch,
			Program:       &prog,
			PushBranch:    true,
			Remotes:       data.remotes,
		})
	}
	prog.Add(&opcodes.CreateAndCheckoutBranchExistingParent{
		Ancestors: data.newBranchParentCandidates,
		Branch:    data.targetBranch,
	})
	// set the parent of the newly created branch
	prog.Add(&opcodes.SetExistingParent{
		Branch:    data.targetBranch,
		Ancestors: data.newBranchParentCandidates,
	})
	// set the parent of the branch prepended to
	prog.Add(&opcodes.SetParentIfBranchExists{
		Branch: data.initialBranch,
		Parent: data.targetBranch,
	})
	if data.remotes.HasOrigin() && data.config.Config.ShouldPushNewBranches() && data.config.Config.IsOnline() {
		prog.Add(&opcodes.CreateTrackingBranch{Branch: data.targetBranch})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{data.previousBranch},
	})
	return prog
}
