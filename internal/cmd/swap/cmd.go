package swap

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/cmd/sync"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const swapDesc = "Swap the position of this branch with its parent"

const swapHelp = `
The "swap" command moves the current branch
one position forward in the stack,
i.e. switches the position of the current branch
with its parent.

Consider this stack:

main
 \
  branch-1
   \
*   branch-2
     \
      branch-3

We are on the "branch-2" branch.
After running "git town swap",
you end up with this stack:

main
 \
  branch-2
   \
*   branch-1
     \
      branch-3
`

const swapCommandName = "swap"

func Cmd() *cobra.Command {
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     swapCommandName,
		Args:    cobra.NoArgs,
		Short:   swapDesc,
		GroupID: cmdhelpers.GroupIDStack,
		Long:    cmdhelpers.Long(swapDesc, swapHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAutoResolve, errDryRun, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       autoResolve,
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          Some(configdomain.Detached(true)),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeSwap(cliConfig)
		},
	}
	addAutoResolveFlag(&cmd)
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwap(cliConfig configdomain.PartialConfig) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := determineSwapData(repo)
	if err != nil {
		return err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit:
		return nil
	case configdomain.ProgramFlowRestart:
		goto Start
	}
	err = validateSwapData(data)
	if err != nil {
		return err
	}
	runProgram := swapProgram(repo, data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               swapCommandName,
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		Inputs:                  data.inputs,
		PendingCommand:          None[string](),
		RootDir:                 repo.RootDir,
		RunState:                runState,
	})
}

type swapData struct {
	branchInfosLastRun    Option[gitdomain.BranchInfos]
	branchesSnapshot      gitdomain.BranchesSnapshot
	children              []swapBranch
	config                config.ValidatedConfig
	connector             Option[forgedomain.Connector]
	currentBranchInfo     gitdomain.BranchInfo
	currentBranchName     gitdomain.LocalBranchName
	currentBranchProposal Option[forgedomain.Proposal]
	currentBranchType     configdomain.BranchType
	grandParentBranch     gitdomain.LocalBranchName
	hasOpenChanges        bool
	initialBranch         gitdomain.LocalBranchName
	inputs                dialogcomponents.Inputs
	nonExistingBranches   gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	parentBranch          gitdomain.LocalBranchName
	parentBranchInfo      gitdomain.BranchInfo
	parentBranchProposal  Option[forgedomain.Proposal]
	parentBranchType      configdomain.BranchType
	previousBranch        Option[gitdomain.LocalBranchName]
	stashSize             gitdomain.StashSize
}

type swapBranch struct {
	info     gitdomain.BranchInfo
	name     gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

func determineSwapData(repo execute.OpenRepoResult) (data swapData, flow configdomain.ProgramFlow, err error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
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
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitLabConnectorType:  config.GitLabConnectorType,
		GitLabToken:          config.GitLabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
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
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return data, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return data, configdomain.ProgramFlowExit, errors.New(messages.SwapRepoHasDetachedHead)
	}
	currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
	if !hasCurrentBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	currentBranchInfo, hasBranchToSwapInfo := branchesSnapshot.Branches.FindByLocalName(currentBranch).Get()
	if !hasBranchToSwapInfo {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchDoesntExist, currentBranch)
	}
	if currentBranchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchOtherWorktree, currentBranch)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{},
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
		return data, configdomain.ProgramFlowExit, err
	}
	currentBranchType := validatedConfig.BranchType(currentBranch)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(currentBranch).Get()
	if !hasParentBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.SwapNoParent)
	}
	parentBranchInfo, hasParentBranchInfo := branchesSnapshot.Branches.FindByLocalName(parentBranch).Get()
	if !hasParentBranchInfo {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.SwapParentNotLocal, parentBranch)
	}
	parentBranchType := validatedConfig.BranchType(parentBranch)
	grandParentBranch, hasGrandParentBranch := validatedConfig.NormalConfig.Lineage.Parent(parentBranch).Get()
	if !hasGrandParentBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.SwapNoGrandParent)
	}
	childBranches := validatedConfig.NormalConfig.Lineage.Children(currentBranch, validatedConfig.NormalConfig.Order)
	children := make([]swapBranch, len(childBranches))
	for c, childBranch := range childBranches {
		proposal := None[forgedomain.Proposal]()
		if connector, hasConnector := connector.Get(); hasConnector {
			if proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder); canFindProposals {
				proposal, err = proposalFinder.FindProposal(childBranch, initialBranch)
				if err != nil {
					return data, configdomain.ProgramFlowExit, err
				}
			}
		}
		childInfo, has := branchesSnapshot.Branches.FindByLocalName(childBranch).Get()
		if !has {
			return data, configdomain.ProgramFlowExit, fmt.Errorf("cannot find branch info for %q", childBranch)
		}
		children[c] = swapBranch{
			info:     *childInfo,
			name:     childBranch,
			proposal: proposal,
		}
	}
	currentbranchProposal := None[forgedomain.Proposal]()
	parentBranchProposal := None[forgedomain.Proposal]()
	if connector, hasConnector := connector.Get(); hasConnector {
		// TODO: load these two proposals concurrently
		if proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder); canFindProposals {
			currentbranchProposal, err = proposalFinder.FindProposal(currentBranch, parentBranch)
			if err != nil {
				return data, configdomain.ProgramFlowExit, err
			}
			parentBranchProposal, err = proposalFinder.FindProposal(parentBranch, grandParentBranch)
			if err != nil {
				return data, configdomain.ProgramFlowExit, err
			}
		}
	}
	branchContainsMerges, err := repo.Git.BranchContainsMerges(repo.Backend, currentBranch, parentBranch)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	if branchContainsMerges {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.SwapNeedsCompress, currentBranch)
	}
	parentContainsMerges, err := repo.Git.BranchContainsMerges(repo.Backend, parentBranch, grandParentBranch)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	if parentContainsMerges {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.SwapNeedsCompress, parentBranch)
	}
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(lineageBranches...)
	return swapData{
		branchInfosLastRun:    branchInfosLastRun,
		branchesSnapshot:      branchesSnapshot,
		children:              children,
		config:                validatedConfig,
		connector:             connector,
		currentBranchInfo:     *currentBranchInfo,
		currentBranchName:     currentBranch,
		currentBranchProposal: currentbranchProposal,
		currentBranchType:     currentBranchType,
		grandParentBranch:     grandParentBranch,
		hasOpenChanges:        repoStatus.OpenChanges,
		initialBranch:         initialBranch,
		inputs:                inputs,
		nonExistingBranches:   nonExistingBranches,
		parentBranch:          parentBranch,
		parentBranchInfo:      *parentBranchInfo,
		parentBranchProposal:  parentBranchProposal,
		parentBranchType:      parentBranchType,
		previousBranch:        previousBranchOpt,
		stashSize:             stashSize,
	}, configdomain.ProgramFlowContinue, nil
}

