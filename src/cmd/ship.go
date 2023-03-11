package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/git-town/git-town/v7/src/validate"
	. "github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func shipCmd(repo *git.ProdRepo) *cobra.Command {
	var commitMessage string
	shipCmd := cobra.Command{
		Use:     "ship",
		GroupID: "basic",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: Ensure(repo, HasGitVersion, IsRepository, IsConfigured),
		Short:   "Deliver a completed feature branch",
		Long: fmt.Sprintf(`Deliver a completed feature branch

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
and Git Town will leave it up to your origin server to delete the remote branch.`, config.GithubTokenKey, config.ShipDeleteRemoteBranchKey),
		RunE: func(cmd *cobra.Command, args []string) error {
			connector, err := hosting.NewConnector(&repo.Config, &repo.Silent, cli.PrintConnectorAction)
			if err != nil {
				return err
			}
			config, err := determineShipConfig(args, connector, repo)
			if err != nil {
				return err
			}
			stepList, err := shipStepList(config, commitMessage, repo)
			if err != nil {
				return err
			}
			runState := runstate.New("ship", stepList)
			return runstate.Execute(runState, repo, connector)
		},
	}
	shipCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Specify the commit message for the squash commit")
	return &shipCmd
}

type shipConfig struct {
	branchToShip             string
	branchToMergeInto        string // TODO: rename to parentBranch
	canShipViaAPI            bool
	childBranches            []string
	defaultProposalMessage   string // TODO: rename to proposalMessage
	deleteOriginBranch       bool
	hasOrigin                bool
	hasTrackingBranch        bool
	initialBranch            string
	isShippingInitialBranch  bool
	isOffline                bool
	proposal                 *hosting.Proposal
	proposalsOfChildBranches []hosting.Proposal
}

func determineShipConfig(args []string, connector hosting.Connector, repo *git.ProdRepo) (*shipConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return nil, err
	}
	deleteOrigin, err := repo.Config.ShouldShipDeleteOriginBranch()
	if err != nil {
		return nil, err
	}
	var branchToShip string
	if len(args) > 0 {
		branchToShip = args[0]
	} else {
		branchToShip = initialBranch
	}
	isShippingInitialBranch := branchToShip == initialBranch
	if isShippingInitialBranch {
		hasOpenChanges, err := repo.Silent.HasOpenChanges()
		if err != nil {
			return nil, err
		}
		if hasOpenChanges {
			return nil, fmt.Errorf("you have uncommitted changes. Did you mean to commit them before shipping?")
		}
	}
	if hasOrigin && !isOffline {
		err := repo.Logging.Fetch()
		if err != nil {
			return nil, err
		}
	}
	if !isShippingInitialBranch {
		hasBranch, err := repo.Silent.HasLocalOrOriginBranch(branchToShip)
		if err != nil {
			return nil, err
		}
		if !hasBranch {
			return nil, fmt.Errorf("there is no branch named %q", branchToShip)
		}
	}
	if !repo.Config.IsFeatureBranch(branchToShip) {
		return nil, fmt.Errorf("the branch %q is not a feature branch. Only feature branches can be shipped", branchToShip)
	}
	err = validate.KnowsParentBranches([]string{branchToShip}, repo)
	if err != nil {
		return nil, err
	}
	err = ensureParentBranchIsMainOrPerennialBranch(branchToShip, repo)
	if err != nil {
		return nil, err
	}
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchToShip)
	if err != nil {
		return nil, err
	}
	branchToMergeInto := repo.Config.ParentBranch(branchToShip)
	canShipViaAPI := false
	defaultProposalMessage := ""
	var proposal *hosting.Proposal
	childBranches := repo.Config.ChildBranches(branchToShip)
	proposalsOfChildBranches := []hosting.Proposal{}
	if !isOffline && connector != nil {
		if hasTrackingBranch {
			proposal, err = connector.FindProposal(branchToShip, branchToMergeInto)
			if err != nil {
				return nil, err
			}
			if proposal != nil {
				canShipViaAPI = true
				defaultProposalMessage = connector.DefaultProposalMessage(*proposal)
			}
		}
		for _, childBranch := range childBranches {
			childProposal, err := connector.FindProposal(childBranch, branchToShip)
			if err != nil {
				return nil, fmt.Errorf("cannot determine proposal for branch %q: %w", branchToShip, err)
			}
			if childProposal != nil {
				proposalsOfChildBranches = append(proposalsOfChildBranches, *childProposal)
			}
		}
	}
	return &shipConfig{
		branchToMergeInto:        branchToMergeInto,
		branchToShip:             branchToShip,
		canShipViaAPI:            canShipViaAPI,
		childBranches:            childBranches,
		defaultProposalMessage:   defaultProposalMessage,
		deleteOriginBranch:       deleteOrigin,
		hasOrigin:                hasOrigin,
		hasTrackingBranch:        hasTrackingBranch,
		initialBranch:            initialBranch,
		isOffline:                isOffline,
		isShippingInitialBranch:  isShippingInitialBranch,
		proposal:                 proposal,
		proposalsOfChildBranches: proposalsOfChildBranches,
	}, nil
}

