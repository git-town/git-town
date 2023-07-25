package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
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
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: false,
	})
	if err != nil {
		return err
	}
	branchesSyncStatus, initialBranch, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineHackConfig(args, promptForParent, &run, branchesSyncStatus, initialBranch)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "hack",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &run, nil)
}

func determineHackConfig(args []string, promptForParent bool, run *git.ProdRunner, branchesSyncStatus git.BranchesSyncStatus, initialBranch string) (*appendConfig, error) {
	fc := failure.Collector{}
	targetBranch := args[0]
	parentBranch := fc.String(determineParentBranch(targetBranch, promptForParent, run))
	hasOrigin := fc.Bool(run.Backend.HasOrigin())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	isOffline := fc.Bool(run.Config.IsOffline())
	mainBranch := run.Config.MainBranch()
	// TODO: inline this variable?
	hasBranch := branchesSyncStatus.Contains(targetBranch)
	pushHook := fc.Bool(run.Config.PushHook())
	if hasBranch {
		return nil, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	branchesToSync := git.BranchesSyncStatus{*branchesSyncStatus.Lookup(mainBranch)}
	return &appendConfig{
		branchesToSync:      branchesToSync,
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOrigin:           hasOrigin,
		mainBranch:          mainBranch,
		shouldNewBranchPush: shouldNewBranchPush,
		noPushHook:          !pushHook,
		isOffline:           isOffline,
	}, fc.Err
}

func determineParentBranch(targetBranch string, promptForParent bool, run *git.ProdRunner) (string, error) {
	if promptForParent {
		parentBranch, err := validate.EnterParent(targetBranch, run.Config.MainBranch(), &run.Backend)
		if err != nil {
			return "", err
		}
		err = validate.KnowsBranchAncestors(parentBranch, run.Config.MainBranch(), &run.Backend)
		if err != nil {
			return "", err
		}
		return parentBranch, nil
	}
	return run.Config.MainBranch(), nil
}
