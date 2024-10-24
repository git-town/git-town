package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/cmd/ship"
	"github.com/git-town/git-town/v16/internal/cmd/sync"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
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
	"github.com/spf13/cobra"
)

const prependDesc = "Create a new feature branch as the parent of the current branch"

const prependHelp = `
Syncs the parent branch, cuts a new feature branch with the given name off the parent branch, makes the new branch the parent of the current branch, pushes the new feature branch to the origin repository (if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.`

func prependCommand() *cobra.Command {
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "prepend <branch>",
		GroupID: "lineage",
		Args:    cobra.ExactArgs(1),
		Short:   prependDesc,
		Long:    cmdhelpers.Long(prependDesc, prependHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePrepend(args, readDetachedFlag(cmd), readDryRunFlag(cmd), readPrototypeFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addPrototypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executePrepend(args []string, detached configdomain.Detached, dryRun configdomain.DryRun, prototype configdomain.Prototype, verbose configdomain.Verbose) error {
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
	data, exit, err := determinePrependData(args, repo, detached, dryRun, prototype, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := prependProgram(data)
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
	branchesToSync      []configdomain.BranchToSync
	config              config.ValidatedConfig
	connector           Option[hostingdomain.Connector]
	dialogTestInputs    components.TestInputs
	dryRun              configdomain.DryRun
	existingParent      gitdomain.LocalBranchName
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	newParentCandidates gitdomain.LocalBranchNames
	preFetchBranchInfos gitdomain.BranchInfos
	previousBranch      Option[gitdomain.LocalBranchName]
	proposal            Option[hostingdomain.Proposal]
	prototype           configdomain.Prototype
	remotes             gitdomain.Remotes
	stashSize           gitdomain.StashSize
	targetBranch        gitdomain.LocalBranchName
}

func determinePrependData(args []string, repo execute.OpenRepoResult, detached configdomain.Detached, dryRun configdomain.DryRun, prototype configdomain.Prototype, verbose configdomain.Verbose) (data prependData, exit bool, err error) {
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
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, gitdomain.RemoteOrigin, print.Logger{})
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
	branchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	if detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchesToSync, err := sync.BranchesToSync(branchNamesToSync, branchesSnapshot, repo, validatedConfig.ValidatedConfig.MainBranch)
	if err != nil {
		return data, false, err
	}
	parentOpt := validatedConfig.NormalConfig.Lineage.Parent(initialBranch)
	parent, hasParent := parentOpt.Get()
	if !hasParent {
		return data, false, fmt.Errorf(messages.SetParentNoFeatureBranch, branchesSnapshot.Active)
	}
	parentAndAncestors := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(parent)
	slices.Reverse(parentAndAncestors)
	proposalOpt := ship.FindProposal(connector, initialBranch, parentOpt)
	return prependData{
		branchInfos:         branchesSnapshot.Branches,
		branchesSnapshot:    branchesSnapshot,
		branchesToSync:      branchesToSync,
		config:              validatedConfig,
		connector:           connector,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		existingParent:      parent,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		newParentCandidates: parentAndAncestors,
		preFetchBranchInfos: prefetchBranchSnapshot.Branches,
		previousBranch:      previousBranch,
		proposal:            proposalOpt,
		prototype:           prototype,
		remotes:             remotes,
		stashSize:           stashSize,
		targetBranch:        targetBranch,
	}, false, fc.Err
}

func prependProgram(data prependData) program.Program {
	prog := NewMutable(&program.Program{})
	if !data.hasOpenChanges {
		sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
			BranchInfos:         data.branchInfos,
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
	if data.prototype.IsTrue() || data.config.NormalConfig.CreatePrototypeBranches.IsTrue() {
		prog.Value.Add(&opcodes.BranchesPrototypeAdd{Branch: data.targetBranch})
	}
	proposal, hasProposal := data.proposal.Get()
	if data.remotes.HasOrigin() && data.config.NormalConfig.IsOnline() && (data.config.NormalConfig.ShouldPushNewBranches() || hasProposal) {
		prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.targetBranch})
	}
	if hasProposal {
		prog.Value.Add(&opcodes.ProposalUpdateBase{
			NewTarget:      data.targetBranch,
			OldTarget:      data.existingParent,
			ProposalNumber: proposal.Number,
		})
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Get()
}
