package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const shipCommand = "ship"

const shipDesc = "Deliver a completed feature branch"

const shipHelp = `
Squash-merges the current branch, or <branch_name> if given, into the main branch, resulting in linear history on the main branch.

- syncs the main branch
- pulls updates for <branch_name>
- merges the main branch into <branch_name>
- squash-merges <branch_name> into the main branch
  with commit message specified by the user
- pushes the main branch to the origin repository
- deletes <branch_name> from the local and origin repositories

Ships direct children of the main branch. To ship a child branch, ship or kill all ancestor branches first.

If you use GitHub, this command can squash merge pull requests via the GitHub API. Setup:

1. Get a GitHub personal access token with the "repo" scope
2. Run 'git config %s <token>' (optionally add the '--global' flag)

Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API. It will also update the base branch for any pull requests against that branch.

If your origin server deletes shipped branches, for example GitHub's feature to automatically delete head branches, run "git config %s false" and Git Town will leave it up to your origin server to delete the tracking branch of the branch you are shipping.`

func shipCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMessageFlag, readMessageFlag := flags.CommitMessage("Specify the commit message for the squash commit")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:   shipCommand,
		Args:  cobra.MaximumNArgs(1),
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, gitconfig.KeyGithubToken, gitconfig.KeyShipDeleteTrackingBranch)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeShip(args, readMessageFlag(cmd), readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	return &cmd
}

