package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/flags"
	"github.com/git-town/git-town/v20/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v20/internal/config"
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/execute"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/undo/undoconfig"
	"github.com/git-town/git-town/v20/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v20/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v20/internal/vm/opcodes"
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
Stops when the shell command exits with an error.
You can continue and undo Git operations.

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
	addAllFlag, readAllFlag := flags.All("sync all local branches")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addStackFlag, readStackFlag := flags.Stack("sync the stack that the current branch belongs to")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     forEachCmd,
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   forEachDesc,
		Long:    cmdhelpers.Long(forEachDesc, forEachHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			allBranches, err := readAllFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			stack, err := readStackFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeForEach(dryRun, allBranches, stack, verbose)
		},
	}
	addAllFlag(&cmd)
	addDryRunFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeForEach(dryRun configdomain.DryRun, allBranches configdomain.AllBranches, fullStack configdomain.FullStack, verbose configdomain.Verbose) error {
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
	data, exit, err := determineForEachData(repo, allBranches, fullStack, verbose)
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
	allBranches        configdomain.AllBranches
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	branchesToIterate  gitdomain.LocalBranchNames
	config             config.ValidatedConfig
	dialogTestInputs   components.TestInputs
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	previousBranch     Option[gitdomain.LocalBranchName]
	fullStack          configdomain.FullStack
	stashSize          gitdomain.StashSize
}

func determineForEachData(repo execute.OpenRepoResult, all configdomain.AllBranches, stack configdomain.FullStack, verbose configdomain.Verbose) (forEachData, bool, error) {
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
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return forEachData{}, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		Connector:          None[forgedomain.Connector](),
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
	perennialBranchNames := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	branchesToIterate := gitdomain.LocalBranchNames{}
	switch {
	case all.Enabled():
		branchesToIterate = localBranches
	case stack.Enabled():
		branchesToIterate = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialBranchNames)
	}
	return forEachData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		branchesToIterate:  branchesToIterate,
		config:             validatedConfig,
		dialogTestInputs:   dialogTestInputs,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		previousBranch:     previousBranch,
		stashSize:          stashSize,
	}, false, err
}

func forEachProgram(data forEachData, dryRun configdomain.DryRun) program.Program {
	prog := NewMutable(&program.Program{})
	for _, branchToIterate := range data.branchesToIterate {
		prog.Value.Add(&opcodes.ExecuteShellCommand{})
	}

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
	if data.allBranches.Enabled() && data.fullStack.Enabled() {
		return fmt.Errorf("Please don't enable both --all or --stack, just one of them")
	}
	if !data.allBranches.Enabled() && !data.fullStack.Enabled() {
		return fmt.Errorf("Please enable either --all or --stack")
	}
	return nil
}
