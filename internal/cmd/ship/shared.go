package ship

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/validate"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// data that all ship strategies use
type sharedShipData struct {
	branchNameToShip         gitdomain.LocalBranchName
	branchToShip             gitdomain.BranchInfo
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

func determineSharedShipData(args []string, repo execute.OpenRepoResult, shipStrategyOverride Option[configdomain.ShipStrategy]) (data sharedShipData, exit dialogdomain.Exit, err error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
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
		return data, false, err
	}
	branchesSnapshot, stashSize, previousBranchInfos, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
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
	if err != nil || exit {
		return data, exit, err
	}
	if branchesSnapshot.DetachedHead {
		return data, false, errors.New(messages.ShipRepoHasDetachedHead)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	var branchNameToShip gitdomain.LocalBranchName
	if len(args) > 0 {
		branchNameToShip = gitdomain.NewLocalBranchName(args[0])
	} else if activeBranch, hasActiveBranch := branchesSnapshot.Active.Get(); hasActiveBranch {
		branchNameToShip = activeBranch
	} else {
		return data, false, errors.New(messages.ShipNoBranchToShip)
	}
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
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{branchNameToShip},
		ConfigSnapshot:     repo.ConfigSnapshot,
		Connector:          data.connector,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		Inputs:             inputs,
		LocalBranches:      localBranches,
		Remotes:            remotes,
		RepoStatus:         repoStatus,
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
	proposalsOfChildBranches := LoadProposalsOfChildBranches(LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connector,
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
	}, false, nil
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
	childBranches := args.Lineage.Children(args.OldBranch)
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