func executeShip(args []string, message Option[gitdomain.CommitMessage], dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineShipData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateData(*data)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               shipCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            shipProgram(data, message),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type shipData struct {
	allBranches              gitdomain.BranchInfos
	branchToShip             gitdomain.BranchInfo
	canShipViaAPI            bool
	childBranches            gitdomain.LocalBranchNames
	config                   config.ValidatedConfig
	connector                hostingdomain.Connector
	dialogTestInputs         components.TestInputs
	dryRun                   bool
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	isShippingInitialBranch  bool
	previousBranch           gitdomain.LocalBranchName
	proposal                 Option[hostingdomain.Proposal]
	proposalMessage          string
	proposalsOfChildBranches []hostingdomain.Proposal
	remotes                  gitdomain.Remotes
	targetBranch             gitdomain.BranchInfo
}

func determineShipData(args []string, repo execute.OpenRepoResult, dryRun, verbose bool) (*shipData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: len(args) == 0,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	remotes, err := repo.Backend.Remotes()
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	branchNameToShip := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToShip, hasBranchToShip := branchesSnapshot.Branches.FindByLocalName(branchNameToShip).Get()
	if hasBranchToShip && branchToShip.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.ShipBranchOtherWorktree, branchNameToShip)
	}
	isShippingInitialBranch := branchNameToShip == branchesSnapshot.Active
	if !hasBranchToShip {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToShip)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{branchToShip.LocalName},
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	if err = validateShippableBranchType(validatedConfig.Config.BranchType(branchNameToShip)); err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	targetBranchName, hasTargetBranch := validatedConfig.Config.Lineage.Parent(branchNameToShip).Get()
	if !hasTargetBranch {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.ShipBranchHasNoParent, branchNameToShip)
	}
	targetBranch, hasTargetBranch := branchesSnapshot.Branches.FindByLocalName(targetBranchName).Get()
	if !hasTargetBranch {
		return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	err = ensureParentBranchIsMainOrPerennialBranch(branchNameToShip, targetBranchName, validatedConfig.Config, validatedConfig.Config.Lineage)
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	var proposalOpt Option[hostingdomain.Proposal]
	childBranches := validatedConfig.Config.Lineage.Children(branchNameToShip)
	proposalsOfChildBranches := []hostingdomain.Proposal{}
	var connector hostingdomain.Connector
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connector, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          *validatedConfig.Config.UnvalidatedConfig,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSize, false, err
		}
	}
	canShipViaAPI := false
	proposalMessage := ""
	if !repo.IsOffline.Bool() && connector != nil {
		if branchToShip.HasTrackingBranch() {
			proposalOpt, err = connector.FindProposal(branchNameToShip, targetBranchName)
			if err != nil {
				return nil, branchesSnapshot, stashSize, false, err
			}
			proposal, hasProposal := proposalOpt.Get()
			if hasProposal {
				canShipViaAPI = true
				proposalMessage = connector.DefaultProposalMessage(proposal)
			}
		}
		for _, childBranch := range childBranches {
			childProposalOpt, err := connector.FindProposal(childBranch, branchNameToShip)
			if err != nil {
				return nil, branchesSnapshot, stashSize, false, fmt.Errorf(messages.ProposalNotFoundForBranch, branchNameToShip, err)
			}
			childProposal, hasChildProposal := childProposalOpt.Get()
			if hasChildProposal {
				proposalsOfChildBranches = append(proposalsOfChildBranches, childProposal)
			}
		}
	}
	return &shipData{
		allBranches:              branchesSnapshot.Branches,
		branchToShip:             branchToShip,
		canShipViaAPI:            canShipViaAPI,
		childBranches:            childBranches,
		config:                   validatedConfig,
		connector:                connector,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            branchesSnapshot.Active,
		isShippingInitialBranch:  isShippingInitialBranch,
		previousBranch:           previousBranch,
		proposal:                 proposalOpt,
		proposalMessage:          proposalMessage,
		proposalsOfChildBranches: proposalsOfChildBranches,
		remotes:                  remotes,
		targetBranch:             targetBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func ensureParentBranchIsMainOrPerennialBranch(branch, parentBranch gitdomain.LocalBranchName, config configdomain.ValidatedConfig, lineage configdomain.Lineage) error {
	if !config.IsMainOrPerennialBranch(parentBranch) {
		ancestors := lineage.Ancestors(branch)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		return fmt.Errorf(messages.ShipChildBranch, stringslice.Connect(ancestorsWithoutMainOrPerennial.Strings()), oldestAncestor)
	}
	return nil
}

func shipProgram(data *shipData, commitMessage Option[gitdomain.CommitMessage]) program.Program {
	prog := program.Program{}
	if data.config.Config.SyncBeforeShip {
		// sync the parent branch
		sync.BranchProgram(data.targetBranch, sync.BranchProgramArgs{
			BranchInfos:   data.allBranches,
			Config:        data.config.Config,
			InitialBranch: data.initialBranch,
			Remotes:       data.remotes,
			Program:       &prog,
			PushBranch:    true,
		})
		// sync the branch to ship (local sync only)
		sync.BranchProgram(data.branchToShip, sync.BranchProgramArgs{
			BranchInfos:   data.allBranches,
			Config:        data.config.Config,
			InitialBranch: data.initialBranch,
			Remotes:       data.remotes,
			Program:       &prog,
			PushBranch:    false,
		})
	}
	prog.Add(&opcodes.EnsureHasShippableChanges{Branch: data.branchToShip.LocalName, Parent: data.config.Config.MainBranch})
	prog.Add(&opcodes.Checkout{Branch: data.targetBranch.LocalName})
	if proposal, hasProposal := data.proposal.Get(); hasProposal && data.canShipViaAPI {
		// update the proposals of child branches
		for _, childProposal := range data.proposalsOfChildBranches {
			prog.Add(&opcodes.UpdateProposalTarget{
				ProposalNumber: childProposal.Number,
				NewTarget:      data.targetBranch.LocalName,
			})
		}
		prog.Add(&opcodes.PushCurrentBranch{CurrentBranch: data.branchToShip.LocalName})
		prog.Add(&opcodes.ConnectorMergeProposal{
			Branch:          data.branchToShip.LocalName,
			ProposalNumber:  proposal.Number,
			CommitMessage:   commitMessage,
			ProposalMessage: data.proposalMessage,
		})
		prog.Add(&opcodes.PullCurrentBranch{})
	} else {
		prog.Add(&opcodes.SquashMerge{Branch: data.branchToShip.LocalName, CommitMessage: commitMessage, Parent: data.targetBranch.LocalName})
	}
	if data.remotes.HasOrigin() && data.config.Config.IsOnline() {
		prog.Add(&opcodes.PushCurrentBranch{CurrentBranch: data.targetBranch.LocalName})
	}
	// NOTE: when shipping via API, we can always delete the tracking branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via API)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if data.canShipViaAPI || (data.branchToShip.HasTrackingBranch() && len(data.childBranches) == 0 && data.config.Config.IsOnline()) {
		if data.config.Config.ShipDeleteTrackingBranch {
			prog.Add(&opcodes.DeleteTrackingBranch{Branch: data.branchToShip.RemoteName})
		}
	}
	prog.Add(&opcodes.DeleteLocalBranch{Branch: data.branchToShip.LocalName})
	if !data.dryRun {
		prog.Add(&opcodes.DeleteParentBranch{Branch: data.branchToShip.LocalName})
	}
	for _, child := range data.childBranches {
		prog.Add(&opcodes.ChangeParent{Branch: child, Parent: data.targetBranch.LocalName})
	}
	if !data.isShippingInitialBranch {
		prog.Add(&opcodes.Checkout{Branch: data.initialBranch})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         !data.isShippingInitialBranch && data.hasOpenChanges,
		PreviousBranchCandidates: data.allBranches.WithNames(data.previousBranch),
	})
	return prog
}

func validateShippableBranchType(branchType configdomain.BranchType) error {
	switch branchType {
	case configdomain.BranchTypeContributionBranch:
		return errors.New(messages.ContributionBranchCannotShip)
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.MainBranchCannotShip)
	case configdomain.BranchTypeObservedBranch:
		return errors.New(messages.ObservedBranchCannotShip)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.PerennialBranchCannotShip)
	}
	panic(fmt.Sprintf("unhandled branch type: %v", branchType))
}

func validateData(data shipData) error {
	if data.branchToShip.LocalName == data.initialBranch {
		return validate.NoOpenChanges(data.hasOpenChanges)
	}
	return nil
}
