package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/spf13/cobra"
)

const pruneBranchesDesc = "Deletes local branches whose tracking branch no longer exists"

const pruneBranchesHelp = `
Deletes branches whose tracking branch no longer exists from the local repository.
This usually means the branch was shipped or killed on another machine.`

func pruneBranchesCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "prune-branches",
		Args:  cobra.NoArgs,
		Short: pruneBranchesDesc,
		Long:  long(pruneBranchesDesc, pruneBranchesHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePruneBranches(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func executePruneBranches(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: true,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determinePruneBranchesConfig(repo, debug)
	if err != nil || exit {
		return err
	}
	runState := runstate.RunState{
		Command:             "prune-branches",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunSteps:            pruneBranchesSteps(config),
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Debug:                   debug,
		Lineage:                 config.lineage,
		NoPushHook:              !config.pushHook,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type pruneBranchesConfig struct {
	branches                  domain.Branches
	branchesWithDeletedRemote domain.LocalBranchNames
	hasOpenChanges            bool
	lineage                   config.Lineage
	mainBranch                domain.LocalBranchName
	previousBranch            domain.LocalBranchName
	pullBranchStrategy        config.PullBranchStrategy
	pushHook                  bool
	remotes                   domain.Remotes
	shouldSyncUpstream        bool
	syncStrategy              config.SyncStrategy
}

func determinePruneBranchesConfig(repo *execute.OpenRepoResult, debug bool) (*pruneBranchesConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Debug:                 debug,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	syncStrategy, err := repo.Runner.Config.SyncStrategy()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	shouldSyncUpstream, err := repo.Runner.Config.ShouldSyncUpstream()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	remotes, err := repo.Runner.Backend.Remotes()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	pullBranchStrategy, err := repo.Runner.Config.PullBranchStrategy()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	return &pruneBranchesConfig{
		branches:                  branches,
		branchesWithDeletedRemote: branches.All.LocalBranchesWithDeletedTrackingBranches().Names(),
		hasOpenChanges:            repoStatus.OpenChanges,
		lineage:                   lineage,
		mainBranch:                repo.Runner.Config.MainBranch(),
		previousBranch:            repo.Runner.Backend.PreviouslyCheckedOutBranch(),
		pullBranchStrategy:        pullBranchStrategy,
		pushHook:                  pushHook,
		remotes:                   remotes,
		shouldSyncUpstream:        shouldSyncUpstream,
		syncStrategy:              syncStrategy,
	}, branchesSnapshot, stashSnapshot, exit, err
}

func pruneBranchesSteps(config *pruneBranchesConfig) steps.List {
	list := steps.List{}
	for _, branchWithDeletedRemote := range config.branchesWithDeletedRemote {
		syncBranchSteps(&list, syncBranchStepsArgs{
			branch:             *config.branches.All.FindByLocalName(branchWithDeletedRemote),
			branchTypes:        config.branches.Types,
			remotes:            config.remotes,
			isOffline:          true,
			lineage:            config.lineage,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		})
		parent := config.lineage.Parent(branchWithDeletedRemote)
		if parent.IsEmpty() {
			parent = config.mainBranch
		}
		list.Add(&step.IfBranchHasChanges{
			Branch: branchWithDeletedRemote,
			Parent: parent,
			IsEmptySteps: []step.Step{
				&step.Checkout{Branch: parent},
				&step.DeleteLocalBranch{
					Branch: branchWithDeletedRemote,
					Force:  false,
					Parent: parent.Location(),
				},
				&step.RemoveBranchFromLineage{
					Branch: branchWithDeletedRemote,
				},
				&step.RemoveFromPerennialBranches{Branch: branchWithDeletedRemote},
			},
			HasChangesSteps: []step.Step{
				&step.QueueMessage{
					Message: fmt.Sprintf(messages.BranchHasUnshippedChanges, branchWithDeletedRemote),
				},
			},
		})
	}
	list.Wrap(steps.WrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return list
}
