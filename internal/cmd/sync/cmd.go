// Package sync provides functionality around syncing Git branches.
package sync

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
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
	"github.com/git-town/git-town/v22/pkg/set"
	"github.com/spf13/cobra"
)

const (
	syncCommand = "sync"
	syncDesc    = "Update the current branch with all relevant changes"
	syncHelp    = `
Synchronizes the current branch with the rest of the world.

When run on a feature branch:
- syncs all ancestor branches
- pulls updates for the current branch
- merges the parent branch into the current branch
- pushes the current branch
- updates branch lineage in proposals

When run on the main branch or a perennial branch:
- pulls and pushes updates for the current branch
- pushes tags

If the repository contains an "upstream" remote, syncs the main branch with its upstream counterpart. You can disable this by running "git config %s false".`
)

func Cmd() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("sync all local branches")
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addGoneFlag, readGoneFlag := flags.Gone()
	addPruneFlag, readPruneFlag := flags.Prune()
	addPushFlag, readPushFlag := flags.Push()
	addStackFlag, readStackFlag := flags.Stack("sync the stack that the current branch belongs to")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     syncCommand,
		GroupID: cmdhelpers.GroupIDBasic,
		Args:    cobra.NoArgs,
		Short:   syncDesc,
		Long:    cmdhelpers.Long(syncDesc, fmt.Sprintf(syncHelp, configdomain.KeySyncUpstream)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			allBranches, errAllBranches := readAllFlag(cmd)
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			detached, errDetached := readDetachedFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			gone, errGone := readGoneFlag(cmd)
			prune, errPrune := readPruneFlag(cmd)
			pushBranches, errPushBranches := readPushFlag(cmd)
			stack, errStack := readStackFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAllBranches, errDetached, errDryRun, errAutoResolve, errGone, errPushBranches, errPrune, errStack, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       autoResolve,
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          detached,
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      pushBranches,
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeSync(executeSyncArgs{
				cliConfig:       cliConfig,
				gone:            gone,
				prune:           prune,
				stack:           stack,
				syncAllBranches: allBranches,
			})
		},
	}
	addAllFlag(&cmd)
	addAutoResolveFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addGoneFlag(&cmd)
	addPruneFlag(&cmd)
	addPushFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeSyncArgs struct {
	cliConfig       configdomain.PartialConfig
	gone            configdomain.Gone
	prune           configdomain.Prune
	stack           configdomain.FullStack
	syncAllBranches configdomain.AllBranches
}

func executeSync(args executeSyncArgs) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := determineSyncData(repo, determineSyncDataArgs{
		gone:            args.gone,
		syncAllBranches: args.syncAllBranches,
		syncStack:       args.stack,
	})
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
	if err = validateSyncData(data); err != nil {
		return err
	}
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, repo.FinalMessages, repo.Backend, data.config.NormalConfig.Order)
	runProgram := NewMutable(&program.Program{})
	branchesToDelete := set.New[gitdomain.LocalBranchName]()
	BranchesProgram(data.branchesToSync, BranchProgramArgs{
		BranchInfos:         data.branchInfos,
		BranchInfosPrevious: data.previousBranchInfos,
		BranchesToDelete:    NewMutable(&branchesToDelete),
		Config:              data.config,
		InitialBranch:       data.initialBranch,
		PrefetchBranchInfos: data.prefetchBranchesSnapshot.Branches,
		Program:             runProgram,
		Prune:               args.prune,
		PushBranches:        data.config.NormalConfig.PushBranches,
		Remotes:             data.remotes,
	})
	previousbranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	finalBranchCandidates := gitdomain.LocalBranchNames{data.initialBranch}
	if previousBranch, hasPreviousBranch := data.previousBranch.Get(); hasPreviousBranch {
		finalBranchCandidates = append(finalBranchCandidates, previousBranch)
	}
	finalBranchCandidates = finalBranchCandidates.AppendAllMissing(data.branchInfos.NamesLocalBranches())
	finalBranchCandidates = finalBranchCandidates.Remove(data.branchInfos.BranchesInOtherWorktrees()...)
	runProgram.Value.Add(&opcodes.CheckoutFirstExisting{
		Branches:   finalBranchCandidates,
		MainBranch: data.config.ValidatedConfigData.MainBranch,
	})
	if data.remotes.HasRemote(data.config.NormalConfig.DevRemote) && data.shouldPushTags && data.config.NormalConfig.Offline.IsOnline() {
		runProgram.Value.Add(&opcodes.PushTags{})
	}
	if data.config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		_ = AddOpcodesToUpdateProposalStack(
			AddOpcodesToUpdateProposalStackArgs{
				Current:   data.initialBranch,
				FullStack: args.stack,
				Program:   runProgram,
				ProposalStackLineageArgs: proposallineage.ProposalStackLineageArgs{
					Connector:                forgedomain.ProposalFinderFromConnector(data.connector),
					CurrentBranch:            data.initialBranch,
					Lineage:                  data.config.NormalConfig.Lineage,
					MainAndPerennialBranches: data.config.MainAndPerennials(),
					Order:                    data.config.NormalConfig.Order,
				},
				SkipUpdateForProposalsWithBaseBranch: gitdomain.NewLocalBranchNames(),
			},
		)
	}

	cmdhelpers.Wrap(runProgram, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousbranchCandidates,
	})
	optimizedProgram := optimizer.Optimize(runProgram.Immutable())
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        0,
		Command:               syncCommand,
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		BranchInfosLastRun:    data.previousBranchInfos,
		RunProgram:            optimizedProgram,
		TouchedBranches:       optimizedProgram.TouchedBranches(),
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

