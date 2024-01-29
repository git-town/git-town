package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/sync"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/spf13/cobra"
)

const syncDesc = "Updates the current branch with all relevant changes"

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
		Use:     "sync",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   syncDesc,
		Long:    cmdhelpers.Long(syncDesc, fmt.Sprintf(syncHelp, gitconfig.KeySyncUpstream)),
		RunE: func(cmd *cobra.Command, args []string) error {
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
		Verbose:          verbose,
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineSyncConfig(all, repo, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := program.Program{}
	sync.BranchesProgram(sync.BranchesProgramArgs{
		BranchProgramArgs: sync.BranchProgramArgs{
			Config:      config.FullConfig,
			BranchInfos: config.allBranches,
			Remotes:     config.remotes,
			Program:     &runProgram,
			PushBranch:  true,
		},
		BranchesToSync: config.branchesToSync,
		DryRun:         dryRun,
		HasOpenChanges: config.hasOpenChanges,
		InitialBranch:  config.initialBranch,
		PreviousBranch: config.previousBranch,
		ShouldPushTags: config.shouldPushTags,
	})
	runState := runstate.RunState{
		Command:             "sync",
		DryRun:              dryRun,
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          runProgram,
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		FullConfig:              config.FullConfig,
		RunState:                &runState,
		Run:                     repo.Runner,
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		Verbose:                 verbose,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type syncConfig struct {
	*configdomain.FullConfig
	allBranches      gitdomain.BranchInfos
	branchesToSync   gitdomain.BranchInfos
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	previousBranch   gitdomain.LocalBranchName
	remotes          gitdomain.Remotes
	shouldPushTags   bool
}

func determineSyncConfig(allFlag bool, repo *execute.OpenRepoResult, verbose bool) (*syncConfig, gitdomain.BranchesStatus, gitdomain.StashSize, bool, error) {
	branchesSnapshot, stashSnapshot, dialogTestInputs, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	remotes, err := repo.Runner.Backend.Remotes()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	var branchNamesToSync gitdomain.LocalBranchNames
	var shouldPushTags bool
	if allFlag {
		localBranches := branchesSnapshot.Branches.LocalBranches()
		err = execute.EnsureKnownBranchesAncestry(execute.EnsureKnownBranchesAncestryArgs{
			Config:           &repo.Runner.FullConfig,
			LocalBranches:    localBranches,
			DialogTestInputs: &dialogTestInputs,
			Runner:           repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
		branchNamesToSync = localBranches.Names()
		shouldPushTags = true
	} else {
		err = execute.EnsureKnownBranchAncestry(branchesSnapshot.Active, execute.EnsureKnownBranchAncestryArgs{
			Config:           &repo.Runner.FullConfig,
			AllBranches:      branchesSnapshot.Branches,
			DefaultBranch:    repo.Runner.MainBranch,
			DialogTestInputs: &dialogTestInputs,
			Runner:           repo.Runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSnapshot, false, err
		}
		branchNamesToSync = gitdomain.LocalBranchNames{branchesSnapshot.Active}
		shouldPushTags = !repo.Runner.IsFeatureBranch(branchesSnapshot.Active)
	}
	allBranchNamesToSync := repo.Runner.Lineage.BranchesAndAncestors(branchNamesToSync)
	branchesToSync, err := branchesSnapshot.Branches.Select(allBranchNamesToSync)
	return &syncConfig{
		FullConfig:       &repo.Runner.FullConfig,
		allBranches:      branchesSnapshot.Branches,
		branchesToSync:   branchesToSync,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		previousBranch:   previousBranch,
		remotes:          remotes,
		shouldPushTags:   shouldPushTags,
	}, branchesSnapshot, stashSnapshot, false, err
}
