package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/cmd/ship"
	"github.com/git-town/git-town/v22/internal/cmd/sync"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/optimizer"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	mergeCmd  = "merge"
	mergeDesc = "Combines the current branch with its parent"
	mergeHelp = `
Merges the current branch with its parent branch.
Both branches must be feature branches.

Consider this stack:

main
 \
  branch-1
   \
    branch-2
     \
*     branch-3
       \
        branch-4

We are on the "branch-3" branch.
After running "git town merge",
the new "branch-3" branch contains the changes
from the old "branch-2" and "branch-3" branches.

main
 \
  branch-1
   \
*   branch-2
     \
      branch-4
`
)

func mergeCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     mergeCmd,
		Args:    cobra.NoArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   mergeDesc,
		Long:    cmdhelpers.Long(mergeDesc, mergeHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			dryRun, errDryRun := readDryRunFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errDryRun, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
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
			return executeMerge(cliConfig)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeMerge(cliConfig configdomain.PartialConfig) error {
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
	data, flow, err := determineMergeData(repo)
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
	if err = validateMergeData(repo, data); err != nil {
		return err
	}
	runProgram := mergeProgram(repo, data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               mergeCmd,
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

type mergeData struct {
	branchInfosLastRun       Option[gitdomain.BranchInfos]
	branchesSnapshot         gitdomain.BranchesSnapshot
	childBranches            gitdomain.LocalBranchNames
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	initialBranchInfo        gitdomain.BranchInfo
	initialBranchSHA         gitdomain.SHA
	initialBranchType        configdomain.BranchType
	inputs                   dialogcomponents.Inputs
	parentBranch             gitdomain.LocalBranchName
	parentBranchInfo         gitdomain.BranchInfo
	parentBranchSHA          gitdomain.SHA
	parentBranchType         configdomain.BranchType
	previousBranch           Option[gitdomain.LocalBranchName]
	proposalsOfChildBranches []forgedomain.Proposal
	stashSize                gitdomain.StashSize
}

func determineMergeData(repo execute.OpenRepoResult) (data mergeData, flow configdomain.ProgramFlow, err error) {
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
		GitHubConnectorType:  config.GitHubConnectorType,
		GitHubToken:          config.GitHubToken,
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
		return data, configdomain.ProgramFlowExit, errors.New(messages.MergeDetachedHead)
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
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
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(initialBranch).Get()
	if !hasParentBranch {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.MergeNoParent, initialBranch)
	}
	grandParentBranch := validatedConfig.NormalConfig.Lineage.Parent(parentBranch)
	if grandParentBranch.IsNone() {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.MergeNoGrandParent, initialBranch, parentBranch)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchInfoNotFound, initialBranch)
	}
	initialBranchSHA, hasInitialBranchSHA := initialBranchInfo.LocalSHA.Get()
	if !hasInitialBranchSHA {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.MergeBranchNotLocal, initialBranch)
	}
	parentBranchInfo, hasParentBranchInfo := branchesSnapshot.Branches.FindByLocalName(parentBranch).Get()
	if !hasParentBranchInfo {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchInfoNotFound, parentBranch)
	}
	parentBranchSHA, hasParentBranchSHA := parentBranchInfo.LocalSHA.Get()
	if !hasParentBranchSHA {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.MergeBranchNotLocal, parentBranch)
	}
	initialBranchType := validatedConfig.BranchType(initialBranch)
	parentBranchType := validatedConfig.BranchType(parentBranch)

	childBranches := validatedConfig.NormalConfig.Lineage.Children(initialBranch, validatedConfig.NormalConfig.Order)
	proposalsOfChildBranches := ship.LoadProposalsOfChildBranches(ship.LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connector,
		Lineage:                    validatedConfig.NormalConfig.Lineage,
		Offline:                    repo.IsOffline,
		OldBranch:                  initialBranch,
		OldBranchHasTrackingBranch: branchesSnapshot.Branches.FindByLocalName(initialBranch).IsSome(),
		Order:                      validatedConfig.NormalConfig.Order,
	})
	return mergeData{
		branchInfosLastRun:       branchInfosLastRun,
		branchesSnapshot:         branchesSnapshot,
		childBranches:            childBranches,
		config:                   validatedConfig,
		connector:                connector,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		initialBranchInfo:        *initialBranchInfo,
		initialBranchSHA:         initialBranchSHA,
		initialBranchType:        initialBranchType,
		inputs:                   inputs,
		parentBranch:             parentBranch,
		parentBranchInfo:         *parentBranchInfo,
		parentBranchSHA:          parentBranchSHA,
		parentBranchType:         parentBranchType,
		previousBranch:           previousBranch,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, configdomain.ProgramFlowContinue, err
}

