package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/flags"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/validate"
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
	repo, exit, err := LoadProdRepo(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: true,
		validateGitversion:    true,
		validateIsRepository:  true,
		validateIsConfigured:  true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineHackConfig(args, promptForParent, &repo)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config, &repo)
	if err != nil {
		return err
	}
	runState := runstate.New("hack", stepList)
	return runstate.Execute(runState, &repo, nil)
}

func determineHackConfig(args []string, promptForParent bool, repo *git.ProdRepo) (*appendConfig, error) {
	ec := runstate.ErrorChecker{}
	targetBranch := args[0]
	parentBranch := ec.String(determineParentBranch(targetBranch, promptForParent, repo))
	hasOrigin := ec.Bool(repo.Backend.HasOrigin())
	shouldNewBranchPush := ec.Bool(repo.Config.ShouldNewBranchPush())
	isOffline := ec.Bool(repo.Config.IsOffline())
	mainBranch := repo.Config.MainBranch()
	if ec.Err == nil && hasOrigin && !isOffline {
		ec.Check(repo.Frontend.Fetch())
	}
	hasBranch := ec.Bool(repo.Backend.HasLocalOrOriginBranch(targetBranch, mainBranch))
	pushHook := ec.Bool(repo.Config.PushHook())
	if hasBranch {
		return nil, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	return &appendConfig{
		ancestorBranches:    []string{},
		targetBranch:        targetBranch,
		parentBranch:        parentBranch,
		hasOrigin:           hasOrigin,
		mainBranch:          mainBranch,
		shouldNewBranchPush: shouldNewBranchPush,
		noPushHook:          !pushHook,
		isOffline:           isOffline,
	}, ec.Err
}

func determineParentBranch(targetBranch string, promptForParent bool, repo *git.ProdRepo) (string, error) {
	if promptForParent {
		parentBranch, err := validate.EnterParent(targetBranch, repo.Config.MainBranch(), &repo.Backend)
		if err != nil {
			return "", err
		}
		err = validate.KnowsBranchAncestry(parentBranch, repo.Config.MainBranch(), &repo.Backend)
		if err != nil {
			return "", err
		}
		return parentBranch, nil
	}
	return repo.Config.MainBranch(), nil
}
