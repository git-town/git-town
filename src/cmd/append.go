package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func appendCmd() *cobra.Command {
	debug := false
	cmd := &cobra.Command{
		Use:     "append <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a new feature branch as a child of the current branch",
		Long: `Creates a new feature branch as a direct child of the current branch.

Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the origin repository
(if and only if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAppend(debug, args)
		},
	}
	debugFlag(cmd, &debug)
	return cmd
}

func runAppend(debug bool, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
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
	config, err := determineAppendConfig(args, &repo)
	if err != nil {
		return err
	}
	stepList, err := appendStepList(config, &repo)
	if err != nil {
		return err
	}
	runState := runstate.New("append", stepList)
	return runstate.Execute(runState, &repo, nil)
}

type appendConfig struct {
	ancestorBranches    []string
	hasOrigin           bool
	isOffline           bool
	mainBranch          string
	noPushHook          bool
	parentBranch        string
	shouldNewBranchPush bool
	targetBranch        string
}

func determineAppendConfig(args []string, repo *git.PublicRepo) (*appendConfig, error) {
	ec := runstate.ErrorChecker{}
	parentBranch := ec.String(repo.CurrentBranch())
	hasOrigin := ec.Bool(repo.HasOrigin())
	isOffline := ec.Bool(repo.Config.IsOffline())
	mainBranch := repo.Config.MainBranch()
	pushHook := ec.Bool(repo.Config.PushHook())
	shouldNewBranchPush := ec.Bool(repo.Config.ShouldNewBranchPush())
	targetBranch := args[0]
	if ec.Err != nil {
		return nil, ec.Err
	}
	if hasOrigin && !isOffline {
		ec.Check(repo.Fetch())
	}
	hasTargetBranch := ec.Bool(repo.HasLocalOrOriginBranch(targetBranch, mainBranch))
	if hasTargetBranch {
		ec.Fail("a branch named %q already exists", targetBranch)
	}
	ec.Check(validate.KnowsBranchAncestry(parentBranch, repo.Config.MainBranch(), repo))
	ancestorBranches := repo.Config.AncestorBranches(parentBranch)
	return &appendConfig{
		ancestorBranches:    ancestorBranches,
		isOffline:           isOffline,
		hasOrigin:           hasOrigin,
		mainBranch:          mainBranch,
		noPushHook:          !pushHook,
		parentBranch:        parentBranch,
		shouldNewBranchPush: shouldNewBranchPush,
		targetBranch:        targetBranch,
	}, ec.Err
}

func appendStepList(config *appendConfig, repo *git.PublicRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range append(config.ancestorBranches, config.parentBranch) {
		updateBranchSteps(&list, branch, true, repo)
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &repo.InternalRepo, config.mainBranch)
	return list.Result()
}
