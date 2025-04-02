package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v18/internal/cli/dialog"
	"github.com/git-town/git-town/v18/internal/cli/dialog/components"
	"github.com/git-town/git-town/v18/internal/cli/flags"
	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v18/internal/cmd/ship"
	"github.com/git-town/git-town/v18/internal/config"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/execute"
	"github.com/git-town/git-town/v18/internal/forge"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	"github.com/git-town/git-town/v18/internal/undo/undoconfig"
	"github.com/git-town/git-town/v18/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v18/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v18/internal/vm/opcodes"
	"github.com/git-town/git-town/v18/internal/vm/optimizer"
	"github.com/git-town/git-town/v18/internal/vm/program"
	"github.com/git-town/git-town/v18/internal/vm/runstate"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	setParentCmd  = "set-parent"
	setParentDesc = "Set the parent branch for the current branch"
	setParentHelp = `
Consider this branch stack:

main
 \
  feature-1
   \
*   feature-B
 \
  feature-A

After running "git town set-parent"
and selecting "feature-A" in the dialog,
we end up with this branch stack:

main
 \
  feature-1
 \
  feature-A
   \
*   feature-B
`
)

func setParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     setParentCmd,
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.NoArgs,
		Short:   setParentDesc,
		Long:    cmdhelpers.Long(setParentDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeSetParent(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSetParent(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSetParentData(repo, verbose)
	if err != nil || exit {
		return err
	}
	err = verifySetParentData(data)
	if err != nil {
		return err
	}
	outcome, selectedBranch, err := dialog.Parent(dialog.ParentArgs{
		Branch:          data.initialBranch,
		DefaultChoice:   data.defaultChoice,
		DialogTestInput: data.dialogTestInputs.Next(),
		Lineage:         data.config.NormalConfig.Lineage,
		LocalBranches:   data.branchesSnapshot.Branches.LocalBranches().Names(),
		MainBranch:      data.mainBranch,
	})
	if err != nil {
		return err
	}
	runProgram, aborted := setParentProgram(outcome, selectedBranch, data)
	if aborted {
		return nil
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               setParentCmd,
		DryRun:                false,
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

type setParentData struct {
	branchesSnapshot gitdomain.BranchesSnapshot
	config           config.ValidatedConfig
	connector        Option[forgedomain.Connector]
	defaultChoice    gitdomain.LocalBranchName
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	mainBranch       gitdomain.LocalBranchName
	proposal         Option[forgedomain.Proposal]
	stashSize        gitdomain.StashSize
}

func determineSetParentData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data setParentData, exit bool, err error) {
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
		Fetch:                 false,
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
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	connectorOpt, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		Connector:          connectorOpt,
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
	mainBranch := validatedConfig.ValidatedConfigData.MainBranch
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	parentOpt := validatedConfig.NormalConfig.Lineage.Parent(initialBranch)
	existingParent, hasParent := parentOpt.Get()
	var defaultChoice gitdomain.LocalBranchName
	if hasParent {
		defaultChoice = existingParent
	} else {
		defaultChoice = mainBranch
	}
	proposalOpt := None[forgedomain.Proposal]()
	if !repo.IsOffline {
		proposalOpt = ship.FindProposal(connectorOpt, initialBranch, parentOpt)
	}
	return setParentData{
		branchesSnapshot: branchesSnapshot,
		config:           validatedConfig,
		connector:        connectorOpt,
		defaultChoice:    defaultChoice,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    initialBranch,
		mainBranch:       mainBranch,
		proposal:         proposalOpt,
		stashSize:        stashSize,
	}, false, nil
}

func verifySetParentData(data setParentData) error {
	if data.config.IsMainOrPerennialBranch(data.initialBranch) {
		return fmt.Errorf(messages.SetParentNoFeatureBranch, data.initialBranch)
	}
	return nil
}

func setParentProgram(dialogOutcome dialog.ParentOutcome, selectedBranch gitdomain.LocalBranchName, data setParentData) (prog program.Program, aborted bool) {
	proposal, hasProposal := data.proposal.Get()
	// update lineage
	switch dialogOutcome {
	case dialog.ParentOutcomeAborted:
		return prog, true
	case dialog.ParentOutcomePerennialBranch:
		prog.Add(&opcodes.BranchTypeOverrideSet{Branch: data.initialBranch, BranchType: configdomain.BranchTypePerennialBranch})
		prog.Add(&opcodes.LineageParentRemove{Branch: data.initialBranch})
	case dialog.ParentOutcomeSelectedParent:
		prog.Add(&opcodes.LineageParentSet{Branch: data.initialBranch, Parent: selectedBranch})
		connector, hasConnector := data.connector.Get()
		connectorCanUpdateProposalTarget := hasConnector && connector.UpdateProposalTargetFn().IsSome()
		if hasProposal && hasConnector && connectorCanUpdateProposalTarget {
			prog.Add(&opcodes.ProposalUpdateTarget{
				NewBranch:      selectedBranch,
				OldBranch:      proposal.Target,
				ProposalNumber: proposal.Number,
			})
		}
		// update commits
		switch data.config.NormalConfig.SyncFeatureStrategy {
		case configdomain.SyncFeatureStrategyMerge:
		case configdomain.SyncFeatureStrategyCompress, configdomain.SyncFeatureStrategyRebase:
			initialBranchInfo, hasInitialBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(data.initialBranch).Get()
			hasRemoteBranch := hasInitialBranchInfo && initialBranchInfo.HasTrackingBranch()
			if hasRemoteBranch {
				prog.Add(
					&opcodes.PullCurrentBranch{},
				)
			}
			parentOpt := data.config.NormalConfig.Lineage.Parent(data.initialBranch)
			prog.Add(
				&opcodes.RebaseOntoRemoveDeleted{
					BranchToRebaseOnto: selectedBranch,
					CommitsToRemove:    data.initialBranch.BranchName(),
					Upstream:           parentOpt,
				},
			)
			if hasRemoteBranch {
				prog.Add(
					&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
				)
			}
			// remove commits from descendents
			descendents := data.config.NormalConfig.Lineage.Descendants(data.initialBranch)
			for _, descendent := range descendents {
				prog.Add(
					&opcodes.CheckoutIfNeeded{
						Branch: descendent,
					},
				)
				descendentBranchInfo, hasDescendentBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get()
				if hasDescendentBranchInfo && descendentBranchInfo.HasTrackingBranch() {
					prog.Add(
						&opcodes.PullCurrentBranch{},
					)
				}
				prog.Add(
					&opcodes.RebaseOntoRemoveDeleted{
						BranchToRebaseOnto: data.initialBranch,
						CommitsToRemove:    descendent.BranchName(),
						Upstream:           parentOpt,
					},
				)
				if hasDescendentBranchInfo && descendentBranchInfo.HasTrackingBranch() {
					prog.Add(
						&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
					)
				}
			}
			prog.Add(
				&opcodes.CheckoutIfNeeded{
					Branch: data.initialBranch,
				},
			)
		}
	}
	return optimizer.Optimize(prog), false
}
