package cmd

import (
	"fmt"

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
	allBranches, initialBranch, err := execute.LoadBranches(&repo.Runner, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineHackConfig(args, promptForParent, &repo.Runner, allBranches, initialBranch)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config, &repo.Runner)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "hack",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &repo.Runner, nil, repo.RootDir)
}

func determineHackConfig(args []string, promptForParent bool, run *git.ProdRunner, allBranches git.BranchesSyncStatus, initialBranch string) (*appendConfig, error) {
	fc := failure.Collector{}
	targetBranch := args[0]
	mainBranch := run.Config.MainBranch()
	parentBranch := fc.String(determineParentBranch(targetBranch, promptForParent, run, mainBranch))
	hasOrigin := fc.Bool(run.Backend.HasOrigin())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	isOffline := fc.Bool(run.Config.IsOffline())
	pushHook := fc.Bool(run.Config.PushHook())
	if allBranches.Contains(targetBranch) {
		return nil, fmt.Errorf(messages.BranchAlreadyExists, targetBranch)
	}
	lineage := run.Config.Lineage()
	branchNamesToSync := lineage.BranchesAndAncestors([]string{parentBranch})
	branchesToSync := fc.BranchesSyncStatus(allBranches.Select(branchNamesToSync))
	shouldSyncUpstream := fc.Bool(run.Config.ShouldSyncUpstream())
	hasUpstream := fc.Bool(run.Backend.HasUpstream())
	return &appendConfig{
		branchesToSync:      branchesToSync,
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOrigin:           hasOrigin,
		hasUpstream:         hasUpstream,
		initialBranch:       initialBranch,
		lineage:             lineage,
		mainBranch:          mainBranch,
		shouldNewBranchPush: shouldNewBranchPush,
		pullBranchStrategy:  fc.PullBranchStrategy(run.Config.PullBranchStrategy()),
		pushHook:            pushHook,
		isOffline:           isOffline,
		shouldSyncUpstream:  shouldSyncUpstream,
		syncStrategy:        fc.SyncStrategy(run.Config.SyncStrategy()),
	}, fc.Err
}

func determineParentBranch(targetBranch string, promptForParent bool, run *git.ProdRunner, mainBranch string) (string, error) {
	if promptForParent {
		parentBranch, err := validate.EnterParent(targetBranch, mainBranch, &run.Backend)
		if err != nil {
			return "", err
		}
		err = validate.KnowsBranchAncestors(parentBranch, mainBranch, &run.Backend)
		if err != nil {
			return "", err
		}
		return parentBranch, nil
	}
	return mainBranch, nil
}
