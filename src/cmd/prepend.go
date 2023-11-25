package cmd

import (
	"fmt"
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

const prependDesc = "Creates a new feature branch as the parent of the current branch"

const prependHelp = `
Syncs the parent branch,
cuts a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the origin repository
(if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.
`

func prependCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "prepend <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   prependDesc,
		Long:    long(prependDesc, prependHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePrepend(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executePrepend(args []string, verbose bool) error {
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
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determinePrependConfig(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "prepend",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          prependProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		NoPushHook:              config.pushHook,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type prependConfig struct {
	branches                  domain.Branches
	branchesToSync            domain.BranchInfos
	hasOpenChanges            bool
	remotes                   domain.Remotes
	isOffline                 bool
	lineage                   config.Lineage
	mainBranch                domain.LocalBranchName
	newBranchParentCandidates domain.LocalBranchNames
	previousBranch            domain.LocalBranchName
	pullBranchStrategy        config.PullBranchStrategy
	pushHook                  bool
	parentBranch              domain.LocalBranchName
	shouldSyncUpstream        bool
	shouldNewBranchPush       bool
	syncStrategy              config.SyncStrategy
	targetBranch              domain.LocalBranchName
}

func determinePrependConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (*prependConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	fc := gohacks.FailureCollector{}
	pushHook := fc.Bool(repo.Runner.Config.PushHook())
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus := fc.RepoStatus(repo.Runner.Backend.RepoStatus())
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	shouldNewBranchPush := fc.Bool(repo.Runner.Config.ShouldNewBranchPush())
	mainBranch := repo.Runner.Config.MainBranch()
	syncStrategy := fc.SyncStrategy(repo.Runner.Config.SyncStrategy())
	pullBranchStrategy := fc.PullBranchStrategy(repo.Runner.Config.PullBranchStrategy())
	shouldSyncUpstream := fc.Bool(repo.Runner.Config.ShouldSyncUpstream())
	targetBranch := domain.NewLocalBranchName(args[0])
	if branches.All.HasLocalBranch(targetBranch) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.HasMatchingTrackingBranchFor(targetBranch) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	if !branches.Types.IsFeatureBranch(branches.Initial) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.SetParentNoFeatureBranch, branches.Initial)
	}
	updated := fc.Bool(validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
		AllBranches:   branches.All,
		Backend:       &repo.Runner.Backend,
		BranchTypes:   branches.Types,
		DefaultBranch: mainBranch,
		MainBranch:    mainBranch,
	}))
	if updated {
		lineage = repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	parent := lineage.Parent(branches.Initial)
	parentAndAncestors := lineage.BranchAndAncestors(parent)
	slices.Reverse(parentAndAncestors)
	return &prependConfig{
		branches:                  branches,
		branchesToSync:            branchesToSync,
		hasOpenChanges:            repoStatus.OpenChanges,
		remotes:                   remotes,
		isOffline:                 repo.IsOffline,
		lineage:                   lineage,
		mainBranch:                mainBranch,
		newBranchParentCandidates: parentAndAncestors,
		previousBranch:            previousBranch,
		pullBranchStrategy:        pullBranchStrategy,
		pushHook:                  pushHook,
		parentBranch:              parent,
		shouldNewBranchPush:       shouldNewBranchPush,
		shouldSyncUpstream:        shouldSyncUpstream,
		syncStrategy:              syncStrategy,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSnapshot, false, fc.Err
}

func prependProgram(config *prependConfig) program.Program {
	prog := program.Program{}
	for _, branchToSync := range config.branchesToSync {
		syncBranchProgram(branchToSync, syncBranchProgramArgs{
			branchTypes:        config.branches.Types,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			program:            &prog,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			remotes:            config.remotes,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		})
	}
	prog.Add(&opcode.CreateBranchExistingParent{
		Ancestors:  config.newBranchParentCandidates,
		Branch:     config.targetBranch,
		MainBranch: config.mainBranch,
	})
	// set the parent of the newly created branch
	prog.Add(&opcode.SetExistingParent{
		Branch:     config.targetBranch,
		Ancestors:  config.newBranchParentCandidates,
		MainBranch: config.mainBranch,
	})
	// set the parent of the branch prepended to
	prog.Add(&opcode.SetParentIfBranchExists{
		Branch: config.branches.Initial,
		Parent: config.targetBranch,
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
