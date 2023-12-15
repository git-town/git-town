package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/spf13/cobra"
)

const hackDesc = "Creates a new feature branch off the main development branch"

const hackHelp = `
Syncs the main branch,
forks a new feature branch with the given name off the main branch,
pushes the new feature branch to origin
(if and only if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func hackCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: "basic",
		Args:    cobra.ExactArgs(1),
		Short:   hackDesc,
		Long:    long(hackDesc, hackHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeHack(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeHack(args []string, verbose bool) error {
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
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineHackConfig(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "hack",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          appendProgram(config),
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

func determineHackConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (*appendConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
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
	targetBranch := domain.NewLocalBranchName(args[0])
	mainBranch := repo.Runner.GitTown.MainBranch()
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	shouldNewBranchPush := fc.NewBranchPush(repo.Runner.GitTown.ShouldNewBranchPush())
	isOffline := fc.Offline(repo.Runner.GitTown.IsOffline())
	if branches.All.HasLocalBranch(targetBranch) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.HasMatchingTrackingBranchFor(targetBranch) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	branchNamesToSync := domain.LocalBranchNames{mainBranch}
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	syncUpstream := fc.SyncUpstream(repo.Runner.GitTown.ShouldSyncUpstream())
	syncPerennialStrategy := fc.SyncPerennialStrategy(repo.Runner.GitTown.SyncPerennialStrategy())
	syncFeatureStrategy := fc.SyncFeatureStrategy(repo.Runner.GitTown.SyncFeatureStrategy())
	return &appendConfig{
		branches:                  branches,
		branchesToSync:            branchesToSync,
		targetBranch:              targetBranch,
		parentBranch:              mainBranch,
		hasOpenChanges:            repoStatus.OpenChanges,
		remotes:                   remotes,
		lineage:                   lineage,
		mainBranch:                mainBranch,
		newBranchParentCandidates: domain.LocalBranchNames{mainBranch},
		shouldNewBranchPush:       shouldNewBranchPush,
		previousBranch:            previousBranch,
		syncPerennialStrategy:     syncPerennialStrategy,
		pushHook:                  pushHook,
		isOnline:                  isOffline.ToOnline(),
		syncUpstream:              syncUpstream,
		syncFeatureStrategy:       syncFeatureStrategy,
	}, branchesSnapshot, stashSnapshot, false, fc.Err
}
