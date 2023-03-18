package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func hackCmd() *cobra.Command {
	debug := false
	promptForParentFlag := false
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: "basic",
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a new feature branch off the main development branch",
		Long: `Creates a new feature branch off the main development branch

Syncs the main branch,
forks a new feature branch with the given name off the main branch,
pushes the new feature branch to origin
(if and only if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHack(debug, promptForParentFlag, args)
		},
	}
	cmd.Flags().BoolVarP(&promptForParentFlag, "prompt", "p", false, "Prompt for the parent branch")
	return &cmd
}

func runHack(debug, prompt bool, args []string) error {
	repo, err := LoadPublicRepo(RepoArgs{
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
		validateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineHackConfig(args, prompt, &repo)
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

func determineHackConfig(args []string, promptForParent bool, repo *git.PublicRepo) (*appendConfig, error) {
	ec := runstate.ErrorChecker{}
	targetBranch := args[0]
	parentBranch := ec.String(determineParentBranch(targetBranch, promptForParent, repo))
	hasOrigin := ec.Bool(repo.HasOrigin())
	shouldNewBranchPush := ec.Bool(repo.Config.ShouldNewBranchPush())
	isOffline := ec.Bool(repo.Config.IsOffline())
	mainBranch := repo.Config.MainBranch()
	if ec.Err == nil && hasOrigin && !isOffline {
		ec.Check(repo.Fetch())
	}
	hasBranch := ec.Bool(repo.HasLocalOrOriginBranch(targetBranch, mainBranch))
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

func determineParentBranch(targetBranch string, promptForParent bool, repo *git.PublicRepo) (string, error) {
	if promptForParent {
		parentBranch, err := validate.EnterParent(targetBranch, repo.Config.MainBranch(), repo)
		if err != nil {
			return "", err
		}
		err = validate.KnowsBranchAncestry(parentBranch, repo.Config.MainBranch(), repo)
		if err != nil {
			return "", err
		}
		return parentBranch, nil
	}
	return repo.Config.MainBranch(), nil
}
