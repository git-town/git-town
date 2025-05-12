package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/flags"
	"github.com/git-town/git-town/v20/internal/cli/print"
	"github.com/git-town/git-town/v20/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v20/internal/config"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/execute"
	"github.com/git-town/git-town/v20/internal/forge"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/undo/undoconfig"
	"github.com/git-town/git-town/v20/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v20/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v20/internal/vm/optimizer"
	"github.com/git-town/git-town/v20/internal/vm/program"
	"github.com/git-town/git-town/v20/internal/vm/runstate"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	forEachCmd  = "for-each"
	forEachDesc = "Executes the given shell command on each branch"
	forEachHelp = `
Executes the given shell command on each branch.

Consider this stack:

main
 \
  branch-1
   \
    branch-2

When running "git town for-each --stack echo hello",
it prints this output.

[main] hello

[branch1] hello

[branch2] hello
`
)

func forEachCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     forEachCmd,
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   forEachDesc,
		Long:    cmdhelpers.Long(forEachDesc, forEachHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeForEach(dryRun, verbose)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeForEach(dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
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
	data, exit, err := determineForEachData(repo, verbose)
	if err != nil || exit {
		return err
	}
	if err = validateForEachData(repo, data); err != nil {
		return err
	}
	runProgram := forEachProgram(data, dryRun)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               forEachCmd,
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
		Connector:               None[forgedomain.Connector](),
		Detached:                true,
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

type forEachData struct {
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	config             config.ValidatedConfig
	// connector                       Option[forgedomain.Connector]
	dialogTestInputs components.TestInputs
	// grandParentBranch               gitdomain.LocalBranchName
	hasOpenChanges bool
	initialBranch  gitdomain.LocalBranchName
	// initialBranchFirstCommitMessage Option[gitdomain.CommitMessage]
	// initialBranchInfo               gitdomain.BranchInfo
	// initialBranchProposal           Option[forgedomain.Proposal]
	// initialBranchType               configdomain.BranchType
	// offline                         configdomain.Offline
	// parentBranch                    gitdomain.LocalBranchName
	// parentBranchFirstCommitMessage  Option[gitdomain.CommitMessage]
	// parentBranchInfo                gitdomain.BranchInfo
	// parentBranchProposal            Option[forgedomain.Proposal]
	// parentBranchType                configdomain.BranchType
	// prefetchBranchesSnapshot        gitdomain.BranchesSnapshot
	previousBranch Option[gitdomain.LocalBranchName]
	// remotes                         gitdomain.Remotes
	stashSize gitdomain.StashSize
}

func determineForEachData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (forEachData, bool, error) {
	preFetchBranchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return forEachData{}, false, err
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return forEachData{}, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Detached:              true,
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
		return forEachData{}, exit, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return forEachData{}, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	connectorOpt, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return forEachData{}, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
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
		return forEachData{}, exit, err
	}
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(initialBranch).Get()
	if !hasParentBranch {
		return forEachData{}, false, fmt.Errorf(messages.MergeNoParent, initialBranch)
	}
	grandParentBranch, hasGrandParentBranch := validatedConfig.NormalConfig.Lineage.Parent(parentBranch).Get()
	if !hasGrandParentBranch {
		return forEachData{}, false, fmt.Errorf(messages.MergeNoGrandParent, initialBranch, parentBranch)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return forEachData{}, false, err
	}
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return forEachData{}, false, fmt.Errorf(messages.BranchInfoNotFound, initialBranch)
	}
	parentBranchInfo, hasParentBranchInfo := branchesSnapshot.Branches.FindByLocalName(parentBranch).Get()
	if !hasParentBranchInfo {
		return forEachData{}, false, fmt.Errorf(messages.BranchInfoNotFound, parentBranch)
	}
	initialBranchFirstCommitMessage, err := repo.Git.FirstCommitMessageInBranch(repo.Backend, initialBranch.BranchName(), parentBranch.BranchName())
	if err != nil {
		return forEachData{}, false, err
	}
	initialBranchType := validatedConfig.BranchType(initialBranch)
	parentBranchType := validatedConfig.BranchType(parentBranch)
	parentBranchFirstCommitMessage, err := repo.Git.FirstCommitMessageInBranch(repo.Backend, parentBranch.BranchName(), grandParentBranch.BranchName())
	if err != nil {
		return forEachData{}, false, err
	}
	initialBranchProposal := None[forgedomain.Proposal]()
	parentBranchProposal := None[forgedomain.Proposal]()
	if connector, hasConnector := connectorOpt.Get(); hasConnector {
		if findProposal, canFindProposal := connector.FindProposalFn().Get(); canFindProposal {
			initialBranchProposal, err = findProposal(initialBranch, parentBranch)
			if err != nil {
				print.Error(err)
			}
			parentBranchProposal, err = findProposal(initialBranch, parentBranch)
			if err != nil {
				print.Error(err)
			}
		}
	}
	return forEachData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		config:             validatedConfig,
		dialogTestInputs:   dialogTestInputs,
		// grandParentBranch:               grandParentBranch,
		hasOpenChanges: repoStatus.OpenChanges,
		initialBranch:  initialBranch,
		// initialBranchFirstCommitMessage: initialBranchFirstCommitMessage,
		// initialBranchInfo:               *initialBranchInfo,
		// initialBranchProposal:           initialBranchProposal,
		// initialBranchType:               initialBranchType,
		// offline:                         repo.IsOffline,
		// parentBranch:                    parentBranch,
		// parentBranchFirstCommitMessage:  parentBranchFirstCommitMessage,
		// parentBranchInfo:                *parentBranchInfo,
		// parentBranchProposal:            parentBranchProposal,
		// parentBranchType:                parentBranchType,
		// prefetchBranchesSnapshot:        preFetchBranchesSnapshot,
		previousBranch: previousBranch,
		// remotes:                         remotes,
		stashSize: stashSize,
	}, false, err
}

func forEachProgram(data forEachData, dryRun configdomain.DryRun) program.Program {
	prog := NewMutable(&program.Program{})

	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

func validateForEachData(repo execute.OpenRepoResult, data forEachData) error {
	return nil
}