func mergeProgram(repo execute.OpenRepoResult, data mergeData) program.Program {
	prog := NewMutable(&program.Program{})
	ship.UpdateChildBranchProposalsToGrandParent(prog.Value, data.proposalsOfChildBranches)
	prog.Value.Add(&opcodes.Checkout{Branch: data.parentBranch})
	if data.initialBranchSHA != data.parentBranchSHA {
		prog.Value.Add(&opcodes.BranchCurrentResetToSHA{SHA: data.initialBranchSHA})
	}
	for _, child := range data.childBranches {
		prog.Value.Add(&opcodes.LineageParentSetToGrandParent{Branch: child})
	}
	prog.Value.Add(&opcodes.LineageParentRemove{
		Branch: data.initialBranch,
	})
	// We need to delete the remote branch before updating the parent branch.
	// This order will make sure that if there is a related proposal where
	// the base is the branch being merged and the target is the parent branch
	// then that proposal gets marked "closed" rather than "merged" by forges
	// like GitHub.
	initialTrackingBranch, initialHasTrackingBranch := data.initialBranchInfo.RemoteName.Get()
	if initialHasTrackingBranch && repo.IsOffline.IsOnline() {
		prog.Value.Add(&opcodes.BranchTrackingDelete{
			Branch: initialTrackingBranch,
		})
	}
	prog.Value.Add(&opcodes.BranchLocalDelete{
		Branch: data.initialBranch,
	})
	parentTracking, parentHasTracking := data.parentBranchInfo.RemoteName.Get()
	if parentHasTracking && repo.IsOffline.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{
			CurrentBranch:   data.parentBranch,
			ForceIfIncludes: true,
			TrackingBranch:  parentTracking,
		})
	}
	if _, hasOverride := data.config.NormalConfig.BranchTypeOverrides[data.initialBranch]; hasOverride {
		prog.Value.Add(&opcodes.BranchTypeOverrideRemove{
			Branch: data.initialBranch,
		})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	if data.config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		_ = sync.AddSyncProposalsProgram(sync.AddSyncProposalsProgramArgs{
			Current:   data.initialBranch,
			FullStack: true,
			Program:   prog,
			ProposalStackLineageArgs: proposallineage.ProposalStackLineageArgs{
				Connector:                forgedomain.ProposalFinderFromConnector(data.connector),
				CurrentBranch:            data.initialBranch,
				Lineage:                  data.config.NormalConfig.Lineage,
				MainAndPerennialBranches: data.config.MainAndPerennials(),
				Order:                    data.config.NormalConfig.Order,
			},
			SkipUpdateForProposalsWithBaseBranch: gitdomain.LocalBranchNames{data.initialBranch},
		})
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

func validateMergeData(repo execute.OpenRepoResult, data mergeData) error {
	if err := verifyBranchType(data.initialBranchType); err != nil {
		return err
	}
	if err := verifyBranchType(data.parentBranchType); err != nil {
		return err
	}
	// ensure all commits on parent branch are contained in the initial branch
	inSyncWithParent, err := repo.Git.BranchInSyncWithParent(repo.Backend, data.initialBranch, data.parentBranch.BranchName())
	if err != nil {
		return err
	}
	if !inSyncWithParent {
		return fmt.Errorf(messages.BranchNotInSyncWithParent, data.initialBranch)
	}
	switch data.initialBranchInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusAhead, gitdomain.SyncStatusBehind, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.MergeNotInSyncWithTracking, data.initialBranch)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.BranchOtherWorktree, data.parentBranch)
	case gitdomain.SyncStatusRemoteOnly:
		// safe to ignore, this cannot happen
	}
	switch data.parentBranchInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusRemoteOnly:
	case gitdomain.SyncStatusAhead, gitdomain.SyncStatusBehind, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusDeletedAtRemote:
		return fmt.Errorf(messages.BranchNotInSyncWithParent, data.parentBranch)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.BranchOtherWorktree, data.parentBranch)
	}
	children := data.config.NormalConfig.Lineage.Children(data.parentBranch, data.config.NormalConfig.Order)
	if len(children) > 1 {
		return fmt.Errorf("branch %q has more than one child", data.parentBranch)
	}
	return nil
}

func verifyBranchType(branchType configdomain.BranchType) error {
	switch branchType {
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.MergeWrongBranchType, branchType)
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	}
	return nil
}
