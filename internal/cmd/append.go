package cmd

import (
	"errors"
	"os"
	"slices"

	"github.com/git-town/git-town/v20/internal/cli/dialog"
	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/flags"
	"github.com/git-town/git-town/v20/internal/cli/print"
	"github.com/git-town/git-town/v20/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v20/internal/cmd/sync"
	"github.com/git-town/git-town/v20/internal/config"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/execute"
	"github.com/git-town/git-town/v20/internal/forge"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/undo/undoconfig"
	"github.com/git-town/git-town/v20/internal/validate"
	"github.com/git-town/git-town/v20/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
	"github.com/git-town/git-town/v20/internal/vm/optimizer"
	"github.com/git-town/git-town/v20/internal/vm/program"
	"github.com/git-town/git-town/v20/internal/vm/runstate"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/git-town/git-town/v20/pkg/set"
	"github.com/spf13/cobra"
)

const (
	appendDesc = "Create a new feature branch as a child of the current branch"
	appendHelp = `
Consider this stack:

main
 \
* feature-1

We are on the "feature-1" branch,
which is a child of branch "main".
After running "git town append feature-2",
the repository will have these branches:

main
 \
  feature-1
   \
*   feature-2

The new branch "feature-2"
is a child of "feature-1".

If there are no uncommitted changes,
it also syncs all affected branches.
`
)

