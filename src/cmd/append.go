package cmd

import (
	"slices"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/sync"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/spf13/cobra"
)

const appendDesc = "Creates a new feature branch as a child of the current branch"

const appendHelp = `
Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the origin repository
(if and only if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func appendCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "append <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    cmdhelpers.Long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAppend(args[0], readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeAppend(arg string, verbose bool) error {
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
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineAppendConfig(gitdomain.NewLocalBranchName(arg), repo, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "append",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          appendProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		Lineage:                 config.lineage,
		NoPushHook:              config.pushHook.Negate(),
	})
}

type appendConfig struct {
	branches                  configdomain.Branches
	branchesToSync            gitdomain.BranchInfos
	hasOpenChanges            bool
	remotes                   gitdomain.Remotes
	isOnline                  configdomain.Online
	lineage                   configdomain.Lineage
	mainBranch                gitdomain.LocalBranchName
	newBranchParentCandidates gitdomain.LocalBranchNames
	pushHook                  configdomain.PushHook
	parentBranch              gitdomain.LocalBranchName
	previousBranch            gitdomain.LocalBranchName
	syncPerennialStrategy     configdomain.SyncPerennialStrategy
	shouldNewBranchPush       configdomain.NewBranchPush
	syncUpstream              configdomain.SyncUpstream
	syncFeatureStrategy       configdomain.SyncFeatureStrategy
	targetBranch              gitdomain.LocalBranchName
}

func determineAppendConfig(targetBranch gitdomain.LocalBranchName, repo *execute.OpenRepoResult, verbose bool) (*appendConfig, undodomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	lineage := repo.Runner.GitTown.Lineage(repo.Runner.Backend.GitTown.RemoveLocalConfigValue)
	fc := execute.FailureCollector{}
	pushHook := repo.Runner.GitTown.PushHook
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		Lineage:               lineage,
		HandleUnfinishedState: true,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	mainBranch := repo.Runner.GitTown.MainBranch
	syncPerennialStrategy := repo.Runner.GitTown.SyncPerennialStrategy
	repoStatus := fc.RepoStatus(repo.Runner.Backend.RepoStatus())
	shouldNewBranchPush := repo.Runner.GitTown.NewBranchPush
	if fc.Err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, fc.Err
	}
	if branches.All.HasLocalBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.HasMatchingTrackingBranchFor(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	branches.Types, lineage, err = execute.EnsureKnownBranchAncestry(branches.Initial, execute.EnsureKnownBranchAncestryArgs{
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
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	syncFeatureStrategy := repo.Runner.GitTown.SyncFeatureStrategy
	syncUpstream := repo.Runner.GitTown.SyncUpstream
	initialAndAncestors := lineage.BranchAndAncestors(branches.Initial)
	slices.Reverse(initialAndAncestors)
	return &appendConfig{
		branches:                  branches,
		branchesToSync:            branchesToSync,
		hasOpenChanges:            repoStatus.OpenChanges,
		remotes:                   remotes,
		isOnline:                  repo.IsOffline.ToOnline(),
		lineage:                   lineage,
		mainBranch:                mainBranch,
		newBranchParentCandidates: initialAndAncestors,
		pushHook:                  pushHook,
		parentBranch:              branches.Initial,
		previousBranch:            previousBranch,
		syncPerennialStrategy:     syncPerennialStrategy,
		shouldNewBranchPush:       shouldNewBranchPush,
		syncUpstream:              syncUpstream,
		syncFeatureStrategy:       syncFeatureStrategy,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSnapshot, false, fc.Err
}

func appendProgram(config *appendConfig) program.Program {
	prog := program.Program{}
	for _, branch := range config.branchesToSync {
		sync.BranchProgram(branch, sync.BranchProgramArgs{
			BranchInfos:           config.branches.All,
			BranchTypes:           config.branches.Types,
			IsOnline:              config.isOnline,
			Lineage:               config.lineage,
			Program:               &prog,
			Remotes:               config.remotes,
			MainBranch:            config.mainBranch,
			SyncPerennialStrategy: config.syncPerennialStrategy,
			PushBranch:            true,
			PushHook:              config.pushHook,
			SyncUpstream:          config.syncUpstream,
			SyncFeatureStrategy:   config.syncFeatureStrategy,
		})
	}
	prog.Add(&opcode.CreateBranchExistingParent{
		Ancestors:  config.newBranchParentCandidates,
		Branch:     config.targetBranch,
		MainBranch: config.mainBranch,
	})
	prog.Add(&opcode.SetExistingParent{
		Branch:     config.targetBranch,
		Ancestors:  config.newBranchParentCandidates,
		MainBranch: config.mainBranch,
	})
	prog.Add(&opcode.Checkout{Branch: config.targetBranch})
	if config.remotes.HasOrigin() && config.shouldNewBranchPush.Bool() && config.isOnline.Bool() {
		prog.Add(&opcode.CreateTrackingBranch{Branch: config.targetBranch, NoPushHook: config.pushHook.Negate()})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.branches.Initial, config.previousBranch},
	})
	return prog
}
