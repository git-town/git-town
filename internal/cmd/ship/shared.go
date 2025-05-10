package ship

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/print"
	"github.com/git-town/git-town/v20/internal/config"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/execute"
	"github.com/git-town/git-town/v20/internal/forge"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/gohacks/slice"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/validate"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// data that all ship strategies use
type sharedShipData struct {
	branchNameToShip         gitdomain.LocalBranchName
	branchToShip             gitdomain.BranchInfo
	branchesSnapshot         gitdomain.BranchesSnapshot
	childBranches            gitdomain.LocalBranchNames
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	dialogTestInputs         components.TestInputs
	dryRun                   configdomain.DryRun
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	isShippingInitialBranch  bool
	previousBranch           Option[gitdomain.LocalBranchName]
	previousBranchInfos      Option[gitdomain.BranchInfos]
	proposalsOfChildBranches []forgedomain.Proposal
	stashSize                gitdomain.StashSize
	targetBranch             gitdomain.BranchInfo
	targetBranchName         gitdomain.LocalBranchName
}

func determineSharedShipData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, shipStrategyOverride Option[configdomain.ShipStrategy], verbose configdomain.Verbose) (data sharedShipData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, previousBranchInfos, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
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
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{branchNameToShip},
		Connector:          data.connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	if shipStrategyOverride, hasShipStrategyOverride := shipStrategyOverride.Get(); hasShipStrategyOverride {
		validatedConfig.NormalConfig.ShipStrategy = shipStrategyOverride
	}
	switch validatedConfig.BranchType(branchNameToShip) {
	case configdomain.BranchTypeContributionBranch:
		return data, false, errors.New(messages.ContributionBranchCannotShip)
	case configdomain.BranchTypeMainBranch:
		return data, false, errors.New(messages.MainBranchCannotShip)
	case configdomain.BranchTypeObservedBranch:
		return data, false, errors.New(messages.ObservedBranchCannotShip)
	case configdomain.BranchTypePerennialBranch:
		return data, false, errors.New(messages.PerennialBranchCannotShip)
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	}
	targetBranchName, hasTargetBranch := validatedConfig.NormalConfig.Lineage.Parent(branchNameToShip).Get()
	if !hasTargetBranch {
		return data, false, fmt.Errorf(messages.ShipBranchHasNoParent, branchNameToShip)
	}
	targetBranch, hasTargetBranch := branchesSnapshot.Branches.FindByLocalName(targetBranchName).Get()
	if !hasTargetBranch {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	childBranches := validatedConfig.NormalConfig.Lineage.Children(branchNameToShip)
	connectorOpt, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	proposalsOfChildBranches := LoadProposalsOfChildBranches(LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connectorOpt,
		Lineage:                    validatedConfig.NormalConfig.Lineage,
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
		previousBranchInfos:      previousBranchInfos,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
		targetBranch:             *targetBranch,
		targetBranchName:         targetBranchName,
	}, false, nil
}

func LoadProposalsOfChildBranches(args LoadProposalsOfChildBranchesArgs) []forgedomain.Proposal {
	connector, hasConnector := args.ConnectorOpt.Get()
	if !hasConnector {
		return []forgedomain.Proposal{}
	}
	findProposal, canFindProposal := connector.FindProposalFn().Get()
	if !canFindProposal {
		return []forgedomain.Proposal{}
	}
	if args.Offline.IsOffline() {
		return []forgedomain.Proposal{}
	}
	if !args.OldBranchHasTrackingBranch {
		return []forgedomain.Proposal{}
	}
	childBranches := args.Lineage.Children(args.OldBranch)
	result := make([]forgedomain.Proposal, 0, len(childBranches))
	for _, childBranch := range childBranches {
		childProposalOpt, err := findProposal(childBranch, args.OldBranch)
		if err != nil {
			print.Error(err)
			continue
		}
		childProposal, hasChildProposal := childProposalOpt.Get()
		if !hasChildProposal {
			continue
		}
		result = append(result, childProposal)
	}
	return result
}

type LoadProposalsOfChildBranchesArgs struct {
	ConnectorOpt               Option[forgedomain.Connector]
	Lineage                    configdomain.Lineage
	Offline                    configdomain.Offline
	OldBranch                  gitdomain.LocalBranchName
	OldBranchHasTrackingBranch bool
}

func FindProposal(connectorOpt Option[forgedomain.Connector], sourceBranch gitdomain.LocalBranchName, targetBranch Option[gitdomain.LocalBranchName]) Option[forgedomain.Proposal] {
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return None[forgedomain.Proposal]()
	}
	target, hasTarget := targetBranch.Get()
	if !hasTarget {
		return None[forgedomain.Proposal]()
	}
	findProposal, canFindProposal := connector.FindProposalFn().Get()
	if !canFindProposal {
		return None[forgedomain.Proposal]()
	}
	proposal, err := findProposal(sourceBranch, target)
	if err != nil {
		print.Error(err)
		return None[forgedomain.Proposal]()
	}
	return proposal
}
