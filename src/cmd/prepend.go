package cmd

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
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
		Run:                     repo.Runner,
		Connector:               nil,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		NoPushHook:              config.pushHook.Negate(),
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
	isOnline                  configdomain.Online
	lineage                   configdomain.Lineage
	mainBranch                domain.LocalBranchName
	newBranchParentCandidates domain.LocalBranchNames
	previousBranch            domain.LocalBranchName
	syncPerennialStrategy     configdomain.SyncPerennialStrategy
	pushHook                  configdomain.PushHook
	parentBranch              domain.LocalBranchName
	syncUpstream              configdomain.SyncUpstream
	shouldNewBranchPush       configdomain.NewBranchPush
	syncFeatureStrategy       configdomain.SyncFeatureStrategy
	targetBranch              domain.LocalBranchName
}

func determinePrependConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (*prependConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.GitTown.Lineage(repo.Runner.Backend.GitTown.RemoveLocalConfigValue)
	fc := configdomain.FailureCollector{}
	pushHook := fc.PushHook(repo.Runner.GitTown.PushHook())
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
	shouldNewBranchPush := fc.NewBranchPush(repo.Runner.GitTown.ShouldNewBranchPush())
	mainBranch := repo.Runner.GitTown.MainBranch()
	syncFeatureStrategy := fc.SyncFeatureStrategy(repo.Runner.GitTown.SyncFeatureStrategy())
	syncPerennialStrategy := fc.SyncPerennialStrategy(repo.Runner.GitTown.SyncPerennialStrategy())
	syncUpstream := fc.SyncUpstream(repo.Runner.GitTown.ShouldSyncUpstream())
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
	parent := lineage.Parent(branches.Initial)
	parentAndAncestors := lineage.BranchAndAncestors(parent)
	slices.Reverse(parentAndAncestors)
	return &prependConfig{
		branches:                  branches,
		branchesToSync:            branchesToSync,
		hasOpenChanges:            repoStatus.OpenChanges,
		remotes:                   remotes,
		isOnline:                  repo.IsOffline.ToOnline(),
		lineage:                   lineage,
		mainBranch:                mainBranch,
		newBranchParentCandidates: parentAndAncestors,
		previousBranch:            previousBranch,
		syncPerennialStrategy:     syncPerennialStrategy,
		pushHook:                  pushHook,
		parentBranch:              parent,
		shouldNewBranchPush:       shouldNewBranchPush,
		syncUpstream:              syncUpstream,
		syncFeatureStrategy:       syncFeatureStrategy,
		targetBranch:              targetBranch,
	}, branchesSnapshot, stashSnapshot, false, fc.Err
}

func prependProgram(config *prependConfig) program.Program {
	prog := program.Program{}
	for _, branchToSync := range config.branchesToSync {
		syncBranchProgram(branchToSync, syncBranchProgramArgs{
			branchInfos:           config.branches.All,
			branchTypes:           config.branches.Types,
			isOnline:              config.isOnline,
			lineage:               config.lineage,
			program:               &prog,
			mainBranch:            config.mainBranch,
			syncPerennialStrategy: config.syncPerennialStrategy,
			pushBranch:            true,
			pushHook:              config.pushHook,
			remotes:               config.remotes,
			syncUpstream:          config.syncUpstream,
			syncFeatureStrategy:   config.syncFeatureStrategy,
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
	if config.remotes.HasOrigin() && config.shouldNewBranchPush.Bool() && config.isOnline.Bool() {
		prog.Add(&opcode.CreateTrackingBranch{Branch: config.targetBranch, NoPushHook: config.pushHook.Negate()})
	}
	wrap(&prog, wrapOptions{
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: domain.LocalBranchNames{config.previousBranch},
	})
	return prog
}
