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
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
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

Ships only direct children of the main branch. To ship a child branch, ship or kill all ancestor branches first.

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
		Args:  cobra.NoArgs,
		Short: shipDesc,
		Long:  cmdhelpers.Long(shipDesc, fmt.Sprintf(shipHelp, gitconfig.KeyGithubToken, gitconfig.KeyShipDeleteTrackingBranch)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeShip(readMessageFlag(cmd), readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	return &cmd
}

func executeShip(message Option[gitdomain.CommitMessage], dryRun, verbose bool) error {
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
	data, exit, err := determineShipData(repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	err = validateData(*data)
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               shipCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
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
	allBranches              gitdomain.BranchInfos
	branchToShip             gitdomain.BranchInfo
	branchesSnapshot         gitdomain.BranchesSnapshot
	canShipViaAPI            bool
	childBranches            gitdomain.LocalBranchNames
	config                   config.ValidatedConfig
	connector                Option[hostingdomain.Connector]
	dialogTestInputs         components.TestInputs
	dryRun                   bool
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	previousBranch           Option[gitdomain.LocalBranchName]
	proposal                 Option[hostingdomain.Proposal]
	proposalMessage          string
	proposalsOfChildBranches []hostingdomain.Proposal
	remotes                  gitdomain.Remotes
	stashSize                gitdomain.StashSize
	targetBranch             gitdomain.BranchInfo
}

func determineShipData(repo execute.OpenRepoResult, dryRun, verbose bool) (*shipData, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return nil, false, err
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
		ValidateNoOpenChanges: true,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return nil, false, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return nil, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchToShip, _ := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if branchToShip.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return nil, false, fmt.Errorf(messages.ShipBranchOtherWorktree, initialBranch)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	if err = validateShippableBranchType(validatedConfig.Config.BranchType(initialBranch)); err != nil {
		return nil, false, err
	}
	targetBranchName, hasTargetBranch := validatedConfig.Config.Lineage.Parent(initialBranch).Get()
	if !hasTargetBranch {
		return nil, false, fmt.Errorf(messages.ShipBranchHasNoParent, initialBranch)
	}
	targetBranch, hasTargetBranch := branchesSnapshot.Branches.FindByLocalName(targetBranchName).Get()
	if !hasTargetBranch {
		return nil, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	err = ensureParentBranchIsMainOrPerennialBranch(initialBranch, targetBranchName, validatedConfig.Config, validatedConfig.Config.Lineage)
	if err != nil {
		return nil, false, err
	}
	var proposalOpt Option[hostingdomain.Proposal]
	childBranches := validatedConfig.Config.Lineage.Children(initialBranch)
	proposalsOfChildBranches := []hostingdomain.Proposal{}
	var connectorOpt Option[hostingdomain.Connector]
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connectorOpt, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          *validatedConfig.Config.UnvalidatedConfig,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return nil, false, err
		}
	}
	canShipViaAPI := false
	proposalMessage := ""
	if connector, hasConnector := connectorOpt.Get(); hasConnector {
		if !repo.IsOffline.Bool() {
			if branchToShip.HasTrackingBranch() {
				proposalOpt, err = connector.FindProposal(initialBranch, targetBranchName)
				if err != nil {
					return nil, false, err
				}
				proposal, hasProposal := proposalOpt.Get()
				if hasProposal {
					canShipViaAPI = true
					proposalMessage = connector.DefaultProposalMessage(proposal)
				}
			}
			for _, childBranch := range childBranches {
				childProposalOpt, err := connector.FindProposal(childBranch, initialBranch)
				if err != nil {
					return nil, false, fmt.Errorf(messages.ProposalNotFoundForBranch, initialBranch, err)
				}
				childProposal, hasChildProposal := childProposalOpt.Get()
				if hasChildProposal {
					proposalsOfChildBranches = append(proposalsOfChildBranches, childProposal)
				}
			}
		}
	}
	return &shipData{
		allBranches:              branchesSnapshot.Branches,
		branchToShip:             branchToShip,
		branchesSnapshot:         branchesSnapshot,
		canShipViaAPI:            canShipViaAPI,
		childBranches:            childBranches,
		config:                   validatedConfig,
		connector:                connectorOpt,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		previousBranch:           previousBranch,
		proposal:                 proposalOpt,
		proposalMessage:          proposalMessage,
		proposalsOfChildBranches: proposalsOfChildBranches,
		remotes:                  remotes,
		stashSize:                stashSize,
		targetBranch:             targetBranch,
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

func shipProgram(data *shipData, commitMessage Option[gitdomain.CommitMessage]) program.Program {
	prog := program.Program{}
	if data.config.Config.Online().Bool() {
		if trackingBranchName, hasTrackingBranch := data.branchToShip.RemoteName.Get(); hasTrackingBranch {
			if data.branchToShip.SyncStatus == gitdomain.SyncStatusNotInSync {
				if data.canShipViaAPI {
					// shipping a branch via API --> push missing local commits to the tracking branch
					// TODO
				} else {
					// shipping a local branch --> pull missing commits from the tracking branch
					switch data.config.Config.SyncFeatureStrategy {
					case configdomain.SyncFeatureStrategyMerge:
						prog.Add(&opcodes.Merge{Branch: trackingBranchName.BranchName()})
					case configdomain.SyncFeatureStrategyRebase:
						prog.Add(&opcodes.RebaseFeatureTrackingBranch{RemoteBranch: trackingBranchName})
					}
				}
			}
		}
	}
	localBranchToShip, hasLocalBranchToShip := data.branchToShip.LocalName.Get()
	localTargetBranch, _ := data.targetBranch.LocalName.Get()
	if hasLocalBranchToShip {
		prog.Add(&opcodes.EnsureHasShippableChanges{Branch: localBranchToShip, Parent: data.config.Config.MainBranch})
		prog.Add(&opcodes.Checkout{Branch: localTargetBranch})
	}
	if proposal, hasProposal := data.proposal.Get(); hasProposal && data.canShipViaAPI {
		// update the proposals of child branches
		for _, childProposal := range data.proposalsOfChildBranches {
			prog.Add(&opcodes.UpdateProposalTarget{
				ProposalNumber: childProposal.Number,
				NewTarget:      localTargetBranch,
			})
		}
		prog.Add(&opcodes.PushCurrentBranch{CurrentBranch: localBranchToShip})
		prog.Add(&opcodes.ConnectorMergeProposal{
			Branch:          localBranchToShip,
			ProposalNumber:  proposal.Number,
			CommitMessage:   commitMessage,
			ProposalMessage: data.proposalMessage,
		})
		prog.Add(&opcodes.PullCurrentBranch{})
	} else {
		prog.Add(&opcodes.SquashMerge{Branch: localBranchToShip, CommitMessage: commitMessage, Parent: localTargetBranch})
	}
	if data.remotes.HasOrigin() && data.config.Config.IsOnline() {
		prog.Add(&opcodes.PushCurrentBranch{CurrentBranch: localTargetBranch})
	}
	// NOTE: when shipping via API, we can always delete the tracking branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via API)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if branchToShipRemoteName, hasRemoteName := data.branchToShip.RemoteName.Get(); hasRemoteName {
		if data.canShipViaAPI || (data.branchToShip.HasTrackingBranch() && len(data.childBranches) == 0 && data.config.Config.IsOnline()) {
			if data.config.Config.ShipDeleteTrackingBranch {
				prog.Add(&opcodes.DeleteTrackingBranch{Branch: branchToShipRemoteName})
			}
		}
	}
	prog.Add(&opcodes.DeleteLocalBranch{Branch: localBranchToShip})
	if !data.dryRun {
		prog.Add(&opcodes.DeleteParentBranch{Branch: localBranchToShip})
	}
	for _, child := range data.childBranches {
		prog.Add(&opcodes.ChangeParent{Branch: child, Parent: localTargetBranch})
	}
	previousBranchCandidates := gitdomain.LocalBranchNames{}
	if previousBranch, hasPreviousBranch := data.previousBranch.Get(); hasPreviousBranch {
		previousBranchCandidates = append(previousBranchCandidates, previousBranch)
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         false,
		PreviousBranchCandidates: previousBranchCandidates,
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
	if localName, hasLocalName := data.branchToShip.LocalName.Get(); hasLocalName {
		if localName == data.initialBranch {
			return validate.NoOpenChanges(data.hasOpenChanges)
		}
	}
	return nil
}
