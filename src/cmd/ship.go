package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/steps"
	"github.com/git-town/git-town/v9/src/stringslice"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const shipDesc = "Deliver a completed feature branch"

const shipHelp = `
Squash-merges the current branch, or <branch_name> if given,
into the main branch, resulting in linear history on the main branch.

- syncs the main branch
- pulls updates for <branch_name>
- merges the main branch into <branch_name>
- squash-merges <branch_name> into the main branch
  with commit message specified by the user
- pushes the main branch to the origin repository
- deletes <branch_name> from the local and origin repositories

Ships direct children of the main branch.
To ship a nested child branch, ship or kill all ancestor branches first.

If you use GitHub, this command can squash merge pull requests via the GitHub API. Setup:
1. Get a GitHub personal access token with the "repo" scope
2. Run 'git config %s <token>' (optionally add the '--global' flag)
Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API.
It will also update the base branch for any pull requests against that branch.

If your origin server deletes shipped branches, for example
GitHub's feature to automatically delete head branches,
run "git config %s false"
and Git Town will leave it up to your origin server to delete the remote branch.`

func shipCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	addMessageFlag, readMessageFlag := flags.String("message", "m", "", "Specify the commit message for the squash commit")
	cmd := cobra.Command{
		Use:     "ship",
		GroupID: "basic",
		Args:    cobra.MaximumNArgs(1),
		Short:   shipDesc,
		Long:    long(shipDesc, fmt.Sprintf(shipHelp, config.GithubTokenKey, config.ShipDeleteRemoteBranchKey)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ship(args, readMessageFlag(cmd), readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	addMessageFlag(&cmd)
	return &cmd
}

func ship(args []string, message string, debug bool) error {
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 true,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: len(args) == 0,
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
	connector, err := hosting.NewConnector(repo.Runner.Config.GitTown, &repo.Runner.Backend, cli.PrintConnectorAction)
	if err != nil {
		return err
	}
	config, err := determineShipConfig(args, connector, &repo.Runner, allBranches, initialBranch, repo.IsOffline)
	if err != nil {
		return err
	}
	if config.branchToShip.Name == initialBranch {
		err = validate.NoOpenChanges(&repo.Runner.Backend)
		if err != nil {
			return err
		}
	}
	stepList, err := shipStepList(config, message, &repo.Runner)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:     "ship",
		RunStepList: stepList,
	}
	return runstate.Execute(&runState, &repo.Runner, connector, repo.RootDir)
}

type shipConfig struct {
	branchToShip             git.BranchSyncStatus
	targetBranch             git.BranchSyncStatus
	canShipViaAPI            bool
	childBranches            []string
	proposalMessage          string
	deleteOriginBranch       bool
	hasOrigin                bool
	hasUpstream              bool
	initialBranch            string
	isShippingInitialBranch  bool
	isOffline                bool
	mainBranch               string
	proposal                 *hosting.Proposal
	proposalsOfChildBranches []hosting.Proposal
	pullBranchStrategy       config.PullBranchStrategy
	pushHook                 bool
	shouldSyncUpstream       bool
	syncStrategy             config.SyncStrategy
}

func determineShipConfig(args []string, connector hosting.Connector, run *git.ProdRunner, allBranches git.BranchesSyncStatus, initialBranch string, isOffline bool) (*shipConfig, error) {
	hasOrigin, err := run.Backend.HasOrigin()
	if err != nil {
		return nil, err
	}
	deleteOrigin, err := run.Config.ShouldShipDeleteOriginBranch()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	branchNameToShip := determineBranchToShip(args, initialBranch)
	branchToShip := allBranches.Lookup(branchNameToShip)
	isShippingInitialBranch := branchNameToShip == initialBranch
	syncStrategy, err := run.Config.SyncStrategy()
	if err != nil {
		return nil, err
	}
	pullBranchStrategy, err := run.Config.PullBranchStrategy()
	if err != nil {
		return nil, err
	}
	shouldSyncUpstream, err := run.Config.ShouldSyncUpstream()
	if err != nil {
		return nil, err
	}
	if !isShippingInitialBranch {
		if branchToShip == nil {
			return nil, fmt.Errorf(messages.BranchDoesntExist, branchNameToShip)
		}
	}
	if !run.Config.IsFeatureBranch(branchNameToShip) {
		return nil, fmt.Errorf(messages.ShipNoFeatureBranch, branchNameToShip)
	}
	err = validate.KnowsBranchAncestors(branchNameToShip, mainBranch, &run.Backend)
	if err != nil {
		return nil, err
	}
	err = ensureParentBranchIsMainOrPerennialBranch(branchNameToShip, run)
	if err != nil {
		return nil, err
	}
	hasUpstream, err := run.Backend.HasUpstream()
	if err != nil {
		return nil, err
	}
	lineage := run.Config.Lineage()
	targetBranchName := lineage.Parent(branchNameToShip)
	targetBranch := allBranches.Lookup(targetBranchName)
	if targetBranch == nil {
		return nil, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	canShipViaAPI := false
	proposalMessage := ""
	var proposal *hosting.Proposal
	childBranches := lineage.Children(branchNameToShip)
	proposalsOfChildBranches := []hosting.Proposal{}
	pushHook, err := run.Config.PushHook()
	if err != nil {
		return nil, err
	}
	if !isOffline && connector != nil {
		if branchToShip.HasTrackingBranch() {
			proposal, err = connector.FindProposal(branchNameToShip, targetBranchName)
			if err != nil {
				return nil, err
			}
			if proposal != nil {
				canShipViaAPI = true
				proposalMessage = connector.DefaultProposalMessage(*proposal)
			}
		}
		for _, childBranch := range childBranches {
			childProposal, err := connector.FindProposal(childBranch, branchNameToShip)
			if err != nil {
				return nil, fmt.Errorf(messages.ProposalNotFoundForBranch, branchNameToShip, err)
			}
			if childProposal != nil {
				proposalsOfChildBranches = append(proposalsOfChildBranches, *childProposal)
			}
		}
	}
	return &shipConfig{
		targetBranch:             *targetBranch,
		branchToShip:             *branchToShip,
		canShipViaAPI:            canShipViaAPI,
		childBranches:            childBranches,
		proposalMessage:          proposalMessage,
		deleteOriginBranch:       deleteOrigin,
		hasOrigin:                hasOrigin,
		hasUpstream:              hasUpstream,
		initialBranch:            initialBranch,
		isOffline:                isOffline,
		isShippingInitialBranch:  isShippingInitialBranch,
		mainBranch:               mainBranch,
		proposal:                 proposal,
		proposalsOfChildBranches: proposalsOfChildBranches,
		pullBranchStrategy:       pullBranchStrategy,
		pushHook:                 pushHook,
		shouldSyncUpstream:       shouldSyncUpstream,
		syncStrategy:             syncStrategy,
	}, nil
}

func ensureParentBranchIsMainOrPerennialBranch(branch string, run *git.ProdRunner) error {
	lineage := run.Config.Lineage()
	parentBranch := lineage.Parent(branch)
	if !run.Config.IsMainBranch(parentBranch) && !run.Config.IsPerennialBranch(parentBranch) {
		ancestors := lineage.Ancestors(branch)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		return fmt.Errorf(`shipping this branch would ship %s as well,
please ship %q first`, stringslice.Connect(ancestorsWithoutMainOrPerennial), oldestAncestor)
	}
	return nil
}

func shipStepList(config *shipConfig, commitMessage string, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	// sync the parent branch
	updateBranchSteps(&list, updateBranchStepsArgs{
		branch:             config.targetBranch,
		hasOrigin:          config.hasOrigin,
		hasUpstream:        config.hasUpstream,
		isOffline:          config.isOffline,
		mainBranch:         config.mainBranch,
		pullBranchStrategy: config.pullBranchStrategy,
		pushBranch:         true,
		pushHook:           config.pushHook,
		run:                run,
		shouldSyncUpstream: config.shouldSyncUpstream,
		syncStrategy:       config.syncStrategy,
	})
	// sync the branch to ship locally only
	updateBranchSteps(&list, updateBranchStepsArgs{
		branch:             config.branchToShip,
		hasOrigin:          config.hasOrigin,
		hasUpstream:        config.hasUpstream,
		isOffline:          config.isOffline,
		mainBranch:         config.mainBranch,
		pullBranchStrategy: config.pullBranchStrategy,
		pushBranch:         false,
		pushHook:           config.pushHook,
		run:                run,
		shouldSyncUpstream: config.shouldSyncUpstream,
		syncStrategy:       config.syncStrategy,
	})
	list.Add(&steps.EnsureHasShippableChangesStep{Branch: config.branchToShip.Name, Parent: config.mainBranch})
	list.Add(&steps.CheckoutStep{Branch: config.targetBranch.Name})
	if config.canShipViaAPI {
		// update the proposals of child branches
		for _, childProposal := range config.proposalsOfChildBranches {
			list.Add(&steps.UpdateProposalTargetStep{
				ProposalNumber: childProposal.Number,
				NewTarget:      config.targetBranch.Name,
				ExistingTarget: childProposal.Target,
			})
		}
		// push
		list.Add(&steps.PushBranchStep{Branch: config.branchToShip.Name})
		list.Add(&steps.ConnectorMergeProposalStep{
			Branch:          config.branchToShip.Name,
			ProposalNumber:  config.proposal.Number,
			CommitMessage:   commitMessage,
			ProposalMessage: config.proposalMessage,
		})
		list.Add(&steps.PullBranchStep{})
	} else {
		list.Add(&steps.SquashMergeStep{Branch: config.branchToShip.Name, CommitMessage: commitMessage, Parent: config.targetBranch.Name})
	}
	if config.hasOrigin && !config.isOffline {
		list.Add(&steps.PushBranchStep{Branch: config.targetBranch.Name, Undoable: true})
	}
	// NOTE: when shipping via API, we can always delete the remote branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via API)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if config.canShipViaAPI || (config.branchToShip.HasTrackingBranch() && len(config.childBranches) == 0 && !config.isOffline) {
		if config.deleteOriginBranch {
			list.Add(&steps.DeleteOriginBranchStep{Branch: config.branchToShip.Name, IsTracking: true})
		}
	}
	list.Add(&steps.DeleteLocalBranchStep{Branch: config.branchToShip.Name, Parent: config.mainBranch})
	list.Add(&steps.DeleteParentBranchStep{Branch: config.branchToShip.Name, Parent: run.Config.Lineage().Parent(config.branchToShip.Name)})
	for _, child := range config.childBranches {
		list.Add(&steps.SetParentStep{Branch: child, ParentBranch: config.targetBranch.Name})
	}
	if !config.isShippingInitialBranch {
		list.Add(&steps.CheckoutStep{Branch: config.initialBranch})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: !config.isShippingInitialBranch}, &run.Backend, config.mainBranch)
	return list.Result()
}

func determineBranchToShip(args []string, initialBranch string) string {
	if len(args) > 0 {
		return args[0]
	}
	return initialBranch
}
