package ship

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks/slice"
	"github.com/git-town/git-town/v15/internal/hosting"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/validate"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

// data that all ship strategies use
type sharedShipData struct {
	branchNameToShip         gitdomain.LocalBranchName
	branchToShip             gitdomain.BranchInfo
	branchesSnapshot         gitdomain.BranchesSnapshot
	childBranches            gitdomain.LocalBranchNames
	config                   config.ValidatedConfig
	connector                Option[hostingdomain.Connector]
	dialogTestInputs         components.TestInputs
	dryRun                   configdomain.DryRun
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	isShippingInitialBranch  bool
	previousBranch           Option[gitdomain.LocalBranchName]
	proposalsOfChildBranches []hostingdomain.Proposal
	stashSize                gitdomain.StashSize
	targetBranch             gitdomain.BranchInfo
	targetBranchName         gitdomain.LocalBranchName
}

func determineSharedShipData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data sharedShipData, exit bool, err error) {
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
	if connector, hasConnector := connectorOpt.Get(); hasConnector {
		if !repo.IsOffline.IsTrue() {
			if branchToShip.HasTrackingBranch() {
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
	}
	return sharedShipData{
		branchNameToShip:         branchNameToShip,
		branchToShip:             *branchToShip,
		branchesSnapshot:         branchesSnapshot,
		childBranches:            childBranches,
		config:                   validatedConfig,
		connector:                connectorOpt,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		isShippingInitialBranch:  isShippingInitialBranch,
		previousBranch:           previousBranch,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
		targetBranch:             *targetBranch,
		targetBranchName:         targetBranchName,
	}, false, nil
}
