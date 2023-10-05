package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/runvm"
	"github.com/git-town/git-town/v9/src/statistics"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const killDesc = "Removes an obsolete feature branch"

const killHelp = `
Deletes the current or provided branch from the local and origin repositories.
Does not delete perennial branches nor the main branch.`

func killCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "kill [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: killDesc,
		Long:  long(killDesc, killHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeKill(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func executeKill(args []string, debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineKillConfig(args, &repo)
	if err != nil || exit {
		return err
	}
	steps, finalUndoSteps := killSteps(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:             "kill",
		CommandsRun:         statistics.NewCommands(),
		RunSteps:            steps,
		InitialActiveBranch: initialBranchesSnapshot.Active,
		MessagesToUser:      statistics.NewMessages(),
		FinalUndoSteps:      finalUndoSteps,
	}
	return runvm.Execute(runvm.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               nil,
		Lineage:                 config.lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              config.noPushHook,
	})
}

type killConfig struct {
	hasOpenChanges bool
	initialBranch  domain.LocalBranchName
	isOffline      bool
	lineage        config.Lineage
	mainBranch     domain.LocalBranchName
	noPushHook     bool
	previousBranch domain.LocalBranchName
	targetBranch   domain.BranchInfo
}

func determineKillConfig(args []string, repo *execute.OpenRepoResult) (*killConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	targetBranchName := domain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	if !branches.Types.IsFeatureBranch(targetBranchName) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.KillOnlyFeatureBranches)
	}
	targetBranch := branches.All.FindByLocalName(targetBranchName)
	if targetBranch == nil {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	if targetBranch.IsLocal() {
		updated, err := validate.KnowsBranchAncestors(targetBranchName, validate.KnowsBranchAncestorsArgs{
			DefaultBranch: mainBranch,
			Backend:       &repo.Runner.Backend,
			AllBranches:   branches.All,
			BranchTypes:   branches.Types,
			MainBranch:    mainBranch,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
		if updated {
			repo.Runner.Config.Reload()
			lineage = repo.Runner.Config.Lineage()
		}
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	return &killConfig{
		hasOpenChanges: repoStatus.OpenChanges,
		initialBranch:  branches.Initial,
		isOffline:      repo.IsOffline,
		lineage:        lineage,
		mainBranch:     mainBranch,
		noPushHook:     !pushHook,
		previousBranch: previousBranch,
		targetBranch:   *targetBranch,
	}, branchesSnapshot, stashSnapshot, false, nil
}

func (kc killConfig) isOnline() bool {
	return !kc.isOffline
}

func (kc killConfig) targetBranchParent() domain.LocalBranchName {
	return kc.lineage.Parent(kc.targetBranch.LocalName)
}

func killSteps(config *killConfig) (runSteps, finalUndoSteps steps.List) {
	list := steps.List{}
	killFeatureBranch(&list, &finalUndoSteps, *config)
	list.Wrap(steps.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.initialBranch != config.targetBranch.LocalName && config.targetBranch.LocalName == config.previousBranch && config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return list, finalUndoSteps
}

// killFeatureBranch kills the given feature branch everywhere it exists (locally and remotely).
func killFeatureBranch(list *steps.List, finalUndoList *steps.List, config killConfig) {
	if config.targetBranch.HasTrackingBranch() && config.isOnline() {
		list.Add(&step.DeleteTrackingBranch{Branch: config.targetBranch.RemoteName})
	}
	if config.initialBranch == config.targetBranch.LocalName {
		if config.hasOpenChanges {
			list.Add(&step.CommitOpenChanges{})
			// update the registered initial SHA for this branch so that undo restores the just committed changes
			list.Add(&step.UpdateInitialBranchLocalSHA{Branch: config.initialBranch})
			// when undoing, manually undo the just committed changes so that they are uncommitted again
			finalUndoList.Add(&step.Checkout{Branch: config.targetBranch.LocalName})
			finalUndoList.Add(&step.UndoLastCommit{})
		}
		list.Add(&step.Checkout{Branch: config.targetBranchParent()})
	}
	list.Add(&step.DeleteLocalBranch{Branch: config.targetBranch.LocalName, Parent: config.mainBranch.Location(), Force: true})
	removeBranchFromLineage(removeBranchFromLineageArgs{
		branch:  config.targetBranch.LocalName,
		lineage: config.lineage,
		list:    list,
		parent:  config.targetBranchParent(),
	})
}

func removeBranchFromLineage(args removeBranchFromLineageArgs) {
	childBranches := args.lineage.Children(args.branch)
	for _, child := range childBranches {
		args.list.Add(&step.SetParent{Branch: child, ParentBranch: args.parent})
	}
	args.list.Add(&step.DeleteParentBranch{Branch: args.branch})
}

type removeBranchFromLineageArgs struct {
	branch  domain.LocalBranchName
	lineage config.Lineage
	list    *steps.List
	parent  domain.LocalBranchName
}