type syncData struct {
	branchInfos              gitdomain.BranchInfos
	branchesSnapshot         gitdomain.BranchesSnapshot
	branchesToSync           configdomain.BranchesToSync
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	inputs                   dialogcomponents.Inputs
	nonExistingBranches      gitdomain.LocalBranchNames
	prefetchBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch           Option[gitdomain.LocalBranchName]
	previousBranchInfos      Option[gitdomain.BranchInfos]
	remotes                  gitdomain.Remotes
	shouldPushTags           bool
	stashSize                gitdomain.StashSize
}

type determineSyncDataArgs struct {
	gone            configdomain.Gone
	syncAllBranches configdomain.AllBranches
	syncStack       configdomain.FullStack
}

func determineSyncData(repo execute.OpenRepoResult, args determineSyncDataArgs) (data syncData, flow configdomain.ProgramFlow, err error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	preFetchBranchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
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
		return data, configdomain.ProgramFlowExit, errors.New(messages.SyncRepoHasDetachedHead)
	}
	previousBranch, hasPreviousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend).Get()
	var previousBranchOpt Option[gitdomain.LocalBranchName]
	if hasPreviousBranch {
		if previousBranchInfo, hasPreviousBranchInfo := branchesSnapshot.Branches.FindByLocalName(previousBranch).Get(); hasPreviousBranchInfo {
			switch previousBranchInfo.SyncStatus {
			case
				gitdomain.SyncStatusLocalOnly,
				gitdomain.SyncStatusNotInSync,
				gitdomain.SyncStatusAhead,
				gitdomain.SyncStatusBehind,
				gitdomain.SyncStatusUpToDate:
				previousBranchOpt = previousBranchInfo.LocalName
			case
				gitdomain.SyncStatusDeletedAtRemote,
				gitdomain.SyncStatusRemoteOnly,
				gitdomain.SyncStatusOtherWorktree:
				previousBranchOpt = None[gitdomain.LocalBranchName]()
			}
		}
	} else {
		previousBranchOpt = None[gitdomain.LocalBranchName]()
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
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
	perennialAndMain := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToSync gitdomain.LocalBranchNames
	switch {
	case args.gone.Enabled():
		branchNamesToSync = branchesSnapshot.Branches.BranchesDeletedAtRemote()
	case args.syncAllBranches.Enabled() && repo.UnvalidatedConfig.NormalConfig.Detached.ShouldWorkDetached():
		branchNamesToSync = localBranches.Remove(perennialAndMain...)
	case args.syncAllBranches.Enabled():
		branchNamesToSync = localBranches
	case args.syncStack.Enabled():
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialAndMain, validatedConfig.NormalConfig.Order)
	default:
		branchNamesToSync = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesAndTypes = repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	validatedConfig, exit, err = validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: branchNamesToSync,
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
	var shouldPushTags bool
	switch {
	case !validatedConfig.NormalConfig.SyncTags.ShouldSyncTags():
		shouldPushTags = false
	case args.syncAllBranches.Enabled():
		shouldPushTags = true
	default:
		shouldPushTags = validatedConfig.IsMainOrPerennialBranch(initialBranch)
	}
	allBranchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchesAndAncestors(branchNamesToSync, validatedConfig.NormalConfig.Order)
	if repo.UnvalidatedConfig.NormalConfig.Detached {
		allBranchNamesToSync = allBranchNamesToSync.Remove(perennialAndMain...)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(allBranchNamesToSync...)
	branchesToSync, err := BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	return syncData{
		branchInfos:              branchesSnapshot.Branches,
		branchesSnapshot:         branchesSnapshot,
		branchesToSync:           branchesToSync,
		config:                   validatedConfig,
		connector:                connector,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		inputs:                   inputs,
		nonExistingBranches:      nonExistingBranches,
		prefetchBranchesSnapshot: preFetchBranchesSnapshot,
		previousBranch:           previousBranchOpt,
		previousBranchInfos:      previousBranchInfos,
		remotes:                  remotes,
		shouldPushTags:           shouldPushTags,
		stashSize:                stashSize,
	}, configdomain.ProgramFlowContinue, err
}

func BranchesToSync(branchInfosToSync gitdomain.BranchInfos, allBranchInfos gitdomain.BranchInfos, repo execute.OpenRepoResult, mainBranch gitdomain.LocalBranchName) (configdomain.BranchesToSync, error) {
	result := make(configdomain.BranchesToSync, len(branchInfosToSync))
	for b, branchInfo := range branchInfosToSync {
		branchNameToSync := branchInfo.GetLocalOrRemoteName()
		if branchNameToSync.LocalName() == mainBranch {
			result[b] = configdomain.BranchToSync{
				BranchInfo:         branchInfo,
				FirstCommitMessage: None[gitdomain.CommitMessage](),
			}
			continue
		}
		parentLocalName, hasParentName := repo.UnvalidatedConfig.NormalConfig.Lineage.Parent(branchNameToSync.LocalName()).Get()
		if !hasParentName {
			parentLocalName = mainBranch
		}
		parentBranchInfo, hasParentBranchInfo := allBranchInfos.FindLocalOrRemote(parentLocalName).Get()
		if !hasParentBranchInfo {
			result[b] = configdomain.BranchToSync{
				BranchInfo:         branchInfo,
				FirstCommitMessage: None[gitdomain.CommitMessage](),
			}
			continue
		}
		parentBranchName := parentBranchInfo.GetLocalOrRemoteName()
		firstCommitMessage, err := repo.Git.FirstCommitMessageInBranch(repo.Backend, branchNameToSync, parentBranchName)
		if err != nil {
			return result, err
		}
		result[b] = configdomain.BranchToSync{
			BranchInfo:         branchInfo,
			FirstCommitMessage: firstCommitMessage,
		}
	}
	return result, nil
}

func validateSyncData(data syncData) error {
	// ensure any branch that uses the ff-only sync strategy does not have unpushed local commits
	if data.config.NormalConfig.SyncPerennialStrategy == configdomain.SyncPerennialStrategyFFOnly {
		perennialBranchesToSync := data.config.BranchesOfType(data.branchesToSync.BranchNames(), configdomain.BranchTypePerennialBranch)
		for _, perennialBranchToSync := range perennialBranchesToSync {
			if branchInfo, hasBranchInfo := data.branchInfos.FindByLocalName(perennialBranchToSync).Get(); hasBranchInfo {
				switch branchInfo.SyncStatus {
				case gitdomain.SyncStatusAhead, gitdomain.SyncStatusNotInSync:
					return fmt.Errorf(messages.SyncPerennialBranchHasUnpushedCommits, perennialBranchToSync)
				case gitdomain.SyncStatusBehind, gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusOtherWorktree, gitdomain.SyncStatusRemoteOnly, gitdomain.SyncStatusUpToDate:
					// no problem with these sync statuses
				}
			}
		}
	}
	return nil
}
