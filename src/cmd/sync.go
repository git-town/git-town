package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/optimizer"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const syncCommand = "sync"

const syncDesc = "Update the current branch with all relevant changes"

const syncHelp = `
Synchronizes the current branch with the rest of the world.

When run on a feature branch:
- syncs all ancestor branches
- pulls updates for the current branch
- merges the parent branch into the current branch
- pushes the current branch

When run on the main branch or a perennial branch:
- pulls and pushes updates for the current branch
- pushes tags

If the repository contains an "upstream" remote, syncs the main branch with its upstream counterpart. You can disable this by running "git config %s false".`

func syncCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addAllFlag, readAllFlag := flags.All()
	addNoPushFlag, readNoPushFlag := flags.NoPush()
	cmd := cobra.Command{
		Use:     syncCommand,
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   syncDesc,
		Long:    cmdhelpers.Long(syncDesc, fmt.Sprintf(syncHelp, gitconfig.KeySyncUpstream)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeSync(readAllFlag(cmd), readDryRunFlag(cmd), readVerboseFlag(cmd), readNoPushFlag(cmd))
		},
	}
	addAllFlag(&cmd)
	addVerboseFlag(&cmd)
	addDryRunFlag(&cmd)
	addNoPushFlag(&cmd)
	return &cmd
}

func executeSync(all bool, dryRun configdomain.DryRun, verbose configdomain.Verbose, pushBranches configdomain.PushBranches) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSyncData(all, repo, verbose)
	if err != nil || exit {
		return err
	}

	// remove outdated lineage
	if err = data.config.RemoveOutdatedConfiguration(data.allBranches.LocalBranches().Names()); err != nil {
		return err
	}
	if err = cleanupPerennialParentEntries(data.config.Config.Lineage, data.config.Config.PerennialBranches, data.config.GitConfig, repo.FinalMessages); err != nil {
		return err
	}

	runProgram := program.Program{}
	sync.BranchesProgram(sync.BranchesProgramArgs{
		BranchProgramArgs: sync.BranchProgramArgs{
			BranchInfos:   data.allBranches,
			Config:        data.config.Config,
			InitialBranch: data.initialBranch,
			Remotes:       data.remotes,
			Program:       NewMutable(&runProgram),
			PushBranch:    pushBranches,
		},
		BranchesToSync: data.branchesToSync,
		DryRun:         dryRun,
		HasOpenChanges: data.hasOpenChanges,
		InitialBranch:  data.initialBranch,
		PreviousBranch: data.previousBranch,
		ShouldPushTags: data.shouldPushTags,
	})
	runProgram = optimizer.Optimize(runProgram)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        0,
		Command:               syncCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               None[hostingdomain.Connector](),
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type syncData struct {
	allBranches      gitdomain.BranchInfos
	branchesSnapshot gitdomain.BranchesSnapshot
	branchesToSync   gitdomain.BranchInfos
	config           config.ValidatedConfig
	dialogTestInputs Mutable[components.TestInputs]
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	previousBranch   Option[gitdomain.LocalBranchName]
	remotes          gitdomain.Remotes
	shouldPushTags   bool
	stashSize        gitdomain.StashSize
}

func determineSyncData(allFlag bool, repo execute.OpenRepoResult, verbose configdomain.Verbose) (data syncData, exit bool, err error) {
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
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	previousBranch, hasPreviousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend).Get()
	var previousBranchOpt Option[gitdomain.LocalBranchName]
	if hasPreviousBranch {
		if previousBranchInfo, hasPreviousBranchInfo := branchesSnapshot.Branches.FindByLocalName(previousBranch).Get(); hasPreviousBranchInfo {
			switch previousBranchInfo.SyncStatus {
			case gitdomain.SyncStatusLocalOnly, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusUpToDate:
				previousBranchOpt = previousBranchInfo.LocalName
			case gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusRemoteOnly, gitdomain.SyncStatusOtherWorktree:
				previousBranchOpt = None[gitdomain.LocalBranchName]()
			}
		}
	} else {
		previousBranchOpt = None[gitdomain.LocalBranchName]()
	}
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	var branchNamesToSync gitdomain.LocalBranchNames
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	if allFlag {
		branchNamesToSync = localBranches
	} else {
		branchNamesToSync = gitdomain.LocalBranchNames{initialBranch}
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchNamesToSync,
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
	var shouldPushTags bool
	if allFlag {
		shouldPushTags = true
	} else {
		shouldPushTags = validatedConfig.Config.IsMainOrPerennialBranch(initialBranch)
	}
	allBranchNamesToSync := validatedConfig.Config.Lineage.BranchesAndAncestors(branchNamesToSync)
	branchesToSync, err := branchesSnapshot.Branches.Select(allBranchNamesToSync...)
	return syncData{
		allBranches:      branchesSnapshot.Branches,
		branchesSnapshot: branchesSnapshot,
		branchesToSync:   branchesToSync,
		config:           validatedConfig,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    initialBranch,
		previousBranch:   previousBranchOpt,
		remotes:          remotes,
		shouldPushTags:   shouldPushTags,
		stashSize:        stashSize,
	}, false, err
}

// cleanupPerennialParentEntries removes outdated entries from the configuration.
func cleanupPerennialParentEntries(lineage configdomain.Lineage, perennialBranches gitdomain.LocalBranchNames, access gitconfig.Access, finalMessages stringslice.Collector) error {
	for _, perennialBranch := range perennialBranches {
		if lineage.Parent(perennialBranch).IsSome() {
			if err := access.RemoveLocalConfigValue(gitconfig.NewParentKey(perennialBranch)); err != nil {
				return err
			}
			lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
	return nil
}
