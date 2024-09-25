package cmd

import (
	"errors"
	"os"
	"slices"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/sync"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    cmdhelpers.Long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAppend(args[0], readDetachedFlag(cmd), readDryRunFlag(cmd), readPrototypeFlag(cmd), readVerboseFlag(cmd))
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
	runProgram := appendProgram(data)
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
	previousBranch            Option[gitdomain.LocalBranchName]
	prototype                 configdomain.Prototype
	remotes                   gitdomain.Remotes
	stashSize                 gitdomain.StashSize
	targetBranch              gitdomain.LocalBranchName
}

func determineAppendData(targetBranch gitdomain.LocalBranchName, repo execute.OpenRepoResult, detached configdomain.Detached, dryRun configdomain.DryRun, prototype configdomain.Prototype, verbose configdomain.Verbose) (data appendFeatureData, exit bool, err error) {
	fc := execute.FailureCollector{}
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
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
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
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchNamesToSync := validatedConfig.Config.Lineage.BranchAndAncestors(initialBranch)
	if detached {
		branchNamesToSync = validatedConfig.Config.RemovePerennials(branchNamesToSync)
	}
	branchesToSync, err := branchesToSync(branchNamesToSync, branchesSnapshot, repo, validatedConfig.Config.MainBranch)
	if err != nil {
		return data, false, err
	}
	initialAndAncestors := validatedConfig.Config.Lineage.BranchAndAncestors(initialBranch)
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
		previousBranch:            previousBranch,
		prototype:                 prototype,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	}, false, fc.Err
}

func appendProgram(data appendFeatureData) program.Program {
	prog := NewMutable(&program.Program{})
	if !data.hasOpenChanges {
		for _, branchToSync := range data.branchesToSync {
			sync.BranchProgram(branchToSync.BranchInfo, sync.BranchProgramArgs{
				BranchInfos:        data.branchInfos,
				Config:             data.config.Config,
				FirstCommitMessage: branchToSync.FirstCommitMessage,
				InitialBranch:      data.initialBranch,
				Program:            prog,
				Remotes:            data.remotes,
				PushBranches:       true,
			})
		}
	}
	prog.Value.Add(&opcodes.CreateAndCheckoutBranchExistingParent{
		Ancestors: data.newBranchParentCandidates,
		Branch:    data.targetBranch,
	})
	if data.remotes.HasOrigin() && data.config.Config.ShouldPushNewBranches() && data.config.Config.IsOnline() {
		prog.Value.Add(&opcodes.CreateTrackingBranch{Branch: data.targetBranch})
	}
	prog.Value.Add(&opcodes.SetExistingParent{
		Branch:    data.targetBranch,
		Ancestors: data.newBranchParentCandidates,
	})
	if data.prototype.IsTrue() || data.config.Config.CreatePrototypeBranches.IsTrue() {
		prog.Value.Add(&opcodes.AddToPrototypeBranches{Branch: data.targetBranch})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.initialBranch), data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Get()
}
