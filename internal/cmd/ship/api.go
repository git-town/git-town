package ship

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v20/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
	"github.com/git-town/git-town/v20/internal/vm/program"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// data only needed for shipping via the API
type shipDataAPI struct {
	branchToShipRemoteName gitdomain.RemoteBranchName
	connector              forgedomain.Connector
	proposal               forgedomain.Proposal
}

func determineAPIData(sharedData sharedShipData) (result shipDataAPI, err error) {
	branchToShipRemoteName, hasRemoteBranchToShip := sharedData.branchToShip.RemoteName.Get()
	if !hasRemoteBranchToShip {
		return result, fmt.Errorf(messages.ShipAPINoRemoteBranch, sharedData.branchNameToShip)
	}
	connector, hasConnector := sharedData.connector.Get()
	if !hasConnector {
		return result, errors.New(messages.ShipAPIConnectorRequired)
	}
	findProposal, canFindProposal := connector.FindProposalFn().Get()
	if !canFindProposal {
		return result, errors.New(messages.ShipAPIConnectorUnsupported)
	}
	proposalOpt, err := findProposal(sharedData.branchNameToShip, sharedData.targetBranchName)
	if err != nil {
		return result, err
	}
	proposal, hasProposal := proposalOpt.Get()
	if !hasProposal {
		return result, fmt.Errorf(messages.ShipAPINoProposal, sharedData.branchNameToShip)
	}
	return shipDataAPI{
		branchToShipRemoteName: branchToShipRemoteName,
		connector:              connector,
		proposal:               proposal,
	}, nil
}

func shipAPIProgram(prog Mutable[program.Program], sharedData sharedShipData, apiData shipDataAPI, commitMessage Option[gitdomain.CommitMessage]) error {
	branchToShipLocal, hasLocalBranchToShip := sharedData.branchToShip.LocalName.Get()
	UpdateChildBranchProposalsToGrandParent(prog.Value, sharedData.proposalsOfChildBranches)
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: sharedData.targetBranchName})
	connector, hasConnector := sharedData.connector.Get()
	if !hasConnector {
		return errors.New(messages.ShipAPIConnectorRequired)
	}
	if connector.SquashMergeProposalFn().IsNone() {
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
	if !sharedData.dryRun {
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: branchToShipLocal})
	}
	for _, child := range sharedData.childBranches {
		prog.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{sharedData.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   sharedData.dryRun,
		InitialStashSize:         sharedData.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         !sharedData.isShippingInitialBranch && sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return nil
}
