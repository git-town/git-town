package ship

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/slice"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/validate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
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
	switch validatedConfig.Config.BranchType(branchNameToShip) {
	case configdomain.BranchTypeContributionBranch:
		return data, false, errors.New(messages.ContributionBranchCannotShip)
	case configdomain.BranchTypeMainBranch:
		return data, false, errors.New(messages.MainBranchCannotShip)
	case configdomain.BranchTypeObservedBranch:
		return data, false, errors.New(messages.ObservedBranchCannotShip)
	case configdomain.BranchTypePerennialBranch:
		return data, false, errors.New(messages.PerennialBranchCannotShip)
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
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
	connectorOpt, err := hosting.NewConnector(hosting.NewConnectorArgs{
		Config:          *validatedConfig.Config.UnvalidatedConfig,
		HostingPlatform: validatedConfig.Config.HostingPlatform,
		Log:             print.Logger{},
		RemoteURL:       validatedConfig.OriginURL(),
	})
	if err != nil {
		return data, false, err
	}
	proposalsOfChildBranches := LoadProposalsOfChildBranches(LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connectorOpt,
		Lineage:                    validatedConfig.Config.Lineage,
		Offline:                    repo.IsOffline,
		OldBranch:                  branchNameToShip,
		OldBranchHasTrackingBranch: branchToShip.HasTrackingBranch(),
	})
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

func LoadProposalsOfChildBranches(args LoadProposalsOfChildBranchesArgs) []hostingdomain.Proposal {
	var proposalsOfChildBranches []hostingdomain.Proposal
	childBranches := args.Lineage.Children(args.OldBranch)
	connector, hasConnector := args.ConnectorOpt.Get()
	if hasConnector && connector.CanMakeAPICalls() {
		if !args.Offline.IsTrue() && args.OldBranchHasTrackingBranch {
			for _, childBranch := range childBranches {
				childProposalOpt, err := connector.FindProposal(childBranch, args.OldBranch)
				if err != nil {
					continue
				}
				childProposal, hasChildProposal := childProposalOpt.Get()
				if hasChildProposal {
					proposalsOfChildBranches = append(proposalsOfChildBranches, childProposal)
				}
			}
		}
	}
	return proposalsOfChildBranches
}

type LoadProposalsOfChildBranchesArgs struct {
	ConnectorOpt               Option[hostingdomain.Connector]
	Lineage                    configdomain.Lineage
	Offline                    configdomain.Offline
	OldBranch                  gitdomain.LocalBranchName
	OldBranchHasTrackingBranch bool
}