func appendCmd() *cobra.Command {
	addBeamFlag, readBeamFlag := flags.Beam()
	addCommitFlag, readCommitFlag := flags.Commit()
	addCommitMessageFlag, readCommitMessageFlag := flags.CommitMessage("the commit message")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "append <branch>",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    cmdhelpers.Long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			beam, err := readBeamFlag(cmd)
			if err != nil {
				return err
			}
			commit, err := readCommitFlag(cmd)
			if err != nil {
				return err
			}
			commitMessage, err := readCommitMessageFlag(cmd)
			if err != nil {
				return err
			}
			detached, err := readDetachedFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			propose, err := readProposeFlag(cmd)
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
			if commitMessage.IsSome() || propose.IsTrue() {
				commit = true
			}
			return executeAppend(args[0], beam, commit, commitMessage, detached, dryRun, propose, prototype, verbose)
		},
	}
	addBeamFlag(&cmd)
	addCommitFlag(&cmd)
	addCommitMessageFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeAppend(arg string, beam configdomain.Beam, commit configdomain.Commit, commitMessage Option[gitdomain.CommitMessage], detached configdomain.Detached, dryRun configdomain.DryRun, propose configdomain.Propose, prototype configdomain.Prototype, verbose configdomain.Verbose) error {
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
	data, exit, err := determineAppendData(gitdomain.NewLocalBranchName(arg), beam, repo, commit, commitMessage, detached, dryRun, propose, prototype, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := appendProgram(data, repo.FinalMessages, false)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "append",
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		BranchInfosLastRun:    data.branchInfosLastRun,
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		Detached:                detached,
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
	beam                      configdomain.Beam
	branchInfos               gitdomain.BranchInfos
	branchInfosLastRun        Option[gitdomain.BranchInfos]
	branchesSnapshot          gitdomain.BranchesSnapshot
	branchesToSync            configdomain.BranchesToSync
	commit                    configdomain.Commit
	commitMessage             Option[gitdomain.CommitMessage]
	commitsToBeam             gitdomain.Commits
	config                    config.ValidatedConfig
	connector                 Option[forgedomain.Connector]
	detached                  configdomain.Detached
	dialogTestInputs          components.TestInputs
	dryRun                    configdomain.DryRun
	hasOpenChanges            bool
	initialBranch             gitdomain.LocalBranchName
	initialBranchInfo         *gitdomain.BranchInfo
	newBranchParentCandidates gitdomain.LocalBranchNames
	nonExistingBranches       gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	preFetchBranchInfos       gitdomain.BranchInfos
	previousBranch            Option[gitdomain.LocalBranchName]
	propose                   configdomain.Propose
	prototype                 configdomain.Prototype
	remotes                   gitdomain.Remotes
	stashSize                 gitdomain.StashSize
	targetBranch              gitdomain.LocalBranchName
}

func determineAppendData(targetBranch gitdomain.LocalBranchName, beam configdomain.Beam, repo execute.OpenRepoResult, commit configdomain.Commit, commitMessage Option[gitdomain.CommitMessage], detached configdomain.Detached, dryRun configdomain.DryRun, propose configdomain.Propose, prototype configdomain.Prototype, verbose configdomain.Verbose) (data appendFeatureData, exit bool, err error) {
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
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Detached:              detached,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 !repoStatus.OpenChanges && beam.IsFalse() && commit.IsFalse(),
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
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch, repo.UnvalidatedConfig.NormalConfig.DevRemote) {
		fc.Fail(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	connector, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
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
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	initialAndAncestors := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	slices.Reverse(initialAndAncestors)
	commitsToBeam := []gitdomain.Commit{}
	ancestor, hasAncestor := latestExistingAncestor(initialBranch, branchesSnapshot.Branches, validatedConfig.NormalConfig.Lineage).Get()
	if beam.IsTrue() && hasAncestor {
		commitsInBranch, err := repo.Git.CommitsInFeatureBranch(repo.Backend, initialBranch, ancestor.BranchName())
		if err != nil {
			return data, false, err
		}
		commitsToBeam, exit, err = dialog.CommitsToBeam(commitsInBranch, targetBranch, repo.Git, repo.Backend, dialogTestInputs.Next())
		if err != nil || exit {
			return data, exit, err
		}
	}
	if validatedConfig.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPropose {
		propose = true
	}
	return appendFeatureData{
		beam:                      beam,
		branchInfos:               branchesSnapshot.Branches,
		branchInfosLastRun:        branchInfosLastRun,
		branchesSnapshot:          branchesSnapshot,
		branchesToSync:            branchesToSync,
		commit:                    commit,
		commitMessage:             commitMessage,
		commitsToBeam:             commitsToBeam,
		config:                    validatedConfig,
		connector:                 connector,
		detached:                  detached,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		initialBranchInfo:         initialBranchInfo,
		newBranchParentCandidates: initialAndAncestors,
		nonExistingBranches:       nonExistingBranches,
		preFetchBranchInfos:       preFetchBranchSnapshot.Branches,
		previousBranch:            previousBranch,
		propose:                   propose,
		prototype:                 prototype,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	}, false, fc.Err
}

func appendProgram(data appendFeatureData, finalMessages stringslice.Collector, beamCherryPick bool) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, finalMessages)
	if !data.hasOpenChanges && data.beam.IsFalse() && data.commit.IsFalse() {
		branchesToDelete := set.New[gitdomain.LocalBranchName]()
		sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
			BranchInfos:         data.branchInfos,
			BranchInfosLastRun:  data.branchInfosLastRun,
			BranchesToDelete:    NewMutable(&branchesToDelete),
			Config:              data.config,
			InitialBranch:       data.initialBranch,
			PrefetchBranchInfos: data.preFetchBranchInfos,
			Program:             prog,
			Prune:               false,
			Remotes:             data.remotes,
			PushBranches:        true,
		})
	}
	prog.Value.Add(&opcodes.BranchCreateAndCheckoutExistingParent{
		Ancestors: data.newBranchParentCandidates,
		Branch:    data.targetBranch,
	})
	if data.remotes.HasRemote(data.config.NormalConfig.DevRemote) && data.config.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPush && data.config.NormalConfig.Offline.IsOnline() {
		prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.targetBranch})
	}
	prog.Value.Add(&opcodes.LineageParentSetFirstExisting{
		Branch:    data.targetBranch,
		Ancestors: data.newBranchParentCandidates,
	})
	if data.prototype {
		prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: configdomain.BranchTypePrototypeBranch})
	} else {
		if newBranchType, hasNewBranchType := data.config.NormalConfig.NewBranchType.Get(); hasNewBranchType {
			switch newBranchType {
			case
				configdomain.BranchTypeContributionBranch,
				configdomain.BranchTypeObservedBranch,
				configdomain.BranchTypeParkedBranch,
				configdomain.BranchTypePerennialBranch,
				configdomain.BranchTypePrototypeBranch,
				configdomain.BranchTypeFeatureBranch:
				prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: newBranchType})
			case configdomain.BranchTypeMainBranch:
			}
		}
	}
	if data.commit {
		prog.Value.Add(
			&opcodes.Commit{
				AuthorOverride:                 None[gitdomain.Author](),
				FallbackToDefaultCommitMessage: false,
				Message:                        data.commitMessage,
			},
		)
	}
	moveCommitsToAppendedBranch(prog, data, beamCherryPick)
	if data.propose {
		prog.Value.Add(
			&opcodes.BranchTrackingCreate{
				Branch: data.targetBranch,
			},
			&opcodes.ProposalCreate{
				Branch:        data.targetBranch,
				MainBranch:    data.config.ValidatedConfigData.MainBranch,
				ProposalBody:  "",
				ProposalTitle: gitdomain.ProposalTitle(data.commitMessage.GetOrDefault()),
			},
		)
	}
	if data.commit {
		prog.Value.Add(
			&opcodes.Checkout{Branch: data.initialBranch},
		)
	} else {
		previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.initialBranch), data.previousBranch}
		cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
			DryRun:                   data.dryRun,
			InitialStashSize:         data.stashSize,
			RunInGitRoot:             true,
			StashOpenChanges:         data.hasOpenChanges,
			PreviousBranchCandidates: previousBranchCandidates,
		})
	}
	return optimizer.Optimize(prog.Immutable())
}

func moveCommitsToAppendedBranch(prog Mutable[program.Program], data appendFeatureData, performCherryPick bool) {
	if len(data.commitsToBeam) == 0 {
		return
	}
	if performCherryPick {
		for _, commitToBeam := range data.commitsToBeam {
			prog.Value.Add(
				&opcodes.CherryPick{SHA: commitToBeam.SHA},
			)
		}
	}
	prog.Value.Add(
		&opcodes.Checkout{
			Branch: data.initialBranch,
		},
	)
	for c := len(data.commitsToBeam) - 1; c >= 0; c-- {
		commitToBeam := data.commitsToBeam[c]
		prog.Value.Add(
			&opcodes.CommitRemove{
				SHA: commitToBeam.SHA,
			},
		)
	}
	if data.initialBranchInfo.HasTrackingBranch() {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIgnoreError{},
		)
	}
	prog.Value.Add(
		&opcodes.Checkout{
			Branch: data.targetBranch,
		},
	)
	if !performCherryPick {
		prog.Value.Add(
			&opcodes.RebaseBranch{
				Branch: data.initialBranch.BranchName(),
			},
		)
	}
}
