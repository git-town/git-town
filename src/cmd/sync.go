package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git"
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineSyncConfig(all, repo.UnvalidatedConfig, repo, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := program.Program{}
	sync.BranchesProgram(sync.BranchesProgramArgs{
		BranchProgramArgs: sync.BranchProgramArgs{
			BranchInfos:   config.allBranches,
			Config:        config.config,
			InitialBranch: config.initialBranch,
			Remotes:       config.remotes,
			Program:       &runProgram,
			PushBranch:    true,
		},
		BranchesToSync: config.branchesToSync,
		DryRun:         dryRun,
		HasOpenChanges: config.hasOpenChanges,
		InitialBranch:  config.initialBranch,
		PreviousBranch: config.previousBranch,
		ShouldPushTags: config.shouldPushTags,
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
		Config:                  config.config,
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		HasOpenChanges:          config.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     &prodRunner,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type syncConfig struct {
	allBranches      gitdomain.BranchInfos
	branchesToSync   gitdomain.BranchInfos
	config           configdomain.ValidatedConfig
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	previousBranch   gitdomain.LocalBranchName
	remotes          gitdomain.Remotes
	shouldPushTags   bool
}

func determineSyncConfig(allFlag bool, unvalidatedConfig configdomain.UnvalidatedConfig, repo *execute.OpenRepoResult, verbose bool) (*syncConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	validatedConfig, err := validate.ValidateConfig(repo.UnvalidatedConfig)
	runner := git.ProdRunner{
		Config:          validatedConfig,
		Backend:         repo.BackendCommands,
		Frontend:        repo.Frontend,
		CommandsCounter: repo.CommandsCounter,
		FinalMessages:   &repo.FinalMessages,
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                &unvalidatedConfig,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	previousBranch := runner.Backend.PreviouslyCheckedOutBranch()
	remotes, err := runner.Backend.Remotes()
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	var branchNamesToSync gitdomain.LocalBranchNames
	var shouldPushTags bool
	if allFlag {
		localBranches := branchesSnapshot.Branches.LocalBranches()
		err = execute.EnsureKnownBranchesAncestry(execute.EnsureKnownBranchesAncestryArgs{
			BranchesToVerify: branchesSnapshot.Branches.LocalBranches().Names(),
			Config:           runner.Config,
			DefaultChoice:    validatedConfig.FullConfig.MainBranch,
			DialogTestInputs: &dialogTestInputs,
			LocalBranches:    localBranches,
			MainBranch:       validatedConfig.FullConfig.MainBranch,
			Runner:           runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSize, false, err
		}
		branchNamesToSync = localBranches.Names()
		shouldPushTags = true
	} else {
		err = execute.EnsureKnownBranchesAncestry(execute.EnsureKnownBranchesAncestryArgs{
			BranchesToVerify: gitdomain.LocalBranchNames{branchesSnapshot.Active},
			Config:           validatedConfig,
			DefaultChoice:    validatedConfig.MainBranch,
			DialogTestInputs: &dialogTestInputs,
			LocalBranches:    branchesSnapshot.Branches,
			MainBranch:       validatedConfig.MainBranch,
			Runner:           runner,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSize, false, err
		}
		branchNamesToSync = gitdomain.LocalBranchNames{branchesSnapshot.Active}
		shouldPushTags = validatedConfig.IsMainOrPerennialBranch(branchesSnapshot.Active)
	}
	allBranchNamesToSync := validatedConfig.Lineage.BranchesAndAncestors(branchNamesToSync)
	branchesToSync, err := branchesSnapshot.Branches.Select(allBranchNamesToSync...)
	return &syncConfig{
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
