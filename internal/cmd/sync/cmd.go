// Package sync provides functionality around syncing Git branches.
package sync

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"strings"

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
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/optimizer"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/git-town/git-town/v21/pkg/set"
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
- updates proposal(s)

When run on the main branch or a perennial branch:
- pulls and pushes updates for the current branch
- pushes tags

If the repository contains an "upstream" remote, syncs the main branch with its upstream counterpart. You can disable this by running "git config %s false".`
)

func Cmd() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("sync all local branches")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addNoPushFlag, readNoPushFlag := flags.NoPush()
	addPruneFlag, readPruneFlag := flags.Prune()
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
			detached, errDetached := readDetachedFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			pushBranches, errPushBranches := readNoPushFlag(cmd)
			prune, errPrune := readPruneFlag(cmd)
			stack, errStack := readStackFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAllBranches, errDetached, errDryRun, errAutoResolve, errPushBranches, errPrune, errStack, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: autoResolve,
				DryRun:      dryRun,
				Verbose:     verbose,
			})
			return executeSync(executeSyncArgs{
				cliConfig:       cliConfig,
				detached:        detached,
				prune:           prune,
				pushBranches:    pushBranches,
				stack:           stack,
				syncAllBranches: allBranches,
			})
		},
	}
	addAllFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addAutoResolveFlag(&cmd)
	addNoPushFlag(&cmd)
	addPruneFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeSyncArgs struct {
	cliConfig       configdomain.PartialConfig
	detached        configdomain.Detached
	prune           configdomain.Prune
	pushBranches    configdomain.PushBranches
	stack           configdomain.FullStack
	syncAllBranches configdomain.AllBranches
}

func executeSync(args executeSyncArgs) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSyncData(repo, determineSyncDataArgs{
		detached:        args.detached,
		syncAllBranches: args.syncAllBranches,
		syncStack:       args.stack,
	})
	if err != nil || exit {
		return err
	}
	if err = validateSyncData(data); err != nil {
		return err
	}
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, repo.FinalMessages, repo.Backend)
	runProgram := NewMutable(&program.Program{})
	branchesToDelete := set.New[gitdomain.LocalBranchName]()
	BranchesProgram(data.branchesToSync, BranchProgramArgs{
		BranchInfos:         data.branchInfos,
		BranchInfosLastRun:  data.previousBranchInfos,
		BranchesToDelete:    NewMutable(&branchesToDelete),
		Config:              data.config,
		InitialBranch:       data.initialBranch,
		PrefetchBranchInfos: data.prefetchBranchesSnapshot.Branches,
		Program:             runProgram,
		Prune:               args.prune,
		PushBranches:        args.pushBranches,
		Remotes:             data.remotes,
	})
	previousbranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	finalBranchCandidates := gitdomain.LocalBranchNames{data.initialBranch}
	if previousBranch, hasPreviousBranch := data.previousBranch.Get(); hasPreviousBranch {
		finalBranchCandidates = append(finalBranchCandidates, previousBranch)
	}
	runProgram.Value.Add(&opcodes.CheckoutFirstExisting{
		Branches:   finalBranchCandidates,
		MainBranch: data.config.ValidatedConfigData.MainBranch,
	})
	if data.remotes.HasRemote(data.config.NormalConfig.DevRemote) && data.shouldPushTags && data.config.NormalConfig.Offline.IsOnline() {
		runProgram.Value.Add(&opcodes.PushTags{})
	}

	if args.cliConfig.ProposalsShowLineage.GetOrDefault() == configdomain.ProposalsShowLineageCLI {
		explainationText := "Current dependencies on/for this pull-request\n\n"
		if data.config.NormalConfig.ForgeType.GetOrDefault() == forgedomain.ForgeTypeGitLab {
			explainationText = strings.ReplaceAll(explainationText, "pull-request", "merge-request")
		}
		proposalStackLineageArgs := configdomain.ProposalStackLineageArgs{
			AfterStackDisplay: []string{
				"\n\nStack generated by [Git-Town](https://github.com/git-town/git-town)\n",
			},
			BeforeStackDisplay: []string{
				"\n-------------------------\n",
				explainationText,
			},
			Connector:                data.connector,
			CurrentBranch:            data.initialBranch,
			CurrentBranchIndicator:   ":point_left:",
			IndentMarker:             "-",
			Lineage:                  data.config.NormalConfig.Lineage,
			MainAndPerennialBranches: Some(data.config.MainAndPerennials()),
		}
		proposalStackLineageBuilder := configdomain.NewProposalStackLineageBuilder(&proposalStackLineageArgs)
		if builder, hasBuilder := proposalStackLineageBuilder.Get(); hasBuilder {
			if branchProposalData, hasBranchProposalData := builder.GetProposal(data.initialBranch).Get(); hasBranchProposalData {
				runProgram.Value.Add(&opcodes.ProposalUpdateBody{
					Proposal:    branchProposalData,
					UpdatedBody: configdomain.ProposalBodyUpdateWithStackLineage(branchProposalData.Body.GetOrDefault(), builder.Build(&proposalStackLineageArgs).GetOrDefault()),
				})
			} else {
				fmt.Printf(messages.ProposalNotFoundForUpdate, data.initialBranch)
			}
		}
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
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
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
		Detached:                args.detached,
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
	detached                 configdomain.Detached
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
	detached        configdomain.Detached
	syncAllBranches configdomain.AllBranches
	syncStack       configdomain.FullStack
}

func determineSyncData(repo execute.OpenRepoResult, args determineSyncDataArgs) (data syncData, exit dialogdomain.Exit, err error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	preFetchBranchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
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
	branchesSnapshot, stashSize, previousBranchInfos, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Detached:              args.detached,
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
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
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
		return data, exit, err
	}
	perennialAndMain := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToSync gitdomain.LocalBranchNames
	switch {
	case args.syncAllBranches.Enabled() && args.detached.IsTrue():
		branchNamesToSync = localBranches.Remove(perennialAndMain...)
	case args.syncAllBranches.Enabled():
		branchNamesToSync = localBranches
	case args.syncStack.Enabled():
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialAndMain)
	default:
		branchNamesToSync = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesAndTypes = repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
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
		return data, exit, err
	}
	var shouldPushTags bool
	switch {
	case validatedConfig.NormalConfig.SyncTags.IsFalse():
		shouldPushTags = false
	case args.syncAllBranches.Enabled():
		shouldPushTags = true
	default:
		shouldPushTags = validatedConfig.IsMainOrPerennialBranch(initialBranch)
	}
	allBranchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchesAndAncestors(branchNamesToSync)
	if args.detached {
		allBranchNamesToSync = allBranchNamesToSync.Remove(perennialAndMain...)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, allBranchNamesToSync...)
	branchesToSync, err := BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	return syncData{
		branchInfos:              branchesSnapshot.Branches,
		branchesSnapshot:         branchesSnapshot,
		branchesToSync:           branchesToSync,
		config:                   validatedConfig,
		connector:                connector,
		detached:                 args.detached,
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
	}, false, err
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
		parentBranchInfo, hasParentBranchInfo := allBranchInfos.FindLocalOrRemote(parentLocalName, repo.UnvalidatedConfig.NormalConfig.DevRemote).Get()
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
