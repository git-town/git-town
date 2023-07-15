package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
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
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: true,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
		ValidateIsConfigured:  true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineHackConfig(args, promptForParent, &run)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config, &run)
	if err != nil {
		return err
	}
	runState := runstate.New("hack", stepList)
	return runstate.Execute(runState, &run, nil)
}

func determineHackConfig(args []string, promptForParent bool, run *git.ProdRunner) (*appendConfig, error) {
	fc := failure.Collector{}
	targetBranch := args[0]
	parentBranch := fc.BranchWithParent(determineParentBranch(targetBranch, promptForParent, run))
	hasOrigin := fc.Bool(run.Backend.HasOrigin())
	shouldNewBranchPush := fc.Bool(run.Config.ShouldNewBranchPush())
	isOffline := fc.Bool(run.Config.IsOffline())
	mainBranch := run.Config.MainBranch()
	if fc.Err == nil && hasOrigin && !isOffline {
		fc.Check(run.Frontend.Fetch())
	}
	hasBranch := fc.Bool(run.Backend.HasLocalOrOriginBranch(targetBranch, mainBranch))
	pushHook := fc.Bool(run.Config.PushHook())
	if hasBranch {
		return nil, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	return &appendConfig{
		ancestorBranches:    config.Lineage{},
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOrigin:           hasOrigin,
		mainBranch:          mainBranch,
		shouldNewBranchPush: shouldNewBranchPush,
		noPushHook:          !pushHook,
		isOffline:           isOffline,
	}, fc.Err
}

func determineParentBranch(targetBranch string, promptForParent bool, run *git.ProdRunner) (config.BranchWithParent, error) {
	lineage := run.Config.Lineage()
	if promptForParent {
		parentBranchName, err := validate.EnterParent(targetBranch, run.Config.MainBranch(), &run.Backend)
		if err != nil {
			return config.BranchWithParent{}, err
		}
		err = validate.KnowsBranchAncestors(parentBranchName, run.Config.MainBranch(), &run.Backend)
		if err != nil {
			return config.BranchWithParent{}, err
		}
		parentBranch := lineage.Lookup(parentBranchName)
		return parentBranch, nil
	}
	mainbranch := lineage.Lookup(run.Config.MainBranch())
	return mainbranch, nil
}
