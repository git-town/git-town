package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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

func swapCommand() *cobra.Command {
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     swapCommandName,
		Args:    cobra.NoArgs,
		Short:   swapDesc,
		GroupID: cmdhelpers.GroupIDStack,
		Long:    cmdhelpers.Long(swapDesc, swapHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAutoResolve, errDryRun, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: autoResolve,
				DryRun:      dryRun,
				Verbose:     verbose,
			})
			return executeSwap(args, cliConfig)
		},
	}
	addAutoResolveFlag(&cmd)
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwap(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSwapData(args, repo)
	if err != nil || exit {
		return err
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
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
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
		Detached:                true,
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
	branchInfosLastRun  Option[gitdomain.BranchInfos]
	branchToSwapInfo    gitdomain.BranchInfo
	branchToSwapName    gitdomain.LocalBranchName
	branchToSwapType    configdomain.BranchType
	branchesSnapshot    gitdomain.BranchesSnapshot
	children            []swapChildBranch
	config              config.ValidatedConfig
	connector           Option[forgedomain.Connector]
	grandParentBranch   gitdomain.LocalBranchName
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	inputs              dialogcomponents.Inputs
	nonExistingBranches gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	parentBranch        gitdomain.LocalBranchName
	parentBranchInfo    gitdomain.BranchInfo
	parentBranchType    configdomain.BranchType
	previousBranch      Option[gitdomain.LocalBranchName]
	stashSize           gitdomain.StashSize
}

type swapChildBranch struct {
	info     gitdomain.BranchInfo
	name     gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

func determineSwapData(args []string, repo execute.OpenRepoResult) (data swapData, exit dialogdomain.Exit, err error) {
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
		CodebergToken:        config.CodebergToken,
		ForgeType:            config.ForgeType,
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
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Detached:              true,
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
	if err != nil || exit {
		return data, exit, err
	}
	branchNameToSwap := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToSwapInfo, hasBranchToSwapInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToSwap).Get()
	if !hasBranchToSwapInfo {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToSwap)
	}
	if branchToSwapInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, exit, fmt.Errorf(messages.BranchOtherWorktree, branchNameToSwap)
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
		return data, exit, err
	}
	branchTypeToSwap := validatedConfig.BranchType(branchNameToSwap)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(branchNameToSwap).Get()
	if !hasParentBranch {
		return data, false, errors.New(messages.SwapNoParent)
	}
	parentBranchInfo, hasParentBranchInfo := branchesSnapshot.Branches.FindByLocalName(parentBranch).Get()
	if !hasParentBranchInfo {
		return data, false, fmt.Errorf(messages.SwapParentNotLocal, parentBranch)
	}
	parentBranchType := validatedConfig.BranchType(parentBranch)
	grandParentBranch, hasGrandParentBranch := validatedConfig.NormalConfig.Lineage.Parent(parentBranch).Get()
	if !hasGrandParentBranch {
		return data, false, errors.New(messages.SwapNoGrandParent)
	}
	childBranches := validatedConfig.NormalConfig.Lineage.Children(branchNameToSwap)
	children := make([]swapChildBranch, len(childBranches))
	for c, childBranch := range childBranches {
		proposal := None[forgedomain.Proposal]()
		if connector, hasConnector := connector.Get(); hasConnector {
			if findProposal, canFindProposal := connector.FindProposalFn().Get(); canFindProposal {
				proposal, err = findProposal(childBranch, initialBranch)
				if err != nil {
					return data, false, err
				}
			}
		}
		childInfo, has := branchesSnapshot.Branches.FindByLocalName(childBranch).Get()
		if !has {
			return data, false, fmt.Errorf("cannot find branch info for %q", childBranch)
		}
		children[c] = swapChildBranch{
			info:     *childInfo,
			name:     childBranch,
			proposal: proposal,
		}
	}
	branchContainsMerges, err := repo.Git.BranchContainsMerges(repo.Backend, branchNameToSwap, parentBranch)
	if err != nil {
		return data, false, err
	}
	if branchContainsMerges {
		return data, false, fmt.Errorf(messages.SwapNeedsCompress, branchNameToSwap)
	}
	parentContainsMerges, err := repo.Git.BranchContainsMerges(repo.Backend, parentBranch, grandParentBranch)
	if err != nil {
		return data, false, err
	}
	if parentContainsMerges {
		return data, false, fmt.Errorf(messages.SwapNeedsCompress, parentBranch)
	}
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	return swapData{
		branchInfosLastRun:  branchInfosLastRun,
		branchToSwapInfo:    *branchToSwapInfo,
		branchToSwapName:    branchNameToSwap,
		branchToSwapType:    branchTypeToSwap,
		branchesSnapshot:    branchesSnapshot,
		children:            children,
		config:              validatedConfig,
		connector:           connector,
		grandParentBranch:   grandParentBranch,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		inputs:              inputs,
		nonExistingBranches: nonExistingBranches,
		parentBranch:        parentBranch,
		parentBranchInfo:    *parentBranchInfo,
		parentBranchType:    parentBranchType,
		previousBranch:      previousBranchOpt,
		stashSize:           stashSize,
	}, false, nil
}

