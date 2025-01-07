package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v17/internal/cli/dialog"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/cmd/ship"
	"github.com/git-town/git-town/v17/internal/cmd/sync"
	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/execute"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v17/internal/hosting"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/undo/undoconfig"
	"github.com/git-town/git-town/v17/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v17/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v17/internal/vm/opcodes"
	"github.com/git-town/git-town/v17/internal/vm/optimizer"
	"github.com/git-town/git-town/v17/internal/vm/program"
	"github.com/git-town/git-town/v17/internal/vm/runstate"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/git-town/git-town/v17/pkg/set"
	"github.com/spf13/cobra"
)

const prependDesc = "Create a new feature branch as the parent of the current branch"

const prependHelp = `
Syncs the parent branch, cuts a new feature branch with the given name off the parent branch, makes the new branch the parent of the current branch, pushes the new feature branch to the origin repository (if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.`

func prependCommand() *cobra.Command {
	addBeamFlag, readBeamFlag := flags.Beam()
	addBodyFlag, readBodyFlag := flags.ProposalBody("")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addTitleFlag, readTitleFlag := flags.ProposalTitle()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "prepend <branch>",
		GroupID: "stack",
		Args:    cobra.ExactArgs(1),
		Short:   prependDesc,
		Long:    cmdhelpers.Long(prependDesc, prependHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			beam, err := readBeamFlag(cmd)
			if err != nil {
				return err
			}
			bodyText, err := readBodyFlag(cmd)
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
			title, err := readTitleFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executePrepend(args, beam, bodyText, detached, dryRun, propose, prototype, title, verbose)
		},
	}
	addBeamFlag(&cmd)
	addBodyFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addTitleFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executePrepend(args []string, beam configdomain.Beam, proposalBody gitdomain.ProposalBody, detached configdomain.Detached, dryRun configdomain.DryRun, propose configdomain.Propose, prototype configdomain.Prototype, proposalTitle gitdomain.ProposalTitle, verbose configdomain.Verbose) error {
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
	data, exit, err := determinePrependData(args, repo, beam, detached, dryRun, proposalBody, proposalTitle, propose, prototype, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := prependProgram(data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "prepend",
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
		Connector:               data.connector,
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

type prependData struct {
	branchInfos         gitdomain.BranchInfos
	branchesSnapshot    gitdomain.BranchesSnapshot
	branchesToSync      configdomain.BranchesToSync
	commitsToBeam       gitdomain.Commits
	config              config.ValidatedConfig
	connector           Option[hostingdomain.Connector]
	dialogTestInputs    components.TestInputs
	dryRun              configdomain.DryRun
	existingParent      gitdomain.LocalBranchName
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	newParentCandidates gitdomain.LocalBranchNames
	nonExistingBranches gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	preFetchBranchInfos gitdomain.BranchInfos
	previousBranch      Option[gitdomain.LocalBranchName]
	proposal            Option[hostingdomain.Proposal]
	proposalBody        gitdomain.ProposalBody
	proposalTitle       gitdomain.ProposalTitle
	propose             configdomain.Propose
	prototype           configdomain.Prototype
	remotes             gitdomain.Remotes
	stashSize           gitdomain.StashSize
	targetBranch        gitdomain.LocalBranchName
}

func determinePrependData(args []string, repo execute.OpenRepoResult, beam configdomain.Beam, detached configdomain.Detached, dryRun configdomain.DryRun, propasalBody gitdomain.ProposalBody, proposalTitle gitdomain.ProposalTitle, propose configdomain.Propose, prototype configdomain.Prototype, verbose configdomain.Verbose) (data prependData, exit bool, err error) {
	prefetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	fc := execute.FailureCollector{}
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
	targetBranch := gitdomain.NewLocalBranchName(args[0])
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch, repo.UnvalidatedConfig.NormalConfig.DevRemote) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	ancestorOpt := latestExistingAncestor(initialBranch, branchesSnapshot.Branches, validatedConfig.NormalConfig.Lineage)
	ancestor, hasAncestor := ancestorOpt.Get()
	if !hasAncestor {
		return data, false, fmt.Errorf(messages.SetParentNoFeatureBranch, branchesSnapshot.Active)
	}
	commitsToBeam := []gitdomain.Commit{}
	if beam {
		commitsInBranch, err := repo.Git.CommitsInFeatureBranch(repo.Backend, initialBranch, ancestor)
		if err != nil {
			return data, false, err
		}
		commitsToBeam, exit, err = dialog.CommitsToBeam(commitsInBranch, ancestor, dialogTestInputs.Next())
		if err != nil || exit {
			return data, exit, err
		}
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
	parentAndAncestors := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(ancestor)
	slices.Reverse(parentAndAncestors)
	proposalOpt := ship.FindProposal(connector, initialBranch, Some(ancestor))
	return prependData{
		branchInfos:         branchesSnapshot.Branches,
		branchesSnapshot:    branchesSnapshot,
		branchesToSync:      branchesToSync,
		commitsToBeam:       commitsToBeam,
		config:              validatedConfig,
		connector:           connector,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		existingParent:      ancestor,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		newParentCandidates: parentAndAncestors,
		nonExistingBranches: nonExistingBranches,
		preFetchBranchInfos: prefetchBranchSnapshot.Branches,
		previousBranch:      previousBranch,
		proposal:            proposalOpt,
		proposalBody:        propasalBody,
		proposalTitle:       proposalTitle,
		propose:             propose,
		prototype:           prototype,
		remotes:             remotes,
		stashSize:           stashSize,
		targetBranch:        targetBranch,
	}, false, fc.Err
}

func prependProgram(data prependData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	if !data.hasOpenChanges {
		data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, finalMessages)
		branchesToDelete := set.New[gitdomain.LocalBranchName]()
		sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
			BranchInfos:         data.branchInfos,
			BranchesToDelete:    NewMutable(&branchesToDelete),
			Config:              data.config,
			InitialBranch:       data.initialBranch,
			PrefetchBranchInfos: data.preFetchBranchInfos,
			Program:             prog,
			PushBranches:        true,
			Remotes:             data.remotes,
		})
	}
	prog.Value.Add(&opcodes.BranchCreateAndCheckoutExistingParent{
		Ancestors: data.newParentCandidates,
		Branch:    data.targetBranch,
	})
	// set the parent of the newly created branch
	prog.Value.Add(&opcodes.LineageParentSetFirstExisting{
		Branch:    data.targetBranch,
		Ancestors: data.newParentCandidates,
	})
	// set the parent of the branch prepended to
	prog.Value.Add(&opcodes.LineageParentSetIfExists{
		Branch: data.initialBranch,
		Parent: data.targetBranch,
	})
	if data.prototype.IsTrue() {
		prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: configdomain.BranchTypePrototypeBranch})
	} else {
		if newBranchType, hasNewBranchType := data.config.NormalConfig.NewBranchType.Get(); hasNewBranchType {
			switch newBranchType {
			case
				configdomain.BranchTypePrototypeBranch,
				configdomain.BranchTypeContributionBranch,
				configdomain.BranchTypeObservedBranch,
				configdomain.BranchTypeParkedBranch,
				configdomain.BranchTypePerennialBranch:
				prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: newBranchType})
			case configdomain.BranchTypeFeatureBranch:
			case configdomain.BranchTypeMainBranch:
			}
		}
	}
	proposal, hasProposal := data.proposal.Get()
	if data.remotes.HasRemote(data.config.NormalConfig.DevRemote) && data.config.NormalConfig.IsOnline() && (data.config.NormalConfig.ShouldPushNewBranches() || hasProposal) {
		prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.targetBranch})
	}
	connector, hasConnector := data.connector.Get()
	connectorCanUpdateProposalTargets := hasConnector && connector.UpdateProposalTargetFn().IsSome()
	if hasProposal && hasConnector && connectorCanUpdateProposalTargets {
		prog.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch:      data.targetBranch,
			OldBranch:      data.existingParent,
			ProposalNumber: proposal.Number,
		})
	}
	moveCommitsToNewBranch(prog, data)
	if data.propose {
		prog.Value.Add(
			&opcodes.PushCurrentBranchIfLocal{
				CurrentBranch: data.targetBranch,
			},
			&opcodes.ProposalCreate{
				Branch:        data.targetBranch,
				MainBranch:    data.config.ValidatedConfigData.MainBranch,
				ProposalBody:  data.proposalBody,
				ProposalTitle: data.proposalTitle,
			})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

// provides the strategy to use to sync a branch after beaming some of its commits to its new parent branch
func afterBeamToParentSyncStrategy(branchType configdomain.BranchType, config configdomain.NormalConfigData) Option[configdomain.SyncStrategy] {
	switch branchType {
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return None[configdomain.SyncStrategy]()
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch:
		return Some(config.SyncFeatureStrategy.SyncStrategy())
	case configdomain.BranchTypePrototypeBranch:
		return Some(config.SyncPrototypeStrategy.SyncStrategy())
	}
	panic("unhandled branch type: " + branchType.String())
}

// provides the name of the youngest ancestor branch of the given branch that actually exists,
// either locally or remotely.
func latestExistingAncestor(branch gitdomain.LocalBranchName, branchInfos gitdomain.BranchInfos, lineage configdomain.Lineage) Option[gitdomain.LocalBranchName] {
	for {
		parent, hasParent := lineage.Parent(branch).Get()
		if !hasParent {
			return None[gitdomain.LocalBranchName]()
		}
		if branchInfos.HasBranch(parent) {
			return Some(parent)
		}
		branch = parent
	}
}

func moveCommitsToNewBranch(prog Mutable[program.Program], data prependData) {
	if len(data.commitsToBeam) > 0 {
		for _, commitToBeam := range data.commitsToBeam {
			prog.Value.Add(
				&opcodes.CherryPick{SHA: commitToBeam.SHA},
			)
		}
		// sync the initial branch with the new parent branch to remove the moved commits from the initial branch
		prog.Value.Add(
			&opcodes.Checkout{Branch: data.initialBranch},
		)
		initialBranchType := data.config.BranchType(data.initialBranch)
		syncWithParent(prog, data.targetBranch, initialBranchType, data.config.NormalConfig.NormalConfigData)
		prog.Value.Add(
			&opcodes.Checkout{Branch: data.targetBranch},
		)
	}
}

// basic sync of the current branch with its parent after beaming some commits into the parent
func syncWithParent(prog Mutable[program.Program], parentName gitdomain.LocalBranchName, initialBranchType configdomain.BranchType, config configdomain.NormalConfigData) {
	if syncStrategy, hasSyncStrategy := afterBeamToParentSyncStrategy(initialBranchType, config).Get(); hasSyncStrategy {
		switch syncStrategy {
		case configdomain.SyncStrategyCompress, configdomain.SyncStrategyMerge:
			prog.Value.Add(
				&opcodes.MergeParent{Parent: parentName.BranchName()},
				&opcodes.PushCurrentBranch{},
			)
		case configdomain.SyncStrategyRebase:
			prog.Value.Add(
				&opcodes.RebaseBranch{Branch: parentName.BranchName()},
				&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
			)
		case configdomain.SyncStrategyFFOnly:
			// the ff-only sync strategy does not sync with the parent
		}
	}
}
