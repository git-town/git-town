package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/slice"
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
	repo, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := determineShipConfig(args, &repo)
	if err != nil || exit {
		return err
	}
	if config.branchToShip.Name == config.branches.Initial {
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
		Branches:  config.branches.All,
	})
}

type shipConfig struct {
	branches                 domain.Branches
	branchToShip             domain.BranchInfo
	connector                hosting.Connector
	targetBranch             domain.BranchInfo
	canShipViaAPI            bool
	childBranches            domain.LocalBranchNames
	proposalMessage          string
	deleteOriginBranch       bool
	hasOpenChanges           bool
	remotes                  domain.Remotes
	isShippingInitialBranch  bool
	isOffline                bool
	lineage                  config.Lineage
	mainBranch               domain.LocalBranchName
	previousBranch           domain.LocalBranchName
	proposal                 *hosting.Proposal
	proposalsOfChildBranches []hosting.Proposal
	pullBranchStrategy       config.PullBranchStrategy
	pushHook                 bool
	shouldSyncUpstream       bool
	syncStrategy             config.SyncStrategy
}

func determineShipConfig(args []string, repo *execute.RepoData) (*shipConfig, bool, error) {
	branches, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: len(args) == 0,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	hasOpenChanges, err := repo.Runner.Backend.HasOpenChanges()
	if err != nil {
		return nil, false, err
	}
	remotes, err := repo.Runner.Backend.Remotes()
	if err != nil {
		return nil, false, err
	}
	deleteOrigin, err := repo.Runner.Config.ShouldShipDeleteOriginBranch()
	if err != nil {
		return nil, false, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	branchNameToShip := domain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	branchToShip := branches.All.FindLocalBranch(branchNameToShip)
	isShippingInitialBranch := branchNameToShip == branches.Initial
	syncStrategy, err := repo.Runner.Config.SyncStrategy()
	if err != nil {
		return nil, false, err
	}
	pullBranchStrategy, err := repo.Runner.Config.PullBranchStrategy()
	if err != nil {
		return nil, false, err
	}
	shouldSyncUpstream, err := repo.Runner.Config.ShouldSyncUpstream()
	if err != nil {
		return nil, false, err
	}
	if !isShippingInitialBranch {
		if branchToShip == nil {
			return nil, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToShip)
		}
	}
	if !branches.Types.IsFeatureBranch(branchNameToShip) {
		return nil, false, fmt.Errorf(messages.ShipNoFeatureBranch, branchNameToShip)
	}
	lineage := repo.Runner.Config.Lineage()
	updated, err := validate.KnowsBranchAncestors(branchNameToShip, validate.KnowsBranchAncestorsArgs{
		DefaultBranch: mainBranch,
		Backend:       &repo.Runner.Backend,
		AllBranches:   branches.All,
		Lineage:       lineage,
		BranchTypes:   branches.Types,
		MainBranch:    mainBranch,
	})
	if err != nil {
		return nil, false, err
	}
	if updated {
		lineage = repo.Runner.Config.Lineage()
	}
	err = ensureParentBranchIsMainOrPerennialBranch(branchNameToShip, branches.Types, lineage)
	if err != nil {
		return nil, false, err
	}
	targetBranchName := lineage.Parent(branchNameToShip)
	targetBranch := branches.All.FindLocalBranch(targetBranchName)
	if targetBranch == nil {
		return nil, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	canShipViaAPI := false
	proposalMessage := ""
	var proposal *hosting.Proposal
	childBranches := lineage.Children(branchNameToShip)
	proposalsOfChildBranches := []hosting.Proposal{}
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, false, err
	}
	originURL := repo.Runner.Config.OriginURL()
	hostingService, err := repo.Runner.Config.HostingService()
	if err != nil {
		return nil, false, err
	}
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetShaForBranch: repo.Runner.Backend.ShaForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.Config.GiteaToken(),
		GithubAPIToken:  repo.Runner.Config.GitHubToken(),
		GitlabAPIToken:  repo.Runner.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             cli.PrintingLog{},
	})
	if err != nil {
		return nil, false, err
	}
	if !repo.IsOffline && connector != nil {
		if branchToShip.HasTrackingBranch() {
			proposal, err = connector.FindProposal(branchNameToShip, targetBranchName)
			if err != nil {
				return nil, false, err
			}
			if proposal != nil {
				canShipViaAPI = true
				proposalMessage = connector.DefaultProposalMessage(*proposal)
			}
		}
		for _, childBranch := range childBranches {
			childProposal, err := connector.FindProposal(childBranch, branchNameToShip)
			if err != nil {
				return nil, false, fmt.Errorf(messages.ProposalNotFoundForBranch, branchNameToShip, err)
			}
			if childProposal != nil {
				proposalsOfChildBranches = append(proposalsOfChildBranches, *childProposal)
			}
		}
	}
	return &shipConfig{
		branches:                 branches,
		connector:                connector,
		targetBranch:             *targetBranch,
		branchToShip:             *branchToShip,
		canShipViaAPI:            canShipViaAPI,
		childBranches:            childBranches,
		proposalMessage:          proposalMessage,
		deleteOriginBranch:       deleteOrigin,
		hasOpenChanges:           hasOpenChanges,
		remotes:                  remotes,
		isOffline:                repo.IsOffline,
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
	}, false, nil
}

