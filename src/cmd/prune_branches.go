package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/gohacks"
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
	isOffline                 bool
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
	fc := gohacks.FailureCollector{}
	repoStatus := fc.RepoStatus(repo.Runner.Backend.RepoStatus())
	syncStrategy := fc.SyncStrategy(repo.Runner.Config.SyncStrategy())
	shouldSyncUpstream := fc.Bool(repo.Runner.Config.ShouldSyncUpstream())
	remotes := fc.Remotes(repo.Runner.Backend.Remotes())
	pullBranchStrategy := fc.PullBranchStrategy(repo.Runner.Config.PullBranchStrategy())
	isOffline := fc.Bool(repo.Runner.Config.IsOffline())
	return &pruneBranchesConfig{
		branches:                  branches,
		branchesWithDeletedRemote: branches.All.LocalBranchesWithDeletedTrackingBranches().Names(),
		hasOpenChanges:            repoStatus.OpenChanges,
		isOffline:                 isOffline,
		lineage:                   lineage,
		mainBranch:                repo.Runner.Config.MainBranch(),
		previousBranch:            repo.Runner.Backend.PreviouslyCheckedOutBranch(),
		pullBranchStrategy:        pullBranchStrategy,
		pushHook:                  pushHook,
		remotes:                   remotes,
		shouldSyncUpstream:        shouldSyncUpstream,
		syncStrategy:              syncStrategy,
	}, branchesSnapshot, stashSnapshot, exit, fc.Err
}

func pruneBranchesSteps(config *pruneBranchesConfig) steps.List {
	list := steps.List{}
	for _, branchWithDeletedRemote := range config.branchesWithDeletedRemote {
		parent := config.lineage.Parent(branchWithDeletedRemote)
		if !parent.IsEmpty() {
			parentInfo := config.branches.All.FindByLocalName(parent)
			syncBranchSteps(*parentInfo, syncBranchStepsArgs{
				branchTypes:        config.branches.Types,
				remotes:            config.remotes,
				isOffline:          config.isOffline,
				lineage:            config.lineage,
				list:               &list,
				mainBranch:         config.mainBranch,
				pullBranchStrategy: config.pullBranchStrategy,
				pushBranch:         true,
				pushHook:           config.pushHook,
				shouldSyncUpstream: config.shouldSyncUpstream,
				syncStrategy:       config.syncStrategy,
			})
		}
		if parent.IsEmpty() {
			parent = config.mainBranch
		}
		pullParentBranchOfCurrentFeatureBranchStep(&list, branchWithDeletedRemote, config.syncStrategy)
		list.Add(&step.IfElse{
			Condition: func(backend *git.BackendCommands) (bool, error) {
				return backend.BranchHasUnmergedChanges(branchWithDeletedRemote)
			},
			TrueSteps: []step.Step{
				&step.QueueMessage{
					Message: fmt.Sprintf(messages.BranchDeletedHasUnmergedChanges, branchWithDeletedRemote),
				},
			},
			FalseSteps: []step.Step{
				&step.Checkout{Branch: parent},
				&step.DeleteLocalBranch{
					Branch: branchWithDeletedRemote,
					Force:  false,
				},
				&step.RemoveBranchFromLineage{
					Branch: branchWithDeletedRemote,
				},
				&step.RemoveFromPerennialBranches{Branch: branchWithDeletedRemote},
			},
		})
	}
	list.Add(&step.CheckoutIfExists{Branch: config.branches.Initial})
	list.Wrap(steps.WrapOptions{
		RunInGitRoot:     false,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return list
}
