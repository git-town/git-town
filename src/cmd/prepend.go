package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func prependCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "prepend <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a new feature branch as the parent of the current branch",
		Long: `Creates a new feature branch as the parent of the current branch

Syncs the parent branch,
cuts a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the origin repository
(if "push-new-branches" is true),
and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrepend(debug, args)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runPrepend(debug bool, args []string) error {
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
	config, err := determinePrependConfig(args, &repo)
	if err != nil {
		return err
	}
	stepList, err := prependStepList(config, &repo)
	if err != nil {
		return err
	}
	runState := runstate.New("prepend", stepList)
	return runstate.Execute(runState, &repo, nil)
}

type prependConfig struct {
	ancestorBranches    []string
	hasOrigin           bool
	initialBranch       string
	isOffline           bool
	mainBranch          string
	noPushHook          bool
	parentBranch        string
	shouldNewBranchPush bool
	targetBranch        string
}

func determinePrependConfig(args []string, repo *git.PublicRepo) (*prependConfig, error) {
	ec := runstate.ErrorChecker{}
	initialBranch := ec.String(repo.CurrentBranch())
	hasOrigin := ec.Bool(repo.HasOrigin())
	shouldNewBranchPush := ec.Bool(repo.Config.ShouldNewBranchPush())
	pushHook := ec.Bool(repo.Config.PushHook())
	isOffline := ec.Bool(repo.Config.IsOffline())
	mainBranch := repo.Config.MainBranch()
	if ec.Err != nil {
		return nil, ec.Err
	}
	if hasOrigin && !isOffline {
		err := repo.Fetch()
		if err != nil {
			return nil, err
		}
	}
	targetBranch := args[0]
	hasBranch, err := repo.HasLocalOrOriginBranch(targetBranch, mainBranch)
	if err != nil {
		return nil, err
	}
	if hasBranch {
		return nil, fmt.Errorf("a branch named %q already exists", targetBranch)
	}
	if !repo.Config.IsFeatureBranch(initialBranch) {
		return nil, fmt.Errorf("the branch %q is not a feature branch. Only feature branches can have parent branches", initialBranch)
	}
	err = validate.KnowsBranchAncestry(initialBranch, repo.Config.MainBranch(), repo)
	if err != nil {
		return nil, err
	}
	return &prependConfig{
		hasOrigin:           hasOrigin,
		initialBranch:       initialBranch,
		isOffline:           isOffline,
		mainBranch:          mainBranch,
		noPushHook:          !pushHook,
		parentBranch:        repo.Config.ParentBranch(initialBranch),
		ancestorBranches:    repo.Config.AncestorBranches(initialBranch),
		shouldNewBranchPush: shouldNewBranchPush,
		targetBranch:        targetBranch,
	}, nil
}

func prependStepList(config *prependConfig, repo *git.PublicRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	for _, branch := range config.ancestorBranches {
		updateBranchSteps(&list, branch, true, repo)
	}
	list.Add(&steps.CreateBranchStep{Branch: config.targetBranch, StartingPoint: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.targetBranch, ParentBranch: config.parentBranch})
	list.Add(&steps.SetParentStep{Branch: config.initialBranch, ParentBranch: config.targetBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch})
	if config.hasOrigin && config.shouldNewBranchPush && !config.isOffline {
		list.Add(&steps.CreateTrackingBranchStep{Branch: config.targetBranch, NoPushHook: config.noPushHook})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: true}, &repo.InternalRepo, config.mainBranch)
	return list.Result()
}
