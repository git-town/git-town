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
		Long:    long(shipDesc, fmt.Sprintf(shipHelp, config.KeyGithubToken, config.KeyShipDeleteRemoteBranch)),
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
	config, err := determineShipConfig(args, &repo.Runner, repo.IsOffline)
	if err != nil {
		return err
	}
	if config.branchToShip.Name == config.initialBranch {
		hasOpenChanges, err := repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
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
	return runstate.Execute(runstate.ExecuteArgs{
		RunState:  &runState,
		Run:       &repo.Runner,
		Connector: config.connector,
		RootDir:   repo.RootDir,
	})
}

type shipConfig struct {
	branchDurations          config.BranchDurations
	branchToShip             git.BranchSyncStatus
	connector                hosting.Connector
	targetBranch             git.BranchSyncStatus
	canShipViaAPI            bool
	childBranches            []string
	proposalMessage          string
	deleteOriginBranch       bool
	hasOpenChanges           bool
	remotes                  config.Remotes
	initialBranch            string
	isShippingInitialBranch  bool
	isOffline                bool
	lineage                  config.Lineage
	mainBranch               string
	previousBranch           string
	proposal                 *hosting.Proposal
	proposalsOfChildBranches []hosting.Proposal
	pullBranchStrategy       config.PullBranchStrategy
	pushHook                 bool
	shouldSyncUpstream       bool
	syncStrategy             config.SyncStrategy
}

func determineShipConfig(args []string, run *git.ProdRunner, isOffline bool) (*shipConfig, error) {
	branches, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return nil, err
	}
	previousBranch := run.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := run.Backend.HasOpenChanges()
	if err != nil {
		return nil, err
	}
	remotes, err := run.Backend.Remotes()
	if err != nil {
		return nil, err
	}
	deleteOrigin, err := run.Config.ShouldShipDeleteOriginBranch()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	branchNameToShip := stringslice.FirstElementOr(args, branches.Initial)
	branchToShip := branches.All.Lookup(branchNameToShip)
	isShippingInitialBranch := branchNameToShip == branches.Initial
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
	if !branches.Durations.IsFeatureBranch(branchNameToShip) {
		return nil, fmt.Errorf(messages.ShipNoFeatureBranch, branchNameToShip)
	}
	lineage := run.Config.Lineage()
	updated, err := validate.KnowsBranchAncestors(branchNameToShip, validate.KnowsBranchAncestorsArgs{
		DefaultBranch:   mainBranch,
		Backend:         &run.Backend,
		AllBranches:     branches.All,
		Lineage:         lineage,
		BranchDurations: branches.Durations,
		MainBranch:      mainBranch,
	})
	if err != nil {
		return nil, err
	}
	if updated {
		lineage = run.Config.Lineage()
	}
	err = ensureParentBranchIsMainOrPerennialBranch(branchNameToShip, branches.Durations, lineage)
	if err != nil {
		return nil, err
	}
	targetBranchName := lineage.Parent(branchNameToShip)
	targetBranch := branches.All.Lookup(targetBranchName)
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
	originURL := run.Config.OriginURL()
	hostingService, err := run.Config.HostingService()
	if err != nil {
		return nil, err
	}
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetShaForBranch: run.Backend.ShaForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   run.Config.GiteaToken(),
		GithubAPIToken:  run.Config.GitHubToken(),
		GitlabAPIToken:  run.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             cli.PrintingLog{},
	})
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
		branchDurations:          branches.Durations,
		connector:                connector,
		targetBranch:             *targetBranch,
		branchToShip:             *branchToShip,
		canShipViaAPI:            canShipViaAPI,
		childBranches:            childBranches,
		proposalMessage:          proposalMessage,
		deleteOriginBranch:       deleteOrigin,
		hasOpenChanges:           hasOpenChanges,
		remotes:                  remotes,
		initialBranch:            branches.Initial,
		isOffline:                isOffline,
		isShippingInitialBranch:  isShippingInitialBranch,
		lineage:                  lineage,
		mainBranch:               mainBranch,
		previousBranch:           previousBranch,
		proposal:                 proposal,
		proposalsOfChildBranches: proposalsOfChildBranches,
		pullBranchStrategy:       pullBranchStrategy,
		pushHook:                 pushHook,
		shouldSyncUpstream:       shouldSyncUpstream,
		syncStrategy:             syncStrategy,
	}, nil
}

func ensureParentBranchIsMainOrPerennialBranch(branch string, branchDurations config.BranchDurations, lineage config.Lineage) error {
	parentBranch := lineage.Parent(branch)
	if !branchDurations.IsMainBranch(parentBranch) && !branchDurations.IsPerennialBranch(parentBranch) {
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
	syncBranchSteps(&list, syncBranchStepsArgs{
		branch:             config.targetBranch,
		branchDurations:    config.branchDurations,
		remotes:            config.remotes,
		isOffline:          config.isOffline,
		lineage:            config.lineage,
		mainBranch:         config.mainBranch,
		pullBranchStrategy: config.pullBranchStrategy,
		pushBranch:         true,
		pushHook:           config.pushHook,
		shouldSyncUpstream: config.shouldSyncUpstream,
		syncStrategy:       config.syncStrategy,
	})
	// sync the branch to ship locally only
	syncBranchSteps(&list, syncBranchStepsArgs{
		branch:             config.branchToShip,
		branchDurations:    config.branchDurations,
		remotes:            config.remotes,
		isOffline:          config.isOffline,
		lineage:            config.lineage,
		mainBranch:         config.mainBranch,
		pullBranchStrategy: config.pullBranchStrategy,
		pushBranch:         false,
		pushHook:           config.pushHook,
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
	if config.remotes.HasOrigin() && !config.isOffline {
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
	list.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: !config.isShippingInitialBranch && config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.initialBranch,
		PreviousBranch:   config.previousBranch,
	})
	return list.Result()
}
