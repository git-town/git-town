package ship

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// data only needed for shipping via the API
type shipDataAPI struct {
	branchToShipRemoteName gitdomain.RemoteBranchName
	connector              forgedomain.Connector
	proposal               forgedomain.Proposal
}

func determineAPIData(sharedData sharedShipData) (result shipDataAPI, err error) {
	branchToShipRemoteName, hasRemoteBranchToShip := sharedData.branchToShipInfo.RemoteName.Get()
	if !hasRemoteBranchToShip {
		return result, fmt.Errorf(messages.ShipAPINoRemoteBranch, sharedData.branchToShip)
	}
	connector, hasConnector := sharedData.connector.Get()
	if !hasConnector {
		return result, errors.New(messages.ShipAPIConnectorRequired)
	}
	proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder)
	if !canFindProposals {
		return result, errors.New(messages.ShipAPIConnectorUnsupported)
	}
	proposalOpt, err := proposalFinder.FindProposal(sharedData.branchToShip, sharedData.targetBranchName)
	proposal, hasProposal := proposalOpt.Get()
	if !hasProposal {
		return result, fmt.Errorf(messages.ShipAPINoProposal, sharedData.branchToShip)
	}
	return shipDataAPI{
		branchToShipRemoteName: branchToShipRemoteName,
		connector:              connector,
		proposal:               proposal,
	}, err
}

func shipAPIProgram(prog Mutable[program.Program], repo execute.OpenRepoResult, sharedData sharedShipData, apiData shipDataAPI, commitMessage Option[gitdomain.CommitMessage]) error {
	branchToShipLocal, hasLocalBranchToShip := sharedData.branchToShipInfo.LocalName.Get()
	UpdateChildBranchProposalsToGrandParent(prog.Value, sharedData.proposalsOfChildBranches)
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.targetBranchName})
	connector, hasConnector := sharedData.connector.Get()
	if !hasConnector {
		return errors.New(messages.ShipAPIConnectorRequired)
	}
	_, canMergeProposals := connector.(forgedomain.ProposalMerger)
	if !canMergeProposals {
		return errors.New(messages.ShipAPIConnectorUnsupported)
	}
	prog.Value.Add(&opcodes.ConnectorProposalMerge{
		Branch:        branchToShipLocal,
		Proposal:      apiData.proposal,
		CommitMessage: commitMessage,
	})
	if sharedData.config.NormalConfig.ShipDeleteTrackingBranch {
		prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: apiData.branchToShipRemoteName})
	}
	if hasLocalBranchToShip {
		prog.Value.Add(&opcodes.BranchLocalDelete{Branch: branchToShipLocal})
	}
	for _, child := range sharedData.childBranches {
		prog.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	if !repo.UnvalidatedConfig.NormalConfig.DryRun {
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: branchToShipLocal})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{sharedData.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   repo.UnvalidatedConfig.NormalConfig.DryRun,
		InitialStashSize:         sharedData.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         !sharedData.isShippingInitialBranch && sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return nil
}
