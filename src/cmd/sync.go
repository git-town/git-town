package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
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
	addAllFlag, readAllFlag := flags.Bool("all", "a", "Sync all local branches", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:     syncCommand,
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   syncDesc,
		Long:    cmdhelpers.Long(syncDesc, fmt.Sprintf(syncHelp, gitconfig.KeySyncUpstream)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeSync(readAllFlag(cmd), readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addAllFlag(&cmd)
	addVerboseFlag(&cmd)
	addDryRunFlag(&cmd)
	return &cmd
}

func executeSync(all, dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineSyncData(all, repo, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := program.Program{}
	sync.BranchesProgram(sync.BranchesProgramArgs{
		BranchProgramArgs: sync.BranchProgramArgs{
			BranchInfos:   data.allBranches,
			Config:        data.config.Config,
			InitialBranch: data.initialBranch,
			Remotes:       data.remotes,
			Program:       &runProgram,
			PushBranch:    true,
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
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        0,
		Command:               syncCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            runProgram,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        &data.dialogTestInputs,
		FinalMessages:           &repo.FinalMessages,
		Frontend:                repo.Frontend,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type syncData struct {
	allBranches      gitdomain.BranchInfos
	branchesToSync   gitdomain.BranchInfos
	config           config.ValidatedConfig
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	previousBranch   gitdomain.LocalBranchName
	remotes          gitdomain.Remotes
	shouldPushTags   bool
}

func emptySyncData() syncData {
	return syncData{} //exhaustruct:ignore
}

func determineSyncData(allFlag bool, repo execute.OpenRepoResult, verbose bool) (syncData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return emptySyncData(), gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               &repo.Backend,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		Frontend:              &repo.Frontend,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return emptySyncData(), branchesSnapshot, stashSize, exit, err
	}
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	remotes, err := repo.Backend.Remotes()
	if err != nil {
		return emptySyncData(), branchesSnapshot, stashSize, false, err
	}
	var branchNamesToSync gitdomain.LocalBranchNames
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	if allFlag {
		branchNamesToSync = localBranches
	} else {
		branchNamesToSync = gitdomain.LocalBranchNames{branchesSnapshot.Active}
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            &repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchNamesToSync,
		CommandsCounter:    repo.CommandsCounter,
		ConfigSnapshot:     repo.ConfigSnapshot,
		DialogTestInputs:   dialogTestInputs,
		FinalMessages:      &repo.FinalMessages,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		RootDir:            repo.RootDir,
		StashSize:          stashSize,
		TestInputs:         &dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
		Verbose:            verbose,
	})
	if err != nil || exit {
		return emptySyncData(), branchesSnapshot, stashSize, exit, err
	}
	var shouldPushTags bool
	if allFlag {
		shouldPushTags = true
	} else {
		shouldPushTags = validatedConfig.Config.IsMainOrPerennialBranch(branchesSnapshot.Active)
	}
	allBranchNamesToSync := validatedConfig.Config.Lineage.BranchesAndAncestors(branchNamesToSync)
	branchesToSync, err := branchesSnapshot.Branches.Select(allBranchNamesToSync...)
	return syncData{
		allBranches:      branchesSnapshot.Branches,
		branchesToSync:   branchesToSync,
		config:           validatedConfig,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		previousBranch:   previousBranch,
		remotes:          remotes,
		shouldPushTags:   shouldPushTags,
	}, branchesSnapshot, stashSize, false, err
}