func swapProgram(repo execute.OpenRepoResult, data swapData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages, repo.Frontend)
	prog.Value.Add(
		&opcodes.RebaseOntoKeepDeleted{
			BranchToRebaseOnto: data.grandParentBranch.BranchName(),
			CommitsToRemove:    data.parentBranch.Location(),
			Upstream:           None[gitdomain.LocalBranchName](),
		},
	)
	if data.branchToSwapInfo.HasTrackingBranch() {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: data.branchToSwapName, ForceIfIncludes: true},
		)
	}
	prog.Value.Add(
		&opcodes.Checkout{
			Branch: data.parentBranch,
		},
		&opcodes.RebaseOntoKeepDeleted{
			BranchToRebaseOnto: data.branchToSwapName.BranchName(),
			CommitsToRemove:    data.grandParentBranch.Location(),
			Upstream:           None[gitdomain.LocalBranchName](),
		},
	)
	if data.parentBranchInfo.HasTrackingBranch() {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: data.parentBranch, ForceIfIncludes: true},
		)
	}
	for _, child := range data.children {
		prog.Value.Add(
			&opcodes.Checkout{
				Branch: child.name,
			},
		)
		oldBranchSHA, hasOldBranchSHA := data.branchToSwapInfo.LocalSHA.Get()
		if !hasOldBranchSHA {
			oldBranchSHA = data.branchToSwapInfo.RemoteSHA.GetOrDefault()
		}
		prog.Value.Add(
			&opcodes.RebaseOntoKeepDeleted{
				BranchToRebaseOnto: data.parentBranch.BranchName(),
				CommitsToRemove:    oldBranchSHA.Location(),
				Upstream:           None[gitdomain.LocalBranchName](),
			},
		)
		if child.info.HasTrackingBranch() {
			prog.Value.Add(
				&opcodes.PushCurrentBranchForceIfNeeded{
					CurrentBranch:   child.name,
					ForceIfIncludes: true,
				},
			)
		}
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.branchToSwapName})
	if !data.config.NormalConfig.DryRun {
		prog.Value.Add(
			&opcodes.LineageParentSet{
				Branch: data.branchToSwapName,
				Parent: data.grandParentBranch,
			},
			&opcodes.LineageParentSet{
				Branch: data.parentBranch,
				Parent: data.branchToSwapName,
			},
		)
		for _, child := range data.children {
			prog.Value.Add(
				&opcodes.LineageParentSet{
					Branch: child.name,
					Parent: data.parentBranch,
				},
			)
		}
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         false,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch},
	})
	return prog.Immutable()
}

func validateSwapData(data swapData) error {
	switch data.branchToSwapInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusAhead, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusBehind:
		return errors.New(messages.SwapNeedsSync)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.SwapOtherWorkTree, data.branchToSwapName)
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
	switch data.branchToSwapType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
	case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.SwapUnsupportedBranchType, data.branchToSwapName, data.branchToSwapType)
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
