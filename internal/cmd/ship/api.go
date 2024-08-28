package ship

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

// data only needed for shipping via the API
type shipDataAPI struct {
	branchToShipRemoteName gitdomain.RemoteBranchName
	connector              hostingdomain.Connector
	proposal               hostingdomain.Proposal
	proposalMessage        string
}

func determineAPIData(sharedData sharedShipData) (result shipDataAPI, err error) {
	proposalMessage := ""
	connector, hasConnector := sharedData.connector.Get()
	if !hasConnector {
		return result, fmt.Errorf(messages.ShipAPIConnectorRequired)
	}
	proposalOpt, err := connector.FindProposal(sharedData.branchNameToShip, sharedData.targetBranchName)
	if err != nil {
		return result, err
	}
	proposal, hasProposal := proposalOpt.Get()
	if !hasProposal {
		return result, fmt.Errorf(messages.ShipAPINoProposal, sharedData.branchNameToShip)
	}
	proposalMessage = connector.DefaultProposalMessage(proposal)
	branchToShipRemoteName, hasRemoteBranchToShip := sharedData.branchToShip.RemoteName.Get()
	if !hasRemoteBranchToShip {
		return result, fmt.Errorf(messages.ShipAPINoRemoteBranch, sharedData.branchNameToShip)
	}
	return shipDataAPI{
		branchToShipRemoteName: branchToShipRemoteName,
		connector:              connector,
		proposal:               proposal,
		proposalMessage:        proposalMessage,
	}, nil
}

func shipAPIProgram(sharedData sharedShipData, apiData shipDataAPI, commitMessage Option[gitdomain.CommitMessage]) program.Program {
	prog := NewMutable(&program.Program{})
	branchToShipLocal, hasLocalBranchToShip := sharedData.branchToShip.LocalName.Get()
	localTargetBranch, _ := sharedData.targetBranch.LocalName.Get()
	// update the proposals of child branches
	for _, childProposal := range sharedData.proposalsOfChildBranches {
		prog.Value.Add(&opcodes.UpdateProposalTarget{
			ProposalNumber: childProposal.Number,
			NewTarget:      localTargetBranch,
		})
	}
	prog.Value.Add(&opcodes.ConnectorMergeProposal{
		Branch:          branchToShipLocal,
		ProposalNumber:  apiData.proposal.Number,
		CommitMessage:   commitMessage,
		ProposalMessage: apiData.proposalMessage,
	})
	if sharedData.config.Config.ShipDeleteTrackingBranch {
		prog.Value.Add(&opcodes.DeleteTrackingBranch{Branch: apiData.branchToShipRemoteName})
	}
	if hasLocalBranchToShip {
		prog.Value.Add(&opcodes.DeleteLocalBranch{Branch: branchToShipLocal})
	}
	if !sharedData.dryRun {
		prog.Value.Add(&opcodes.DeleteParentBranch{Branch: branchToShipLocal})
	}
	for _, child := range sharedData.childBranches {
		prog.Value.Add(&opcodes.ChangeParent{Branch: child, Parent: localTargetBranch})
	}
	previousBranchCandidates := gitdomain.LocalBranchNames{}
	if previousBranch, hasPreviousBranch := sharedData.previousBranch.Get(); hasPreviousBranch {
		previousBranchCandidates = append(previousBranchCandidates, previousBranch)
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   sharedData.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         !sharedData.isShippingInitialBranch && sharedData.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Get()
}