func swapProgram(repo execute.OpenRepoResult, data swapData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages, repo.Frontend, data.config.NormalConfig.Order)
	swapGitOperationsProgram(swapGitOperationsProgramArgs{
		children: data.children,
		current: swapBranch{
			info:     data.currentBranchInfo,
			name:     data.currentBranchName,
			proposal: data.currentBranchProposal,
		},
		grandParent: data.grandParentBranch,
		parent: swapBranch{
			info:     data.parentBranchInfo,
			name:     data.parentBranch,
			proposal: data.parentBranchProposal,
		},
		program: prog,
	})
	if !data.config.NormalConfig.DryRun {
		swapLineageParentSetsProgram(swapLineageParentSetsProgramArg{
			children:    data.children,
			current:     data.currentBranchName,
			grandParent: data.grandParentBranch,
			parent:      data.parentBranch,
			program:     prog,
		})
	}
	if data.config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		sync.AddSyncProposalsProgram(sync.AddSyncProposalsProgramArgs{
			ChangedBranches: gitdomain.LocalBranchNames{data.initialBranch},
			Config:          data.config,
			Program:         prog,
		})
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         false, // TODO: stash if open changes here?
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch},
	})
	return prog.Immutable()
}

func validateSwapData(data swapData) error {
	switch data.currentBranchInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusAhead, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusBehind:
		return errors.New(messages.SwapNeedsSync)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.SwapOtherWorkTree, data.currentBranchName)
	case gitdomain.SyncStatusRemoteOnly:
		return errors.New(messages.SwapRemoteBranch)
	}
	switch data.parentBranchInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusAhead, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusBehind:
		return errors.New(messages.SwapNeedsSync)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.SwapOtherWorkTree, data.parentBranch)
	case gitdomain.SyncStatusRemoteOnly:
		return fmt.Errorf(messages.SwapRemoteBranch, data.parentBranch)
	}
	switch data.currentBranchType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.SwapUnsupportedBranchType, data.currentBranchName, data.currentBranchType)
	}
	switch data.parentBranchType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.SwapParentWrongBranchType, data.parentBranch, data.parentBranchType)
	default:
		panic(fmt.Sprintf("unexpected configdomain.BranchType: %#v", data.parentBranchType))
	}
	for _, child := range data.children {
		switch child.info.SyncStatus {
		case gitdomain.SyncStatusAhead, gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusUpToDate:
		case gitdomain.SyncStatusBehind, gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusRemoteOnly:
			return errors.New(messages.SwapNeedsSync)
		case gitdomain.SyncStatusOtherWorktree:
			return fmt.Errorf(messages.SwapOtherWorkTree, child.name)
		}
	}
	return nil
}
