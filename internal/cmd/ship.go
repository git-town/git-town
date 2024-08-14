package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/gohacks/slice"
	"github.com/git-town/git-town/v15/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v15/internal/hosting"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/undo/undoconfig"
	"github.com/git-town/git-town/v15/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v15/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	"github.com/git-town/git-town/v15/internal/vm/runstate"
	"github.com/spf13/cobra"
)

const shipCommand = "ship"

const shipDesc = "Deliver a completed feature branch"

const shipHelp = `
Squash-merges the current branch, or <branch_name> if given, into the main branch, resulting in linear history on the main branch.

Ships only direct children of the main branch. To ship a child branch, ship or kill all ancestor branches first.

If you use GitHub, this command can squash merge pull requests via the GitHub API. Setup:

1. Get a GitHub personal access token with the "repo" scope
2. Run 'git config %s <token>' (optionally add the '--global' flag)

Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API. It will also update the base branch for any pull requests against that branch.

If your origin server deletes shipped branches, for example GitHub's feature to automatically delete head branches, run "git config %s false" and Git Town will leave it up to your origin server to delete the tracking branch of the branch you are shipping.`

func shipCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message for the squash commit")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addToParentFlag, readToParentFlag := flags.ShipIntoNonPerennialParent()
	cmd := cobra.Command{
		Use:   shipCommand,
		Args:  cobra.MaximumNArgs(1),
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, configdomain.KeyGithubToken, configdomain.KeyShipDeleteTrackingBranch)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeShip(args, readMessageFlag(cmd), readDryRunFlag(cmd), readVerboseFlag(cmd), readToParentFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	addToParentFlag(&cmd)
	return &cmd
}

func executeShip(args []string, message Option[gitdomain.CommitMessage], dryRun configdomain.DryRun, verbose configdomain.Verbose, toParent configdomain.ShipIntoNonperennialParent) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineShipData(args, repo, dryRun, verbose, toParent)
	if err != nil || exit {
		return err
	}
	err = validateData(data)
	if err != nil {
		return err
	}
	runProgram := shipProgram(data, message)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               shipCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type shipData struct {
	allBranches                gitdomain.BranchInfos
	branchToShip               gitdomain.BranchInfo
	branchesSnapshot           gitdomain.BranchesSnapshot
	canShipViaAPI              bool
	childBranches              gitdomain.LocalBranchNames
	config                     config.ValidatedConfig
	connector                  Option[hostingdomain.Connector]
	dialogTestInputs           components.TestInputs
	dryRun                     configdomain.DryRun
	hasOpenChanges             bool
	initialBranch              gitdomain.LocalBranchName
	isShippingInitialBranch    bool
	previousBranch             Option[gitdomain.LocalBranchName]
	proposal                   Option[hostingdomain.Proposal]
	proposalMessage            string
	proposalsOfChildBranches   []hostingdomain.Proposal
	remotes                    gitdomain.Remotes
	shipIntoNonPerennialParent configdomain.ShipIntoNonperennialParent
	stashSize                  gitdomain.StashSize
	targetBranch               gitdomain.BranchInfo
}

func determineShipData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose, shipIntoNonPerennialParent configdomain.ShipIntoNonperennialParent) (data shipData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: len(args) == 0,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchNameToShip := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToShip, hasBranchToShip := branchesSnapshot.Branches.FindByLocalName(branchNameToShip).Get()
	if hasBranchToShip && branchToShip.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, false, fmt.Errorf(messages.ShipBranchOtherWorktree, branchNameToShip)
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	isShippingInitialBranch := branchNameToShip == initialBranch
	if !hasBranchToShip {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToShip)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{branchNameToShip},
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return data, exit, err
	}
	if err = validateShippableBranchType(validatedConfig.Config.BranchType(branchNameToShip)); err != nil {
		return data, false, err
	}
	targetBranchName, hasTargetBranch := validatedConfig.Config.Lineage.Parent(branchNameToShip).Get()
	if !hasTargetBranch {
		return data, false, fmt.Errorf(messages.ShipBranchHasNoParent, branchNameToShip)
	}
	targetBranch, hasTargetBranch := branchesSnapshot.Branches.FindByLocalName(targetBranchName).Get()
	if !hasTargetBranch {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	var proposalOpt Option[hostingdomain.Proposal]
	childBranches := validatedConfig.Config.Lineage.Children(branchNameToShip)
	proposalsOfChildBranches := []hostingdomain.Proposal{}
	var connectorOpt Option[hostingdomain.Connector]
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connectorOpt, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          *validatedConfig.Config.UnvalidatedConfig,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			RemoteURL:       originURL,
		})
		if err != nil {
			return data, false, err
		}
	}
	canShipViaAPI := false
	proposalMessage := ""
	if connector, hasConnector := connectorOpt.Get(); hasConnector {
		if !repo.IsOffline.IsTrue() {
			if branchToShip.HasTrackingBranch() {
				proposalOpt, err = connector.FindProposal(branchNameToShip, targetBranchName)
				if err != nil {
					return data, false, err
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
					return data, false, fmt.Errorf(messages.ProposalNotFoundForBranch, branchNameToShip, err)
				}
				childProposal, hasChildProposal := childProposalOpt.Get()
				if hasChildProposal {
					proposalsOfChildBranches = append(proposalsOfChildBranches, childProposal)
				}
			}
		}
	}
	return shipData{
		allBranches:                branchesSnapshot.Branches,
		branchToShip:               *branchToShip,
		branchesSnapshot:           branchesSnapshot,
		canShipViaAPI:              canShipViaAPI,
		childBranches:              childBranches,
		config:                     validatedConfig,
		connector:                  connectorOpt,
		dialogTestInputs:           dialogTestInputs,
		dryRun:                     dryRun,
		hasOpenChanges:             repoStatus.OpenChanges,
		initialBranch:              initialBranch,
		isShippingInitialBranch:    isShippingInitialBranch,
		previousBranch:             previousBranch,
		proposal:                   proposalOpt,
		proposalMessage:            proposalMessage,
		proposalsOfChildBranches:   proposalsOfChildBranches,
		remotes:                    remotes,
		shipIntoNonPerennialParent: shipIntoNonPerennialParent,
		stashSize:                  stashSize,
		targetBranch:               *targetBranch,
	}, false, nil
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

