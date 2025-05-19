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
	walkCmd  = "walk"
	walkDesc = "Perform a shell operation on each branch"
	walkHelp = `
If a shell command is given, executes it on each branch.
Stops when the shell command exits with an error.

If no shell command is given, exits to the shell on each branch.
In that case, you can run "git town continue", "git town skip", or "git town undo"
to retry, ignore, or abort the iteration.


Consider this stack:

main
 \
  branch-1
   \
    branch-2

When running "git town walk --stack echo hello",
it prints this output:

[main] hello

[branch-1] hello

[branch-2] hello
`
)

func walkCommand() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("iterate all local branches")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addStackFlag, readStackFlag := flags.Stack("iterate all branches in the current stack")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     walkCmd,
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   walkDesc,
		Long:    cmdhelpers.Long(walkDesc, walkHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
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
			return executeWalk(args, dryRun, allBranches, stack, verbose)
		},
	}
	addAllFlag(&cmd)
	addDryRunFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeWalk(args []string, dryRun configdomain.DryRun, allBranches configdomain.AllBranches, fullStack configdomain.FullStack, verbose configdomain.Verbose) error {
	err := validateArgs(allBranches, fullStack)
	if err != nil {
		return err
	}
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
	data, exit, err := determineWalkData(args, repo, allBranches, fullStack, verbose)
	if err != nil || exit {
		return err
	}
	// if err = validateWalkData(repo, data); err != nil {
	// 	return err
	// }
	runProgram := walkProgram(data, dryRun)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               walkCmd,
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

type walkData struct {
	allBranches        configdomain.AllBranches
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	branchesToWalk     gitdomain.LocalBranchNames
	commandToExecute   []string
	config             config.ValidatedConfig
	connector          Option[forgedomain.Connector]
	dialogTestInputs   components.TestInputs
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	previousBranch     Option[gitdomain.LocalBranchName]
	fullStack          configdomain.FullStack
	stashSize          gitdomain.StashSize
}

func determineWalkData(args []string, repo execute.OpenRepoResult, all configdomain.AllBranches, stack configdomain.FullStack, verbose configdomain.Verbose) (walkData, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return walkData{}, false, err
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
		return walkData{}, exit, err
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return walkData{}, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	connector, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return walkData{}, false, err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
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
		return walkData{}, exit, err
	}
	perennialBranchNames := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	branchesToWalk := gitdomain.LocalBranchNames{}
	switch {
	case all.Enabled():
		branchesToWalk = localBranches
	case stack.Enabled():
		branchesToWalk = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialBranchNames)
	}
	return walkData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		branchesToWalk:     branchesToWalk,
		commandToExecute:   args,
		config:             validatedConfig,
		connector:          connector,
		dialogTestInputs:   dialogTestInputs,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		previousBranch:     previousBranch,
		stashSize:          stashSize,
	}, false, err
}

func walkProgram(data walkData, dryRun configdomain.DryRun) program.Program {
	prog := NewMutable(&program.Program{})
	// for _, branchToWalk := range data.branchesToWalk {
	// 	prog.Value.Add(&opcodes.ExecuteShellCommand{})
	// }
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

func validateArgs(all configdomain.AllBranches, stack configdomain.FullStack) error {
	if all.Enabled() && stack.Enabled() {
		return fmt.Errorf("Please provide either --all or --stack, not both of them")
	}
	if !all.Enabled() && !stack.Enabled() {
		return fmt.Errorf("Please provide either --all or --stack")
	}
	return nil
}