func ensureParentBranchIsMainOrPerennialBranch(branch domain.LocalBranchName, branchTypes domain.BranchTypes, lineage config.Lineage) error {
	parentBranch := lineage.Parent(branch)
	if !branchTypes.IsMainBranch(parentBranch) && !branchTypes.IsPerennialBranch(parentBranch) {
		ancestors := lineage.Ancestors(branch)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		return fmt.Errorf(`shipping this branch would ship %s as well,
please ship %q first`, stringslice.Connect(ancestorsWithoutMainOrPerennial.Strings()), oldestAncestor)
	}
	return nil
}

func shipStepList(config *shipConfig, commitMessage string, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	// sync the parent branch
	syncBranchSteps(&list, syncBranchStepsArgs{
		branch:             config.targetBranch,
		branchTypes:        config.branches.Types,
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
	// sync the branch to ship (local sync only)
	syncBranchSteps(&list, syncBranchStepsArgs{
		branch:             config.branchToShip,
		branchTypes:        config.branches.Types,
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
		list.Add(&steps.PushBranchStep{Branch: config.branchToShip.Name, Remote: config.branchToShip.Remote()})
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
		list.Add(&steps.PushBranchStep{Branch: config.targetBranch.Name, Remote: config.targetBranch.Remote(), Undoable: true})
	}
	// NOTE: when shipping via API, we can always delete the remote branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via API)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if config.canShipViaAPI || (config.branchToShip.HasTrackingBranch() && len(config.childBranches) == 0 && !config.isOffline) {
		if config.deleteOriginBranch {
			remote, remoteBranch := config.branchToShip.RemoteName.Parts()
			list.Add(&steps.DeleteRemoteBranchStep{Branch: remoteBranch, Remote: remote, IsTracking: true})
		}
	}
	list.Add(&steps.DeleteLocalBranchStep{Branch: config.branchToShip.Name, Parent: config.mainBranch.Location()})
	list.Add(&steps.DeleteParentBranchStep{Branch: config.branchToShip.Name, Parent: run.Config.Lineage().Parent(config.branchToShip.Name)})
	for _, child := range config.childBranches {
		list.Add(&steps.SetParentStep{Branch: child, ParentBranch: config.targetBranch.Name})
	}
	if !config.isShippingInitialBranch {
		list.Add(&steps.CheckoutStep{Branch: config.branches.Initial})
	}
	list.Wrap(runstate.WrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: !config.isShippingInitialBranch && config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return list.Result()
}
