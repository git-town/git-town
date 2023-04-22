package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v8/src/cli"
	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/src/execute"
	"github.com/git-town/git-town/v8/src/flags"
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
	"github.com/git-town/git-town/v8/src/runstate"
	"github.com/git-town/git-town/v8/src/steps"
	"github.com/git-town/git-town/v8/src/validate"
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
	connector, err := hosting.NewConnector(run.Config.GitTown, &run.Backend, cli.PrintConnectorAction)
	if err != nil {
		return err
	}
	config, err := determineShipConfig(args, connector, &run)
	if err != nil {
		return err
	}
	stepList, err := shipStepList(config, message, &run)
	if err != nil {
		return err
	}
	runState := runstate.New("ship", stepList)
	return runstate.Execute(runState, &run, connector)
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
	mainBranch               string
	proposal                 *hosting.Proposal
	proposalsOfChildBranches []hosting.Proposal
}

func determineShipConfig(args []string, connector hosting.Connector, run *git.ProdRunner) (*shipConfig, error) {
	hasOrigin, err := run.Backend.HasOrigin()
	if err != nil {
		return nil, err
	}
	isOffline, err := run.Config.IsOffline()
	if err != nil {
		return nil, err
	}
	initialBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return nil, err
	}
	deleteOrigin, err := run.Config.ShouldShipDeleteOriginBranch()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	var branchToShip string
	if len(args) > 0 {
		branchToShip = args[0]
	} else {
		branchToShip = initialBranch
	}
	isShippingInitialBranch := branchToShip == initialBranch
	if isShippingInitialBranch {
		hasOpenChanges, err := run.Backend.HasOpenChanges()
		if err != nil {
			return nil, err
		}
		if hasOpenChanges {
			return nil, fmt.Errorf("you have uncommitted changes. Did you mean to commit them before shipping?")
		}
	}
	if hasOrigin && !isOffline {
		err := run.Frontend.Fetch()
		if err != nil {
			return nil, err
		}
	}
	if !isShippingInitialBranch {
		hasBranch, err := run.Backend.HasLocalOrOriginBranch(branchToShip, mainBranch)
		if err != nil {
			return nil, err
		}
		if !hasBranch {
			return nil, fmt.Errorf("there is no branch named %q", branchToShip)
		}
	}
	if !run.Config.IsFeatureBranch(branchToShip) {
		return nil, fmt.Errorf("the branch %q is not a feature branch. Only feature branches can be shipped", branchToShip)
	}
	err = validate.KnowsBranchAncestry(branchToShip, mainBranch, &run.Backend)
	if err != nil {
		return nil, err
	}
	err = ensureParentBranchIsMainOrPerennialBranch(branchToShip, run)
	if err != nil {
		return nil, err
	}
	hasTrackingBranch, err := run.Backend.HasTrackingBranch(branchToShip)
	if err != nil {
		return nil, err
	}
	branchToMergeInto := run.Config.ParentBranch(branchToShip)
	canShipViaAPI := false
	defaultProposalMessage := ""
	var proposal *hosting.Proposal
	childBranches := run.Config.ChildBranches(branchToShip)
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
		mainBranch:               mainBranch,
		proposal:                 proposal,
		proposalsOfChildBranches: proposalsOfChildBranches,
	}, nil
}

func ensureParentBranchIsMainOrPerennialBranch(branch string, run *git.ProdRunner) error {
	parentBranch := run.Config.ParentBranch(branch)
	if !run.Config.IsMainBranch(parentBranch) && !run.Config.IsPerennialBranch(parentBranch) {
		ancestors := run.Config.AncestorBranches(branch)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		return fmt.Errorf(`shipping this branch would ship %q as well,
please ship %q first`, strings.Join(ancestorsWithoutMainOrPerennial, ", "), oldestAncestor)
	}
	return nil
}

func shipStepList(config *shipConfig, commitMessage string, run *git.ProdRunner) (runstate.StepList, error) {
	list := runstate.StepListBuilder{}
	updateBranchSteps(&list, config.branchToMergeInto, true, run) // sync the parent branch
	updateBranchSteps(&list, config.branchToShip, false, run)     // sync the branch to ship locally only
	list.Add(&steps.EnsureHasShippableChangesStep{Branch: config.branchToShip, Parent: config.mainBranch})
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
	list.Add(&steps.DeleteLocalBranchStep{Branch: config.branchToShip, Parent: config.mainBranch})
	list.Add(&steps.DeleteParentBranchStep{Branch: config.branchToShip})
	for _, child := range config.childBranches {
		list.Add(&steps.SetParentStep{Branch: child, ParentBranch: config.branchToMergeInto})
	}
	if !config.isShippingInitialBranch {
		// TODO: check out the main branch here?
		list.Add(&steps.CheckoutStep{Branch: config.initialBranch})
	}
	list.Wrap(runstate.WrapOptions{RunInGitRoot: true, StashOpenChanges: !config.isShippingInitialBranch}, &run.Backend, config.mainBranch)
	return list.Result()
}
