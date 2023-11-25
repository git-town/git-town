package cmd

import (
	"slices"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/gohacks"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/validate"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/opcode"
	"github.com/git-town/git-town/v10/src/vm/program"
	"github.com/git-town/git-town/v10/src/vm/runstate"
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
		Long:    long(appendDesc, appendHelp),
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
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineAppendConfig(domain.NewLocalBranchName(arg), repo, verbose)
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
		Run:                     &repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		Lineage:                 config.lineage,
		NoPushHook:              !config.pushHook,
	})
}

type appendConfig struct {
	branches                  domain.Branches
	branchesToSync            domain.BranchInfos
	hasOpenChanges            bool
	remotes                   domain.Remotes
	isOffline                 bool
	lineage                   config.Lineage
	mainBranch                domain.LocalBranchName
	newBranchParentCandidates domain.LocalBranchNames
	pushHook                  bool
	parentBranch              domain.LocalBranchName
	previousBranch            domain.LocalBranchName
	pullBranchStrategy        config.PullBranchStrategy
	shouldNewBranchPush       bool
	shouldSyncUpstream        bool
	syncStrategy              config.SyncStrategy
	targetBranch              domain.LocalBranchName
}

func determineAppendConfig(targetBranch domain.LocalBranchName, repo *execute.OpenRepoResult, verbose bool) (*appendConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	fc := gohacks.FailureCollector{}
	pushHook := fc.Bool(repo.Runner.Config.PushHook())
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
	mainBranch := repo.Runner.Config.MainBranch()
	pullBranchStrategy := fc.PullBranchStrategy(repo.Runner.Config.PullBranchStrategy())
	repoStatus := fc.RepoStatus(repo.Runner.Backend.RepoStatus())
	shouldNewBranchPush := fc.Bool(repo.Runner.Config.ShouldNewBranchPush())
	if fc.Err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, fc.Err
	}
	if branches.All.HasLocalBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.HasMatchingTrackingBranchFor(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	updated, err := validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
		AllBranches:   branches.All,
		Backend:       &repo.Runner.Backend,
		BranchTypes:   branches.Types,
		DefaultBranch: mainBranch,
		MainBranch:    mainBranch,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	if updated {
		lineage = repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue) // refresh lineage after ancestry changes
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	syncStrategy := fc.SyncStrategy(repo.Runner.Config.SyncStrategy())
	shouldSyncUpstream := fc.Bool(repo.Runner.Config.ShouldSyncUpstream())
	initialAndAncestors := lineage.BranchAndAncestors(branches.Initial)
	slices.Reverse(initialAndAncestors)
	return &appendConfig{
		branches:                  branches,
		branchesToSync:            branchesToSync,
		hasOpenChanges:            repoStatus.OpenChanges,
		remotes:                   remotes,
		isOffline:                 repo.IsOffline,
		lineage:                   lineage,
		mainBranch:                mainBranch,
		newBranchParentCandidates: initialAndAncestors,
		pushHook:                  pushHook,
		parentBranch:              branches.Initial,
		previousBranch:            previousBranch,
		pullBranchStrategy:        pullBranchStrategy,
		shouldNewBranchPush:       shouldNewBranchPush,
		shouldSyncUpstream:        shouldSyncUpstream,
		syncStrategy:              syncStrategy,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSnapshot, false, fc.Err
}

func appendProgram(config *appendConfig) program.Program {
	prog := program.Program{}
	for _, branch := range config.branchesToSync {
		syncBranchProgram(branch, syncBranchProgramArgs{
			branchTypes:        config.branches.Types,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			program:            &prog,
			remotes:            config.remotes,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
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
	if config.remotes.HasOrigin() && config.shouldNewBranchPush && !config.isOffline {
		prog.Add(&opcode.CreateTrackingBranch{Branch: config.targetBranch, NoPushHook: !config.pushHook})
	}
	wrap(&prog, wrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return prog
}
