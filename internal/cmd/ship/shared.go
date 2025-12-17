package ship

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/validate"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// data that all ship strategies use
type sharedShipData struct {
	branchToShip             gitdomain.LocalBranchName
	branchToShipInfo         gitdomain.BranchInfo
	branchesSnapshot         gitdomain.BranchesSnapshot
	childBranches            gitdomain.LocalBranchNames
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	inputs                   dialogcomponents.Inputs
	isShippingInitialBranch  bool
	previousBranch           Option[gitdomain.LocalBranchName]
	previousBranchInfos      Option[gitdomain.BranchInfos]
	proposalsOfChildBranches []forgedomain.Proposal
	stashSize                gitdomain.StashSize
	targetBranch             gitdomain.BranchInfo
	targetBranchName         gitdomain.LocalBranchName
}

func determineSharedShipData(args []string, repo execute.OpenRepoResult, shipStrategyOverride Option[configdomain.ShipStrategy]) (sharedShipData, configdomain.ProgramFlow, error) {
	var emptyResult sharedShipData
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		Browser:              config.Browser,
		ForgeType:            config.ForgeType,
		ForgejoToken:         config.ForgejoToken,
		Frontend:             repo.Frontend,
		GitHubConnectorType:  config.GitHubConnectorType,
		GitHubToken:          config.GitHubToken,
		GitLabConnectorType:  config.GitLabConnectorType,
		GitLabToken:          config.GitLabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, stashSize, previousBranchInfos, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Inputs:                inputs,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: len(args) == 0,
	})
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return emptyResult, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.ShipRepoHasDetachedHead)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	var branchToShip gitdomain.LocalBranchName
	if len(args) > 0 {
		branchToShip = gitdomain.NewLocalBranchName(args[0])
	} else if activeBranch, hasActiveBranch := branchesSnapshot.Active.Get(); hasActiveBranch {
		branchToShip = activeBranch
	} else {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.ShipNoBranchToShip)
	}
	branchToShipInfo, hasBranchToShipInfo := branchesSnapshot.Branches.FindByLocalName(branchToShip).Get()
	if !hasBranchToShipInfo {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchDoesntExist, branchToShip)
	}
	if branchToShipInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.ShipBranchOtherWorktree, branchToShip)
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	isShippingInitialBranch := branchToShip == initialBranch
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{branchToShip},
		ConfigSnapshot:     repo.ConfigSnapshot,
		Connector:          connector,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		Inputs:             inputs,
		LocalBranches:      localBranches,
		Remotes:            remotes,
		RepoStatus:         repoStatus,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	if shipStrategyOverride, hasShipStrategyOverride := shipStrategyOverride.Get(); hasShipStrategyOverride {
		validatedConfig.NormalConfig.ShipStrategy = shipStrategyOverride
	}
	switch validatedConfig.BranchType(branchToShip) {
	case configdomain.BranchTypeContributionBranch:
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.ContributionBranchCannotShip)
	case configdomain.BranchTypeMainBranch:
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.MainBranchCannotShip)
	case configdomain.BranchTypeObservedBranch:
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.ObservedBranchCannotShip)
	case configdomain.BranchTypePerennialBranch:
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.PerennialBranchCannotShip)
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	}
	targetBranchName, hasTargetBranch := validatedConfig.NormalConfig.Lineage.Parent(branchToShip).Get()
	if !hasTargetBranch {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.ShipBranchHasNoParent, branchToShip)
	}
	targetBranch, hasTargetBranch := branchesSnapshot.Branches.FindByLocalName(targetBranchName).Get()
	if !hasTargetBranch {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	childBranches := validatedConfig.NormalConfig.Lineage.Children(branchToShip, validatedConfig.NormalConfig.Order)
	proposalsOfChildBranches := LoadProposalsOfChildBranches(LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connector,
		Lineage:                    validatedConfig.NormalConfig.Lineage,
		Offline:                    repo.IsOffline,
		OldBranch:                  branchToShip,
		OldBranchHasTrackingBranch: branchToShipInfo.HasTrackingBranch(),
		Order:                      validatedConfig.NormalConfig.Order,
	})
	return sharedShipData{
		branchToShip:             branchToShip,
		branchToShipInfo:         *branchToShipInfo,
		branchesSnapshot:         branchesSnapshot,
		childBranches:            childBranches,
		config:                   validatedConfig,
		connector:                connector,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		inputs:                   inputs,
		isShippingInitialBranch:  isShippingInitialBranch,
		previousBranch:           previousBranch,
		previousBranchInfos:      previousBranchInfos,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
		targetBranch:             *targetBranch,
		targetBranchName:         targetBranchName,
	}, configdomain.ProgramFlowContinue, nil
}

func LoadProposalsOfChildBranches(args LoadProposalsOfChildBranchesArgs) []forgedomain.Proposal {
	connector, hasConnector := args.ConnectorOpt.Get()
	if !hasConnector {
		return []forgedomain.Proposal{}
	}
	proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder)
	if !canFindProposals {
		return []forgedomain.Proposal{}
	}
	if args.Offline.IsOffline() {
		return []forgedomain.Proposal{}
	}
	if !args.OldBranchHasTrackingBranch {
		return []forgedomain.Proposal{}
	}
	childBranches := args.Lineage.Children(args.OldBranch, args.Order)
	result := make([]forgedomain.Proposal, 0, len(childBranches))
	for _, childBranch := range childBranches {
		childProposalOpt, err := proposalFinder.FindProposal(childBranch, args.OldBranch)
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
	Order                      configdomain.Order
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
	proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder)
	if !canFindProposals {
		return None[forgedomain.Proposal]()
	}
	proposal, err := proposalFinder.FindProposal(sourceBranch, target)
	if err != nil {
		print.Error(err)
		return None[forgedomain.Proposal]()
	}
	return proposal
}
