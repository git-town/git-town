package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/validate"
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
	addDebugFlag, readDebugFlag := flags.Debug()
	addPromptFlag, readPromptFlag := flags.Bool("prompt", "p", "Prompt for the parent branch")
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: "basic",
		Args:    cobra.ExactArgs(1),
		Short:   hackDesc,
		Long:    long(hackDesc, hackHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return hack(args, readPromptFlag(cmd), readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	addPromptFlag(&cmd)
	return &cmd
}

func hack(args []string, promptForParent, debug bool) error {
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 true,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineHackConfig(args, promptForParent, &repo.Runner)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "hack",
		RunStepList: stepList,
	}
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: nil,
		RootDir:   repo.RootDir,
	})
}

func determineHackConfig(args []string, promptForParent bool, run *git.ProdRunner) (*appendConfig, error) {
	fc := failure.Collector{}
	branches := fc.Branches(execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	}))
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges := fc.Bool(run.Backend.HasOpenChanges())
	targetBranch := domain.NewLocalBranchName(args[0])
	mainBranch := run.Config.MainBranch()
	lineage := run.Config.Lineage()
	parentBranch, updated, err := determineParentBranch(determineParentBranchArgs{
		backend:         &run.Backend,
		branches:        branches,
		lineage:         lineage,
		mainBranch:      mainBranch,
		promptForParent: promptForParent,
		targetBranch:    targetBranch,
	})
	if err != nil {
		return nil, err
	}
	if updated {
		lineage = run.Config.Lineage()
	}
	remotes := fc.Strings(run.Backend.Remotes())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	isOffline := fc.Bool(run.Config.IsOffline())
	pushHook := fc.Bool(run.Config.PushHook())
	if branches.All.HasLocalBranch(targetBranch) {
		return nil, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branches.All.HasRemoteBranchFor(targetBranch) {
		return nil, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	branchNamesToSync := lineage.BranchesAndAncestors(domain.LocalBranchNames{parentBranch})
	branchesToSync := fc.BranchesSyncStatus(branches.All.Select(branchNamesToSync))
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	pullBranchStrategy := fc.PullBranchStrategy(run.Config.PullBranchStrategy())
	syncStrategy := fc.SyncStrategy(run.Config.SyncStrategy())
	return &appendConfig{
		durations:           branches.Durations,
		branchesToSync:      branchesToSync,
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOpenChanges:      hasOpenChanges,
		remotes:             remotes,
		initialBranch:       branches.Initial,
		lineage:             lineage,
		mainBranch:          mainBranch,
		shouldNewBranchPush: shouldNewBranchPush,
		previousBranch:      previousBranch,
		pullBranchStrategy:  pullBranchStrategy,
		pushHook:            pushHook,
		isOffline:           isOffline,
		shouldSyncUpstream:  shouldSyncUpstream,
		syncStrategy:        syncStrategy,
	}, fc.Err
}

func determineParentBranch(args determineParentBranchArgs) (parentBranch domain.LocalBranchName, updated bool, err error) {
	if !args.promptForParent {
		return args.mainBranch, false, nil
	}
	parentBranch, err = dialog.EnterParent(args.targetBranch, args.mainBranch, args.lineage, args.branches.All)
	if err != nil {
		return domain.LocalBranchName{}, true, err
	}
	_, err = validate.KnowsBranchAncestors(parentBranch, validate.KnowsBranchAncestorsArgs{
		AllBranches:     args.branches.All,
		Backend:         args.backend,
		BranchDurations: args.branches.Durations,
		DefaultBranch:   args.mainBranch,
		Lineage:         args.lineage,
		MainBranch:      args.mainBranch,
	})
	if err != nil {
		return domain.LocalBranchName{}, true, err
	}
	return parentBranch, true, nil
}

type determineParentBranchArgs struct {
	backend         *git.BackendCommands
	branches        git.Branches
	lineage         config.Lineage
	mainBranch      domain.LocalBranchName
	promptForParent bool
	targetBranch    domain.LocalBranchName
}
