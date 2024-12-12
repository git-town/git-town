package cmd

import (
	"errors"
	"os"
	"slices"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/cmd/sync"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/git-town/git-town/v16/pkg/set"
	"github.com/spf13/cobra"
)

const appendDesc = "Create a new feature branch as a child of the current branch"

const appendHelp = `
Syncs the current branch, forks a new feature branch with the given name off the current branch, makes the new branch a child of the current branch, pushes the new feature branch to the origin repository (if and only if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func appendCmd() *cobra.Command {
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "append <branch>",
		GroupID: "stack",
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    cmdhelpers.Long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			detached, err := readDetachedFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			prototype, err := readPrototypeFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeAppend(args[0], detached, dryRun, prototype, verbose)
		},
	}
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addPrototypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeAppend(arg string, detached configdomain.Detached, dryRun configdomain.DryRun, prototype configdomain.Prototype, verbose configdomain.Verbose) error {
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
	data, exit, err := determineAppendData(gitdomain.NewLocalBranchName(arg), repo, detached, dryRun, prototype, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := appendProgram(data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "append",
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
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

type appendFeatureData struct {
	branchInfos               gitdomain.BranchInfos
	branchesSnapshot          gitdomain.BranchesSnapshot
	branchesToSync            []configdomain.BranchToSync
	config                    config.ValidatedConfig
	dialogTestInputs          components.TestInputs
	dryRun                    configdomain.DryRun
	hasOpenChanges            bool
	initialBranch             gitdomain.LocalBranchName
	newBranchParentCandidates gitdomain.LocalBranchNames
	nonExistingBranches       gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	preFetchBranchInfos       gitdomain.BranchInfos
	previousBranch            Option[gitdomain.LocalBranchName]
	prototype                 configdomain.Prototype
	remotes                   gitdomain.Remotes
	stashSize                 gitdomain.StashSize
	targetBranch              gitdomain.LocalBranchName
}

func determineAppendData(targetBranch gitdomain.LocalBranchName, repo execute.OpenRepoResult, detached configdomain.Detached, dryRun configdomain.DryRun, prototype configdomain.Prototype, verbose configdomain.Verbose) (data appendFeatureData, exit bool, err error) {
	fc := execute.FailureCollector{}
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
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
		Fetch:                 !repoStatus.OpenChanges,
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
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes := fc.Remotes(repo.Git.Remotes(repo.Backend))
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, gitdomain.RemoteOrigin, print.Logger{})
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      branchesSnapshot.Branches.LocalBranches().Names(),
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	if detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	initialAndAncestors := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	slices.Reverse(initialAndAncestors)
	return appendFeatureData{
		branchInfos:               branchesSnapshot.Branches,
		branchesSnapshot:          branchesSnapshot,
		branchesToSync:            branchesToSync,
		config:                    validatedConfig,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		newBranchParentCandidates: initialAndAncestors,
		nonExistingBranches:       nonExistingBranches,
		preFetchBranchInfos:       preFetchBranchSnapshot.Branches,
		previousBranch:            previousBranch,
		prototype:                 prototype,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	}, false, fc.Err
}

func appendProgram(data appendFeatureData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, finalMessages)
	if !data.hasOpenChanges {
		branchesToDelete := set.New[gitdomain.LocalBranchName]()
		sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
			BranchInfos:         data.branchInfos,
			BranchesToDelete:    NewMutable(&branchesToDelete),
			Config:              data.config,
			InitialBranch:       data.initialBranch,
			PrefetchBranchInfos: data.preFetchBranchInfos,
			Program:             prog,
			Remotes:             data.remotes,
			PushBranches:        true,
		})
	}
	prog.Value.Add(&opcodes.BranchCreateAndCheckoutExistingParent{
		Ancestors: data.newBranchParentCandidates,
		Branch:    data.targetBranch,
	})
	if data.remotes.HasOrigin() && data.config.NormalConfig.ShouldPushNewBranches() && data.config.NormalConfig.IsOnline() {
		prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.targetBranch})
	}
	prog.Value.Add(&opcodes.LineageParentSetFirstExisting{
		Branch:    data.targetBranch,
		Ancestors: data.newBranchParentCandidates,
	})
	if data.prototype.IsTrue() {
		prog.Value.Add(&opcodes.BranchesPrototypeAdd{Branch: data.targetBranch})
	} else {
		switch data.config.NormalConfig.NewBranchType {
		case configdomain.BranchTypeContributionBranch:
			prog.Value.Add(&opcodes.BranchesContributionAdd{Branch: data.targetBranch})
		case configdomain.BranchTypeFeatureBranch:
		case configdomain.BranchTypeMainBranch:
		case configdomain.BranchTypeObservedBranch:
			prog.Value.Add(&opcodes.BranchesObservedAdd{Branch: data.targetBranch})
		case configdomain.BranchTypeParkedBranch:
			prog.Value.Add(&opcodes.BranchesParkedAdd{Branch: data.targetBranch})
		case configdomain.BranchTypePerennialBranch:
			prog.Value.Add(&opcodes.BranchesPerennialAdd{Branch: data.targetBranch})
		case configdomain.BranchTypePrototypeBranch:
			prog.Value.Add(&opcodes.BranchesPrototypeAdd{Branch: data.targetBranch})
		}
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.initialBranch), data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Immutable()
}