func ensureParentBranchIsMainOrPerennialBranch(branch string, repo *git.ProdRepo) error {
	parentBranch := repo.Config.ParentBranch(branch)
	if !repo.Config.IsMainBranch(parentBranch) && !repo.Config.IsPerennialBranch(parentBranch) {
		ancestors := repo.Config.AncestorBranches(branch)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		return fmt.Errorf(`shipping this branch would ship %q as well,
please ship %q first`, strings.Join(ancestorsWithoutMainOrPerennial, ", "), oldestAncestor)
	}
	return nil
}

func shipStepList(config *shipConfig, commitMessage string, repo *git.ProdRepo) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	updateBranchSteps(&list, config.branchToMergeInto, true, repo) // sync the parent branch
	updateBranchSteps(&list, config.branchToShip, false, repo)     // sync the branch to ship locally only
	list.Add(&steps.EnsureHasShippableChangesStep{Branch: config.branchToShip})
	list.Add(&steps.CheckoutStep{Branch: config.branchToMergeInto})
	if config.canShipViaAPI {
		// update the proposals of child branches
		for _, childProposal := range config.proposalsOfChildBranches {
			list.Add(&steps.UpdateProposalTargetStep{
				ProposalNumber: childProposal.Number,
				NewTarget:      config.branchToMergeInto,
				ExistingTarget: childProposal.Target,
			})
		}
		// push
		list.Add(&steps.PushBranchStep{Branch: config.branchToShip})
		list.Add(&steps.ConnectorMergeProposalStep{
			Branch:                 config.branchToShip,
			ProposalNumber:         config.proposal.Number,
			CommitMessage:          commitMessage,
			DefaultProposalMessage: config.defaultProposalMessage,
		})
		list.Add(&steps.PullBranchStep{})
	} else {
		list.Add(&steps.SquashMergeStep{Branch: config.branchToShip, CommitMessage: commitMessage, Parent: config.branchToMergeInto})
	}
	if config.hasOrigin && !config.isOffline {
		list.Add(&steps.PushBranchStep{Branch: config.branchToMergeInto, Undoable: true})
	}
	// NOTE: when shipping via API, we can always delete the remote branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via API)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if config.canShipViaAPI || (config.hasTrackingBranch && len(config.childBranches) == 0 && !config.isOffline) {
		if config.deleteOriginBranch {
			list.Add(&steps.DeleteOriginBranchStep{Branch: config.branchToShip, IsTracking: true})
		}
	}
	list.Add(&steps.DeleteLocalBranchStep{Branch: config.branchToShip})
	list.Add(&steps.DeleteParentBranchStep{Branch: config.branchToShip})
	for _, child := range config.childBranches {
		list.Add(&steps.SetParentStep{Branch: child, ParentBranch: config.branchToMergeInto})
	}
	if !config.isShippingInitialBranch {
		// TODO: check out the main branch here?
		list.Add(&steps.CheckoutStep{Branch: config.initialBranch})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: !config.isShippingInitialBranch}, repo)
	return list.Result()
}