func shipProgram(data shipData, commitMessage Option[gitdomain.CommitMessage]) program.Program {
	prog := NewMutable(&program.Program{})
	localBranchToShip, hasLocalBranchToShip := data.branchToShip.LocalName.Get()
	localTargetBranch, _ := data.targetBranch.LocalName.Get()
	if hasLocalBranchToShip {
		prog.Value.Add(&opcodes.EnsureHasShippableChanges{Branch: localBranchToShip, Parent: data.config.Config.MainBranch})
		prog.Value.Add(&opcodes.Checkout{Branch: localTargetBranch})
	}
	if proposal, hasProposal := data.proposal.Get(); hasProposal && data.canShipViaAPI {
		// update the proposals of child branches
		for _, childProposal := range data.proposalsOfChildBranches {
			prog.Value.Add(&opcodes.UpdateProposalTarget{
				ProposalNumber: childProposal.Number,
				NewTarget:      localTargetBranch,
			})
		}
		prog.Value.Add(&opcodes.PushCurrentBranch{CurrentBranch: localBranchToShip})
		prog.Value.Add(&opcodes.ConnectorMergeProposal{
			Branch:          localBranchToShip,
			ProposalNumber:  proposal.Number,
			CommitMessage:   commitMessage,
			ProposalMessage: data.proposalMessage,
		})
		prog.Value.Add(&opcodes.PullCurrentBranch{})
	} else {
		prog.Value.Add(&opcodes.SquashMerge{Branch: localBranchToShip, CommitMessage: commitMessage, Parent: localTargetBranch})
	}
	if data.remotes.HasOrigin() && data.config.Config.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranch{CurrentBranch: localTargetBranch})
	}
	// NOTE: when shipping via API, we can always delete the tracking branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via API)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if branchToShipRemoteName, hasRemoteName := data.branchToShip.RemoteName.Get(); hasRemoteName {
		if data.canShipViaAPI || (data.branchToShip.HasTrackingBranch() && len(data.childBranches) == 0 && data.config.Config.IsOnline()) {
			if data.config.Config.ShipDeleteTrackingBranch {
				prog.Value.Add(&opcodes.DeleteTrackingBranch{Branch: branchToShipRemoteName})
			}
		}
	}
	prog.Value.Add(&opcodes.DeleteLocalBranch{Branch: localBranchToShip})
	if !data.dryRun {
		prog.Value.Add(&opcodes.DeleteParentBranch{Branch: localBranchToShip})
	}
	for _, child := range data.childBranches {
		prog.Value.Add(&opcodes.ChangeParent{Branch: child, Parent: localTargetBranch})
	}
	if !data.isShippingInitialBranch {
		prog.Value.Add(&opcodes.Checkout{Branch: data.initialBranch})
	}
	previousBranchCandidates := gitdomain.LocalBranchNames{}
	if previousBranch, hasPreviousBranch := data.previousBranch.Get(); hasPreviousBranch {
		previousBranchCandidates = append(previousBranchCandidates, previousBranch)
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         !data.isShippingInitialBranch && data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Get()
}

func validateShippableBranchType(branchType configdomain.BranchType) error {
	switch branchType {
	case configdomain.BranchTypeContributionBranch:
		return errors.New(messages.ContributionBranchCannotShip)
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
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
	if !data.shipIntoNonPerennialParent {
		err := ensureParentBranchIsMainOrPerennialBranch(data.branchToShip.LocalName.GetOrPanic(), data.targetBranch.LocalName.GetOrPanic(), data.config.Config, data.config.Config.Lineage)
		if err != nil {
			return err
		}
	}
	switch data.branchToShip.SyncStatus {
	case gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.ShipBranchDeletedAtRemote, data.branchToShip.LocalName)
	case gitdomain.SyncStatusNotInSync:
		return fmt.Errorf(messages.ShipBranchNotInSync, data.branchToShip.LocalName)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.ShipBranchIsInOtherWorktree, data.branchToShip.LocalName)
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusRemoteOnly, gitdomain.SyncStatusLocalOnly:
	}
	if localName, hasLocalName := data.branchToShip.LocalName.Get(); hasLocalName {
		if localName == data.initialBranch {
			return validate.NoOpenChanges(data.hasOpenChanges)
		}
	}
	return nil
}
